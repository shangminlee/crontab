package master

import (
    "encoding/json"
    "github.com/shangminlee/crontab/common"
    "net"
    "net/http"
    "strconv"
    "time"
)

type ApiServer struct {
    httpServer *http.Server
}

var G_API_SERVER *ApiServer

func InitApiServer() error {

    // 1. 创建路由
    mux := http.NewServeMux()
    mux.HandleFunc("/job/save", handlerJobSave)
    mux.HandleFunc("/job/delete", handlerJobDelete)
    mux.HandleFunc("/job/list", handlerJobList)
    mux.HandleFunc("/job/kill", handlerJobKill)
    mux.HandleFunc("/job/log", handlerJobLog)
    mux.HandleFunc("/worker/list", handlerWorkerList)

    // 静态页面
    staticDir := http.Dir(G_CONFIG.WebRoot)
    staticHandler := http.FileServer(staticDir)
    mux.Handle("/", http.StripPrefix("/", staticHandler))

    // 2. 创建监听
    listen, err := net.Listen("tcp", ":" + strconv.Itoa(G_CONFIG.ApiPort))
    if err != nil {
        return err
    }

    // 3. 创建服务配置
    httpServer := &http.Server{
        ReadTimeout : time.Duration(G_CONFIG.ApiReadTimeout) * time.Millisecond,
        WriteTimeout: time.Duration(G_CONFIG.ApiWriteTimeout) * time.Millisecond,
        Handler     : mux,
    }

    G_API_SERVER = &ApiServer{
        httpServer: httpServer,
    }

    // 4. 启动监听服务
    go func() {
        err := httpServer.Serve(listen)
        if err != nil {
            panic(err)
        }
    }()

    return nil
}

func handlerWorkerList(resp http.ResponseWriter, req *http.Request) {
    var (
        workerArr []string
        err error
        bytes []byte
    )

    if workerArr, err = G_WORKER_MGR.ListWorkers(); err != nil {
        goto ERR
    }

    // 正常应答
    if bytes, err = common.BuildResponse(0, "success", workerArr); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
    return

ERR:
    if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
}

func handlerJobLog(resp http.ResponseWriter, req *http.Request) {
    var (
        err        error
        name       string // 任务名字
        skipParam  string// 从第几条开始
        limitParam string // 返回多少条
        skip       int
        limit      int
        logArr     []*common.JobLog
        bytes      []byte
    )

    // 解析GET参数
    if err = req.ParseForm(); err != nil {
        goto ERR
    }

    // 获取请求参数 /job/log?name=job10&skip=0&limit=10
    name = req.Form.Get("name")
    skipParam = req.Form.Get("skip")
    limitParam = req.Form.Get("limit")
    if skip, err = strconv.Atoi(skipParam); err != nil {
        skip = 0
    }
    if limit, err = strconv.Atoi(limitParam); err != nil {
        limit = 20
    }

    if logArr, err = G_LOG_MGR.ListLog(name, skip, limit); err != nil {
        goto ERR
    }

    // 正常应答
    if bytes, err = common.BuildResponse(0, "success", logArr); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
    return

ERR:
    if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
}

// 杀死任务
func handlerJobKill(resp http.ResponseWriter, req *http.Request) {
    var (
        err error
        name string
        bytes []byte
    )

    // 解析POST表单
    if err = req.ParseForm(); err != nil {
        goto ERR
    }

    // 要杀死的任务名
    name = req.PostForm.Get("name")

    // 杀死任务
    if err = G_JOB_MGR.KillJob(name); err != nil {
        goto ERR
    }

    // 正常应答
    if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
    return

ERR:
    if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
}

// 任务列表
func handlerJobList(resp http.ResponseWriter, req *http.Request) {
    var (
        jobList []*common.Job
        bytes []byte
        err error
    )

    // 获取任务列表
    if jobList, err = G_JOB_MGR.ListJobs(); err != nil {
        goto ERR
    }

    // 正常应答
    if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
    return

ERR:
    if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
}


// 删除任务
func handlerJobDelete(resp http.ResponseWriter, req *http.Request) {
    var (
        err error	// interface{}
        name string
        oldJob *common.Job
        bytes []byte
    )

    // POST:   a=1&b=2&c=3
    if err = req.ParseForm(); err != nil {
        goto ERR
    }

    // 删除的任务名
    name = req.PostForm.Get("name")

    // 去删除任务
    if oldJob, err = G_JOB_MGR.DeleteJob(name); err != nil {
        goto ERR
    }

    // 正常应答
    if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
    return

ERR:
    if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }
}


// 保存任务接口
func handlerJobSave(resp http.ResponseWriter, req *http.Request)  {
    var (
        err     error
        postJob string
        job     common.Job
        oldJob  *common.Job
    )

    // 1. 解析 POST 表单
    if err = req.ParseForm(); err != nil {
        goto ERR
    }

    // 2. 获取表单中的 job 字段
    postJob = req.PostForm.Get("job")

    // 3, 反序列化job
    err = json.Unmarshal([]byte(postJob), &job);
    if err != nil {
        goto ERR
    }

    // 4, 保存到etcd
    oldJob, err = G_JOB_MGR.SaveJob(&job);
    if err != nil {
        goto ERR
    }

    // 5, 返回正常应答 ({"errno": 0, "msg": "", "data": {....}})
    if bytes, err := common.BuildResponse(0, "success", oldJob); err == nil {
        resp.Header().Add("Content-Type","application/json;charset=utf-8")
        _, _ = resp.Write(bytes)
    }

    return

ERR:
    bytes, err := common.BuildResponse(-1, err.Error(), nil )
    resp.Header().Add("Content-Type","application/json;charset=utf-8")
    if err == nil {
        _, err = resp.Write(bytes)
    }
}