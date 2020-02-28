package channel

import (
	"bufio"
	"bytes"
	"io"
)

func newLenPrefix(r io.Reader, wc io.WriteCloser) *lenprefix {
	return &lenprefix{
		wc:  wc,
		rd:  bufio.NewReader(r),
		buf: bytes.NewBuffer(nil),
	}
}

// A lenprefix implements Channel for length-prefixed framings.
type lenprefix struct {
	wc  io.WriteCloser
	rd  *bufio.Reader
	buf *bytes.Buffer

	// Encode a length into the provided buffer.
	enc func(int, io.Writer) error

	// Decode a length from the provided reader.
	dec func(*bufio.Reader) (int, error)
}

// Send implements part of the Channel interface. It encodes len(msg) using the
// encoding function, concatenates it with the message body, and writes the
// message to the underlying writer.
func (c *lenprefix) Send(msg []byte) error {
	c.buf.Reset()
	if err := c.enc(len(msg), c.buf); err != nil {
		return err
	}
	c.buf.Write(msg)
	_, err := c.wc.Write(c.buf.Next(c.buf.Len()))
	return err
}

// Recv implements part of the Channel interface. It decodes a message length
// using the decoding function, then reads that many bytes from the underlying
// reader.
func (c *lenprefix) Recv() ([]byte, error) {
	ln, err := c.dec(c.rd)
	if err != nil {
		return nil, err
	}
	out := make([]byte, int(ln))
	nr, err := io.ReadFull(c.rd, out)
	return out[:nr], err
}

// Close implements part of the Channel interface.
func (c *lenprefix) Close() error { return c.wc.Close() }
