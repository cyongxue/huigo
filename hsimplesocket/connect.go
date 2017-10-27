package hsimplesocket

import (
    "fmt"
    "net"
    "encoding/binary"
    "bytes"
    "io"
    "time"
)

// @Description: 封装net.Conn，以便对于read、write的封装
type Connection struct {
    Conn            net.Conn
}

const (
    CONNECTION_MAGIC = 0xFFFF
    CONNECTION_SIZE_BUF = 6

    CONN_MAX_DATA_LEN = 8 * 1024 * 1024
)

/**
 @Description：封装协议链接的read操作
 @Param:
 @Return：
        []byte          数据结果输出
        error           错误输出
 */
func (c *Connection) Read() ([]byte, error) {
    // 设置read超时
    c.Conn.SetReadDeadline(time.Now().Add(70 * time.Second))

    // 先读取长度
    lenData := make([]byte, CONNECTION_SIZE_BUF)
    _, err := io.ReadFull(c.Conn, lenData)
    if err != nil {
        return nil, fmt.Errorf("socket read data length error: %s", err.Error())
    }
    // 从byte中解析出l值
    magic := binary.BigEndian.Uint16(lenData[0:2])
    if magic != CONNECTION_MAGIC {
        return nil, fmt.Errorf("socket read data magic error: %x", magic)
    }
    l := binary.BigEndian.Uint32(lenData[2:CONNECTION_SIZE_BUF])
    if l > CONN_MAX_DATA_LEN {
        return nil, fmt.Errorf("data len big: %d", l)
    }

    // 准备读取数据
    d := make([]byte, l)
    realLen, err := io.ReadFull(c.Conn, d)
    if err != nil {
        return nil, fmt.Errorf("socket read data(len=%d) error: %s", l, err.Error())
    }
    if realLen != int(l) {
        return nil, fmt.Errorf("data len is error: reallen(%d) != len(%d)", realLen, l)
    }

    return d, nil
}

/**
 @Description：封装协议链接的write操作
 @Param:
 @Return：
 */
func (c *Connection) Write(data []byte) error {

    // 准备长度buf
    lenData := len(data)
    if lenData > CONN_MAX_DATA_LEN {
        return fmt.Errorf("data length more len")
    }

    bufLen := make([]byte, CONNECTION_SIZE_BUF)
    binary.BigEndian.PutUint16(bufLen[0:2], CONNECTION_MAGIC)
    binary.BigEndian.PutUint32(bufLen[2:CONNECTION_SIZE_BUF], uint32(lenData))

    // 拼接处send buffer
    bufs := make([][]byte, 2)
    bufs[0] = bufLen
    bufs[1] = data
    sep := []byte("")

    buf := bytes.Join(bufs, sep)

    // 发送
    willLen := lenData + CONNECTION_SIZE_BUF
    l, err := c.Conn.Write(buf)
    if err != nil {
        return fmt.Errorf("socket write buf error: %s", err.Error())
    }
    if l != willLen {
        return fmt.Errorf("socket write data length error, %d send %d", willLen, l)
    }

    return nil
}

/**
 @Description：tls connection close
 @Param:
 @Return：
 */
func (c *Connection) Close()  {
    c.Conn.Close()
    return
}
