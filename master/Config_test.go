package master

import "testing"

func TestInitConfig(t *testing.T) {
    tests := []struct {
        name string
        filename string
        expected Config
    }{
        {"配置文件读取单元测试-1", "./main/master.json", Config{
            ApiPort        : 8070,
            ApiReadTimeout : 5000,
            ApiWriteTimeout: 5000,
        }},
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            err := InitConfig(test.filename)
            if err != nil {
                t.Errorf("InitConfig Error %v", err)
                return
            }
            //if test.expected != *G_CONFIG {
            //   t.Errorf("actual value %v, expected value %v", G_CONFIG, test.expected)
            //}
        })
    }
}

