package hrpc

import (
	"encoding/gob"
	"bufio"
	"net"
	"time"
	"log"
)

type GobServerCodec struct {
	rwc    net.Conn
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
	closed bool
}

func NewServerCodec(conn net.Conn) ServerCodec {
	buf := bufio.NewWriter(conn)
	return &GobServerCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
	}
}

func (c *GobServerCodec) RemoteAddr() string {
	return c.rwc.RemoteAddr().String()
}

func (c *GobServerCodec) ReadRequestHeader(r *Request, readTimeout time.Duration) error {
	if readTimeout != 0 {
		c.rwc.SetDeadline(time.Now().Add(readTimeout))
	}
	return c.dec.Decode(r)
}

func (c *GobServerCodec) ReadRequestBody(body interface{}, readTimeout time.Duration) error {
	if readTimeout != 0 {
		c.rwc.SetDeadline(time.Now().Add(readTimeout))
	}
	return c.dec.Decode(body)
}

func (c *GobServerCodec) WriteResponse(r *Response, body interface{}, writeTimeout time.Duration) (err error) {
	if writeTimeout != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(writeTimeout))
	}

	if err = c.enc.Encode(r); err != nil {
		if c.encBuf.Flush() == nil {
			// Gob couldn't encode the header. Should not happen, so if it does,
			// shut down the connection to signal that the connection is broken.
			log.Println("rpc: gob error encoding response:", err)
			c.Close()
		}
		return
	}
	if err = c.enc.Encode(body); err != nil {
		if c.encBuf.Flush() == nil {
			// Was a gob problem encoding the body but the header has been written.
			// Shut down the connection to signal that the connection is broken.
			log.Println("rpc: gob error encoding body:", err)
			c.Close()
		}
		return
	}
	return c.encBuf.Flush()
}

func (c *GobServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}
