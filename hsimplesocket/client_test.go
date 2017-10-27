package hsimplesocket

import (
    "testing"
    "fmt"
)

func TestTlsClient_Connect(t *testing.T) {
    client := TlsClient{
        ServerAddr: "127.0.0.1:1997",
    }

    err := client.Connect()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    err = client.Conn.Write([]byte("hello 你好。。。。。"))
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    data, err := client.Conn.Read()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(string(data))
    return
}