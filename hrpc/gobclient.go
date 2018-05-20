package hrpc

import (
	"encoding/gob"
	"bufio"
	"net"
	"time"
)

func NewClientCodec(conn net.Conn) ClientCodec {
	encBuf := bufio.NewWriter(conn)
	return &GobClientCodec{
		rwc:conn,
		dec: gob.NewDecoder(conn),
		enc: gob.NewEncoder(encBuf),
		encBuf: encBuf}
}

type GobClientCodec struct {
	rwc    net.Conn
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
}

func (c *GobClientCodec) RemoteAddr() string {
	return c.rwc.RemoteAddr().String()
}

func (c *GobClientCodec) WriteRequest(r *Request, body interface{}, writeTimeout time.Duration) (err error) {
	if writeTimeout != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(writeTimeout))
	}

	if err = c.enc.Encode(r); err != nil {
		return
	}

	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *GobClientCodec) ReadResponseHeader(r *Response, readTimeout time.Duration) error {
	if readTimeout != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(readTimeout))
	}
	return c.dec.Decode(r)
}

func (c *GobClientCodec) ReadResponseBody(body interface{}, readTimeout time.Duration) error {
	if readTimeout != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(readTimeout))
	}
	return c.dec.Decode(body)
}

func (c *GobClientCodec) Close() error {
	return c.rwc.Close()
}