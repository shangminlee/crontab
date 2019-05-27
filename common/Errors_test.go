package common

import (
    "testing"
)

var testErrs = []struct{
    err error
}{
    {ERR_LOCK_ALREADY_REQUIRED},
    {ERR_NO_LOCAL_IP_FOUND},
}

func deferRecover(t *testing.T)  {
    r := recover()
    if err, ok := r.(error); ok {
        t.Logf("Test Error Success : %v", err)
    } else {
        t.Errorf("Testing Error Not Pass : %v", err)
    }
}

func testError(t *testing.T, err error) {
    defer deferRecover(t)
    panic(err)
}

func TestErrors(t *testing.T)  {

    for _ , err := range testErrs{
        testError(t, err.err)
    }

}

