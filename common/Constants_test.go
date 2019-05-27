package common

import "testing"

// 测试数据
var testConsData = []struct{
    code int
}{
    {1},
    {2},
    {3},
}

func TestConstants(t *testing.T)  {
    if JOB_EVENT_SAVE != testConsData[0].code {
        t.Errorf("JOB_EVENT_SAVE actual value : %d expected value : %d", JOB_EVENT_SAVE, testConsData[0].code)
    }
    if JOB_EVENT_DELETE != testConsData[1].code {
        t.Errorf("JOB_EVENT_DELETE actual value : %d expected value : %d", JOB_EVENT_DELETE, testConsData[1].code)
    }
    if JOB_EVENT_KILL != testConsData[2].code {
        t.Errorf("JOB_EVENT_KILL actual value : %d expected value : %d", JOB_EVENT_KILL, testConsData[2].code)
    }
}