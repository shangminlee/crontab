package master

import (
    "encoding/json"
    "io/ioutil"
)

var G_CONFIG *Config

// 程序配置
type Config struct {
    ApiPort               int      `json:"apiPort"`
    ApiReadTimeout        int      `json:"apiReadTimeout"`
    ApiWriteTimeout       int      `json:"apiWriteTimeout"`
    EtcdEndpoints         []string `json:"etcdEndpoints"`
    EtcdDialTimeout       int      `json:"etcdDialTimeout"`
    WebRoot               string   `json:"webroot"`
    MongodbUri            string   `json:"mongodbUri"`
    MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
}

// 读取配置文件
func InitConfig(filename string) error {

    // 1. 读取配置文件
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }

    // 2. 反序列化配置文件
    var conf Config
    err = json.Unmarshal(content, &conf )
    if err != nil {
        return err
    }

    G_CONFIG = &conf

    return nil
}
