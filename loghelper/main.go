package loghelper

import (
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"strings"
)

func DumpLog(res interface{}) {
	result := spew.Sdump(res)
	strSlice := strings.Split(result, "\n")
	for _, s := range strSlice {
		log.Debug(s)
	}
}
