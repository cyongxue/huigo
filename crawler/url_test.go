package crawler

import (
    "fmt"
    "testing"
)

func TestUrl_GetFileExt(t *testing.T) {
    u := Url{
        StrUrl: "http://www.anquanbao.com/book/index?id=1&hehui=girl#top",
    }
    err := u.Init()
    if err != nil {
        fmt.Errorf(err.Error())
        return
    }

    fmt.Println(u.GetFileName())
    fmt.Println(u.GetFileExt())
    return
}
