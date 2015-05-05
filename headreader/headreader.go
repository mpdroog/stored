package headreader

import (
	"io"
	"bytes"
)

var SEP = []byte("\r\n\r\n")

type HeadReader struct {
	r io.Reader
	done bool
}

// Just keep on reading until we found that separator
func (d *HeadReader) Read(b []byte) (n int, err error) {
	if d.done {
		return 0, io.EOF
	}
	n, e := d.r.Read(b)
	if idx := bytes.Index(b, SEP); idx != -1 {
		d.done = true
		return idx, io.EOF
	}
	return n, e
}

func New(r io.Reader) *HeadReader {
	return &HeadReader{r: r}
}