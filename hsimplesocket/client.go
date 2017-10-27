package hsimplesocket

import (
    "crypto/tls"
)

// @Description: 客户端的数据结构
type TlsClient struct {
    ServerAddr          string
    Conn                Connection
}

/**
 @Description：tls server的set
 @Param:
 @Return：
 */
func (t *TlsClient) Connect() error {
    conf := &tls.Config{
        InsecureSkipVerify: true,
    }

    conn, err := tls.Dial("tcp", t.ServerAddr, conf)
    if err != nil {
        return err
    }

    t.Conn.Conn = conn
    return nil
}
