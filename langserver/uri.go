package langserver

import (
	"github.com/go-language-server/uri"
)

func absolutePath(in string) (uri.URI, error) {
	return uri.Parse(string(in))
}
