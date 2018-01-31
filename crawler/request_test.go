package crawler

import (
    "testing"
    "fmt"
)

func TestRequest_Run(t *testing.T) {
    req := Request{
        StrURL: "http://www.baidu.com",
        Method: "get",
    }
    err := req.Run()
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println(req.Response)
    return
}
