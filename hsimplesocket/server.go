package hsimplesocket

import (
    "crypto/tls"
    "net"
    "fmt"
    "github.com/golang/glog"
    "time"
)

type ServiceManager interface {
    Handler(conn Connection)
}

/**
 @Description: tcp server的封装
 */
type TlsServer struct {
    ServiceMg    ServiceManager

    Addr         string             // 服务端监听的地址
    PemFile      string
    KeyFile      string

    Listener     net.Listener
}

/**
 @Description：tls server的set
 @Param:
 @Return：
 */
func (t *TlsServer) SetTlsServer(addr string, pemFile string, keyFile string, serviceMg ServiceManager) {
    t.ServiceMg = serviceMg
    t.Addr = addr
    t.PemFile = pemFile
    t.KeyFile = keyFile
}

/**
 @Description：tcp server的listen
 @Param:
 @Return：
 */
func (t *TlsServer) Listen() error {

    if t.Listener != nil {
        return fmt.Errorf("The server is already started")
    }

    // 加载证书
    cert, err := tls.LoadX509KeyPair(t.PemFile, t.KeyFile)
    if err != nil {
        return err
    }
    config := &tls.Config{Certificates: []tls.Certificate{cert}}

    t.Listener, err = tls.Listen("tcp", t.Addr, config)
    if err != nil {
        return err
    }

    // 创建协程执行loop
    go t.loop()
    return nil
}

/**
 @Description：tcp server的Loop，一直循环
 @Param:
 @Return：
 */
func (t *TlsServer) loop() {
    glog.Infof("tls server start loop ......")
    for {
        client, err := t.Listener.Accept()
        if err != nil {
            continue
        }

        //client.(*net.TCPConn).SetKeepAlive(true)
        client.SetReadDeadline(time.Now().Add(70 * time.Second))
        // 处理单个连接connection
        conn := Connection{
            Conn: client,
        }
        go t.ServiceMg.Handler(conn)
    }

    t.Listener.Close()
    return
}


// @Description: tls server的对象
var tlsServer *TlsServer

/**
 @Description：获取tcp server实例
 @Param:
 @Return：
 */
func TlsServerInstance() *TlsServer {
    return tlsServer
}

/**
 @Description：tcp server的初始化接口
 @Param:
 @Return：
 */
func init()  {
    tlsServer = new(TlsServer)
    if tlsServer == nil {
        panic("new TlsServer failed")
    }
}