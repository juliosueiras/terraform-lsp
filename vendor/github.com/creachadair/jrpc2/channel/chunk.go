package channel

import (
	"bufio"
	"bytes"
	"io"
)

const chunkMaxBytes = 65535

type chunked struct {
	r  *bufio.Reader
	wc io.WriteCloser

	// 2-byte length + data bytes
	rbuf [2 + chunkMaxBytes]byte
	wbuf [2 + chunkMaxBytes]byte
}

// Chunked is a framing discipline that transmits messages in fixed-size chunks
// of up to 65535 bytes, each prefixed by a length tag. All chunks except the
// last are exactly 65535 bytes; the last chunk contains the remaining
// (possibly empty) tail of the message.
//
// The length tag is a two-byte unsigned big-endian integer.
func Chunked(r io.Reader, wc io.WriteCloser) Channel {
	return &chunked{r: bufio.NewReader(r), wc: wc}
}

// Send implements part of the Channel interface. It writes the message in one
// or more chunks to the underlying writer.
func (c *chunked) Send(msg []byte) error {
	buf := c.wbuf[:]
	rest := msg
	buf[0] = 0xff
	buf[1] = 0xff
	for len(rest) >= chunkMaxBytes {
		n := copy(buf[2:], rest)
		rest = rest[n:]
		if _, err := c.wc.Write(buf); err != nil {
			return err
		}
	}
	left := len(rest)
	buf[0] = byte((left >> 8) & 0xff)
	buf[1] = byte((left >> 0) & 0xff)
	copy(buf[2:], rest)
	_, err := c.wc.Write(buf[:left+2])
	return err
}

// Recv implements part of the Channel interface.
func (c *chunked) Recv() ([]byte, error) {
	buf := c.rbuf[:]
	var out bytes.Buffer
	for {
		if _, err := io.ReadFull(c.r, buf[:2]); err != nil {
			return nil, err
		}
		clen := int(buf[0])*256 + int(buf[1])
		if clen > 0 {
			if _, err := io.ReadFull(c.r, buf[:clen]); err != nil {
				return nil, err
			}
			out.Write(buf[:clen])
		}
		if clen < chunkMaxBytes {
			break
		}
	}
	return out.Bytes(), nil
}

// Close implements part of the Channel interface.
func (c *chunked) Close() error { return c.wc.Close() }
