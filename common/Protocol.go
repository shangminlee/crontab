package common

import (
    "context"
    "encoding/json"
    "github.com/gorhill/cronexpr"
    "strings"
    "time"
)

// 定时任务
type Job struct {
    Name     string `json:"name"`      // 任务名
    Command  string `json:"command"`   // shell 命令
    CronExpr string `json:"cron_expr"` // cron 表达式
}

// 任务调度计划
type JobSchedulePlan struct {
    Job      *Job                 // 要调度的任务信息
    Expr     *cronexpr.Expression // 解析好的 cron expr 表达式
    NextTime time.Time            // 下次调度时间
}

// 任务执行状态
type JobExecuteInfo struct {
    Job        *Job                // 任务信息
    PlanTime   time.Time           // 计划的调度时间
    RealTime   time.Time           // 实际的调度时间
    CancelCtx  context.Context     // 任务 command 的 context
    CancelFunc context.CancelFunc  // 用于取消 command 执行的 cancel 函数
}

// Http 接口应答
type Response struct {
    ErrorCode int         `json:"error_code"` // 错误代码，0 成功， 大于 0 失败
    Msg       string      `json:"msg"`        // 错误代码
    Data      interface{} `json:"data"`       // 应答数据
}

// 变化事件
 type JobEvent struct {
     EventType int  // SAVE、DELETE、KILL
     Job      *Job // 任务信息
 }

// 任务执行结果
type JobExecuteResult struct {
    ExecuteInfo *JobExecuteInfo // 执行状态
    Output      []byte          // 脚本输出
    Err         error           // 脚本错误原因
    StartTime   time.Time       // 启动时间
    EndTime     time.Time       // 结束时间
}

// 任务执行日志
type JobLog struct {
    JobName      string `json:"job_name"      bson:"job_name"`      // 任务名字
    Command      string `json:"command"       bson:"command"`       // 脚本命令
    Err          string `json:"err"           bson:"err"`           // 错误原因
    Output       string `json:"output"        bson:"output"`        // 脚本输出
    PlanTime     int64  `json:"plan_time"     bson:"plan_time"`     // 计划开始时间
    ScheduleTime int64  `json:"schedule_time" bson:"schedule_time"` // 实际调度时间
    StartTime    int64  `json:"start_time"    bson:"start_time"`    // 任务执行开始时间
    EndTime      int64  `json:"end_time"      bson:"end_time"`      // 任务执行结束时间
}

// 日志批次
type LogBatch struct {
    Logs []interface{} // 多条日志
}

// 任务日志过滤条件
type JobLogFilter struct {
    JobName string `bson:"job_name"`
}

// 任务日志排序条件
type SortLogByStartTime struct {
    SortOrder int `bson:"start_time"`
}

// 应答方法
func BuildResponse(errorCode int, msg string, data interface{})([]byte, error)  {
    resObj := Response{
        ErrorCode: errorCode,
        Msg      : msg,
        Data     : data,
    }
    return json.Marshal(resObj)
}

// 反序列化 Job
func UnpackJob(value []byte) (*Job, error) {
    job := &Job{}
    err := json.Unmarshal(value, job)
    if err != nil {
        return nil, err
    }
    return job, nil
}

// 从 etcd 的 key 中提取任务名
// /cron/jobs/job10 抹掉 /cron/jobs/
func ExtractJobName (jobKey string) string {
    return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// /cron/killer/job10 抹掉 /cron/jobs/
func ExtractKillerName(killerKey string) string  {
    return strings.TrimPrefix(killerKey, JOB_KILLER_DIR)
}

// 任务变化事件 1 更新任务 2 删除任务
func BuildJobEvent(eventType int, job *Job) *JobEvent {
    return &JobEvent{
        EventType: eventType,
        Job     : job,
    }
}

// 构造任务执行计划
func BuildJobSchedulePlan(job *Job) (*JobSchedulePlan, error)  {
    expr, err := cronexpr.Parse(job.CronExpr)
    if err != nil {
        return nil, err
    }
    plan := &JobSchedulePlan{
        Job     : job,
        Expr    : expr,
        NextTime: expr.Next(time.Now()),
    }
    return plan, nil
}

// 构造执行状态信息
func BuildJobExecuteInfo(plan *JobSchedulePlan) *JobExecuteInfo {
    info := &JobExecuteInfo{
        Job     : plan.Job,
        PlanTime: plan.NextTime,
        RealTime: time.Now(),
    }
    info.CancelCtx, info.CancelFunc = context.WithCancel(context.TODO())
    return info
}

// 提取 worker 的 IP
func ExtractWorkerIP(regKey string) string {
    return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}