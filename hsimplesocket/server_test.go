package hsimplesocket

import (
    "testing"
    "fmt"
)

type ServiceManage struct {

}

func (s *ServiceManage) Handler(conn Connection) {

    // 处理数据
    for {
        data, err := conn.Read()
        if err != nil {
            conn.Close()
            return
        }

        fmt.Println(string(data))
        err = conn.Write(data)
        if err != nil {
            fmt.Println(err.Error())
            conn.Close()
            return
        }
    }
    return
}

func TestTlsServer_Listen(t *testing.T) {

    service := new(ServiceManage)
    TlsServerInstance().SetTlsServer("127.0.0.1:1997", "xxxx", "xxxxx", service)
    TlsServerInstance().Listen()
    select {
    }
    return 
}