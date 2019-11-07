package channel

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

// Header defines a framing that transmits and receives messages using a header
// prefix similar to HTTP, in which mimeType describes the content type.
//
// Specifically, each message is sent in the format:
//
//    Content-Type: <mime-type>\r\n
//    Content-Length: <nbytes>\r\n
//    \r\n
//    <payload>
//
// The length (nbytes) is encoded as decimal digits. For example, given a
// mimeType value "application/json", the message "123\n" is transmitted as:
//
//    Content-Type: application/json\r\n
//    Content-Length: 4\r\n
//    \r\n
//    123\n
//
// If mimeType == "", the Content-Type header is omitted. Note, however, that
// the framing function returned by Header does not verify that the encoding of
// messages matches the declared mimeType.
func Header(mimeType string) Framing {
	return func(r io.Reader, wc io.WriteCloser) Channel {
		var ctype string
		if mimeType != "" {
			ctype = "Content-Type: " + mimeType + "\r\n"
		}
		return &hdr{
			mtype: mimeType,
			ctype: ctype,
			wc:    wc,
			rd:    bufio.NewReader(r),
			buf:   bytes.NewBuffer(nil),
		}
	}
}

// An hdr implements Channel. Messages sent on a hdr channel are framed as a
// header/body transaction, similar to HTTP.
type hdr struct {
	mtype string
	ctype string
	wc    io.WriteCloser
	rd    *bufio.Reader
	buf   *bytes.Buffer
	rbuf  []byte
}

// Send implements part of the Channel interface.
func (h *hdr) Send(msg []byte) error {
	h.buf.Reset()
	if h.ctype != "" {
		h.buf.WriteString(h.ctype)
	}
	h.buf.WriteString("Content-Length: ")
	h.buf.WriteString(strconv.Itoa(len(msg)))
	h.buf.WriteString("\r\n\r\n")
	h.buf.Write(msg)
	_, err := h.wc.Write(h.buf.Next(h.buf.Len()))
	return err
}

// Recv implements part of the Channel interface.
func (h *hdr) Recv() ([]byte, error) {
	var contentType, contentLength string
	for {
		raw, err := h.rd.ReadString('\n')
		if err == io.EOF && raw != "" {
			// handle a partial line at EOF
		} else if err != nil {
			return nil, err
		}
		if line := strings.TrimRight(raw, "\r\n"); line == "" {
			break
		} else if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
			// This implementation ignores unknown header fields.
			clean := strings.TrimSpace(parts[1])
			switch strings.ToLower(parts[0]) {
			case "content-type":
				contentType = clean
			case "content-length":
				contentLength = clean
			}
		} else {
			return nil, xerrors.New("invalid header line")
		}
	}

	// Verify that the content-type matches what we expect.
	if contentType != h.mtype {
		return nil, xerrors.New("invalid content-type")
	}

	// Parse out the required content-length field.
	if contentLength == "" {
		return nil, xerrors.New("missing required content-length")
	}
	size, err := strconv.Atoi(contentLength)
	if err != nil || size < 0 {
		return nil, xerrors.New("invalid content-length")
	}

	// We need to use ReadFull here because the buffered reader may not have a
	// big enough buffer to deliver the whole message, and will only issue a
	// single read to the underlying source.
	data := h.rbuf
	if len(data) < size || len(data) > (1<<20) && size < len(data)/4 {
		data = make([]byte, size*2)
		h.rbuf = data
	}
	if _, err := io.ReadFull(h.rd, data[:size]); err != nil {
		return nil, err
	}
	return data[:size], nil
}

// Close implements part of the Channel interface.
func (h *hdr) Close() error { return h.wc.Close() }

// LSP is a header framing (see Header) that transmits and receives messages on
// r and wc using the MIME type application/vscode-jsonrpc. This is the format
// preferred by the Language Server Protocol (LSP), defined by
// https://microsoft.github.io/language-server-protocol
var LSP = Header("application/vscode-jsonrpc; charset=utf-8")
