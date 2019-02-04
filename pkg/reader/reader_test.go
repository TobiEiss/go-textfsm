package reader_test

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/TobiEiss/go-textfsm/pkg/reader"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func TestReader(t *testing.T) {
	filepath := basepath + "/../../testfiles/01.txt"
	ch := make(chan string)
	go reader.ReadLineByLine(filepath, ch)

	for {
		line, ok := <-ch
		if !ok {
			break
		}
		log.Println(line)
	}
}
