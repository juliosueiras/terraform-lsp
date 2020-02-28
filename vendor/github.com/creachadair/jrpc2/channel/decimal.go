package channel

import (
	"bufio"
	"io"
	"strconv"
)

// Decimal is a framing that transmits and receives messages on r and wc, with
// each message prefixed by its length encoded as a line of decimal digits.
//
// For example, the message "empanada\n" is encoded as:
//
//    9\n
//    empanada\n
//
func Decimal(r io.Reader, wc io.WriteCloser) Channel {
	ch := newLenPrefix(r, wc)
	ch.enc = encodeLenDecimal
	ch.dec = decodeLenDecimal
	return ch
}

func encodeLenDecimal(n int, w io.Writer) error {
	var tmp [32]byte
	buf := strconv.AppendInt(tmp[:0], int64(n), 10)
	buf = append(buf, '\n')
	_, err := w.Write(buf)
	return err
}

func decodeLenDecimal(rd *bufio.Reader) (int, error) {
	pfx, err := rd.ReadString('\n')
	if err == io.EOF && pfx != "" {
		// handle a partial line at EOF
	} else if err != nil {
		return 0, err
	}
	return strconv.Atoi(pfx[:len(pfx)-1])
}
