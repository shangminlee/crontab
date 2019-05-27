package master

import (
    "context"
    "github.com/shangminlee/crontab/common"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
)

type LogMgr struct {
    client        *mongo.Client
    logCollection *mongo.Collection
}

var G_LOG_MGR *LogMgr

func InitLogMgr() error{
    ctx, _ := context.WithTimeout(context.TODO(), 10 * time.Second)

    client, err := mongo.Connect(
        ctx,
        options.Client().ApplyURI(G_CONFIG.MongodbUri),
        options.Client().SetConnectTimeout(time.Duration(G_CONFIG.MongodbConnectTimeout) * time.Millisecond),
    )
    if err != nil {
        return err
    }

    G_LOG_MGR = &LogMgr{
        client       : client,
        logCollection: client.Database("cron").Collection("log"),
    }
    return nil
}

func (logMgr *LogMgr) ListLog(name string, skip int, limit int) (logArr []*common.JobLog, err error){
    var (
        filter  *common.JobLogFilter
        logSort *common.SortLogByStartTime
        cursor  *mongo.Cursor
        jobLog  *common.JobLog
    )

    // len(logArr)
    logArr = make([]*common.JobLog, 0)

    // 过滤条件
    filter = &common.JobLogFilter{JobName: name}

    // 按照任务开始时间倒排
    logSort = &common.SortLogByStartTime{SortOrder: -1}

    var findOps = options.FindOptions{}
    // 查询
    findOps.SetSort(logSort).SetSkip(int64(skip)).SetLimit(int64(limit))
    if cursor, err = logMgr.logCollection.Find(context.TODO(),filter, &findOps); err != nil {
        return
    }
    // 延迟释放游标
    defer func() {
        _ = cursor.Close(context.TODO())
    }()

    for cursor.Next(context.TODO()) {
        jobLog = &common.JobLog{}

        // 反序列化BSON
        if err = cursor.Decode(jobLog); err != nil {
            continue // 有日志不合法
        }

        logArr = append(logArr, jobLog)
    }
    return
}