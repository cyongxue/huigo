package main

import (
    "github.com/hprose/hprose-golang/rpc"
    "fmt"
)

type Apple struct {
    Name       string
    Score      int64
    ID         string
}

type RemoteApples struct {
    Init            func()
    GetAll          func() map[string]*Apple
}

func main()  {
    client := rpc.NewHTTPClient("http://:8080")

    //var ro *RemoteObject
    //client.UseService(&ro)

    //ob, err := ro.GetOne("hjkhsbnmn123")
    //if err != nil {
    //    fmt.Println(err.Error())
    //} else {
    //    fmt.Println(ob.ObjectId)
    //    fmt.Println(ob.Score)
    //    fmt.Println(ob.PlayerName)
    //}
    //
    //obs := ro.GetAll()
    //for _, item := range obs {
    //    fmt.Println(item.ObjectId)
    //    fmt.Println(item.Score)
    //    fmt.Println(item.PlayerName)
    //}
    //
    //fmt.Println("*******************")
    var ra *RemoteApples
    client.UseService(&ra)
    ra.Init()
    as := ra.GetAll()
    for _, item := range as {
        fmt.Println(item.ID)
        fmt.Println(item.Score)
        fmt.Println(item.Name)
    }

}