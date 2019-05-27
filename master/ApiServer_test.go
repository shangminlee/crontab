package master

import "testing"

func Test(t *testing.T) {
    tests := []struct {
        name string
    }{
        {"测试初始化服务"},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            err := InitConfig("./main/master.json")
            err = InitApiServer()
            if err != nil {
                t.Errorf("InitApiServer Error %v", err)
                return
            }
            if G_API_SERVER.httpServer == nil {
                t.Errorf("G_API_SERVER is nil")
            }
        })
    }
}
