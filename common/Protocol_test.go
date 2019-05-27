package common

import "testing"

// 测试应答方法
func TestBuildResponse(t *testing.T) {
    tests := []struct {
        name      string
        errorCode int
        msg       string
        data      interface{}
        expected  string
    }{
        {"测试成功",0,"成功",Job{
            Name: "job1",
            Command: "echo hello",
            CronExpr: "* * * * *",
        }, `{"error_code":0,"msg":"成功","data":{"name":"job1","command":"echo hello","cron_expr":"* * * * *"}}`},
        {"测试失败",1,"失败",
            "", `{"error_code":1,"msg":"失败","data":""}`},

    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            content, err :=BuildResponse(test.errorCode, test.msg, test.data)
            if err != nil{
                t.Errorf("Test BuildResponse %v", err)
                return
            }
            if string(content) != test.expected {
                t.Errorf("Test BuildResponse actual value is %s , " +
                    "expected value is %s", content, test.expected)
            }
        })
    }
}

func TestUnpackJob(t *testing.T) {
    tests := []struct {
        name     string
        jobJson  string
        expected *Job
    }{
        {"测试反序列化-1",`{"name":"job1","command":"echo hello","cron_expr":"* * * * *"}`, &Job{
            Name: "job1",
            Command: "echo hello",
            CronExpr: "* * * * *",
        }},
        {"测试反序列化-2", "{}", &Job{}},
        {"测试反序列化失败-1", "oaosf", nil},
        {"测试反序列化失败-2", "{kkkk}", nil},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            job, err := UnpackJob([]byte(test.jobJson))

            if test.expected != nil {
                if err != nil {
                    t.Errorf("Test UnpackJob Error %v", err)
                    return
                }
                if *job != *test.expected {
                    t.Errorf("Test UnpackJob response actual value is %v , "+
                        "expected value is %v", *job, *test.expected)
                }
            } else {
                if err == nil {
                    t.Errorf("Test UnpackJob Error Fail expected err not nil")
                }
                if job != test.expected {
                    t.Errorf("Test UnpackJob response actual value is %v , "+
                        "expected value is %v", job, test.expected)
                }
            }
        })
    }
}

func TestExtractJobName(t *testing.T) {
    tests := []struct {
        name string
        jobKey string
        expected string
    }{
        {"测试获取JobName","/cron/jobs/job10","job10"},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            res := ExtractJobName(test.jobKey)
            if res != test.expected {
                t.Errorf("actual is %s, expected is %s", res, test.expected)
            }
        })
    }
}

func TestExtractKillerName(t *testing.T) {
    tests := []struct {
        name string
        killerKey string
        expected string
    }{
        {"测试获取KillerName","/cron/killer/job10","job10"},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            res := ExtractKillerName(test.killerKey)
            if res != test.expected {
                t.Errorf("actual is %s, expected is %s", res, test.expected)
            }
        })
    }
}

func TestExtractWorkerIP(t *testing.T) {
    tests := []struct {
        name string
        workKey string
        expected string
    }{
        {"测试获取WorkerIP","/cron/workers/127.0.0.1","127.0.0.1"},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            res := ExtractWorkerIP(test.workKey)
            if res != test.expected {
                t.Errorf("actual is %s, expected is %s", res, test.expected)
            }
        })
    }
}

