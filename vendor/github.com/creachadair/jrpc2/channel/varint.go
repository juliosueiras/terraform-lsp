package channel

import (
	"bufio"
	"encoding/binary"
	"io"
)

// Varint is a framing that transmits and receives messages on r and wc, with
// each message prefixed by its length encoded in a varint as defined by the
// encoding/binary package.
func Varint(r io.Reader, wc io.WriteCloser) Channel {
	ch := newLenPrefix(r, wc)
	ch.enc = encodeLenVarint
	ch.dec = decodeLenVarint
	return ch
}

func encodeLenVarint(n int, w io.Writer) error {
	var ln [binary.MaxVarintLen64]byte
	nb := binary.PutUvarint(ln[:], uint64(n))
	_, err := w.Write(ln[:nb])
	return err
}

func decodeLenVarint(rd *bufio.Reader) (int, error) {
	ln, err := binary.ReadUvarint(rd)
	if err != nil {
		return 0, err
	}
	return int(ln), nil
}
