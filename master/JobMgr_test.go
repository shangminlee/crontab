package master

import (
    "fmt"
    "github.com/shangminlee/crontab/common"
    "testing"
)

// 测试 Job Manager
func TestJobMgr(t *testing.T) {
    tests := []struct {
        name string
        job  *common.Job
    }{
        {"测试Etcd管理-1", &common.Job{
            Name: "job1",
            Command: "echo hello world",
            CronExpr: "*/5 * * * *",
        }},
        {"测试Etcd管理-2", &common.Job{
            Name: "job2",
            Command: "echo hello world",
            CronExpr: "*/5 * * * *",
        }},
        {"测试Etcd管理-3", &common.Job{
            Name: "job3",
            Command: "echo hello world",
            CronExpr: "*/5 * * * *",
        }},

    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            err := InitConfig("./main/master.json")
            if err != nil {
                t.Errorf("Test Job Mgr Init Config Error %v", err)
            }

            err = InitJobMgr()

            if err != nil {
                t.Errorf("Test Job Mgr Init JobMgr Error %v", err)
            }

            oldJob, err := G_JOB_MGR.SaveJob(test.job)
            fmt.Println(oldJob, err)

            jobList, err := G_JOB_MGR.ListJobs()
            for _, job := range jobList {
                fmt.Println("JobList : ", *job)
            }

            oldJob, err = G_JOB_MGR.DeleteJob(test.job.Name)
            fmt.Println("OldJob : ",oldJob)
        })
    }
}
