package bodyreader

import (
	"io"
	"bytes"
)

var SEP = []byte("\r\n\r\n")

type BodyReader struct {
	r io.Reader
	c bool
}

// Just keep on reading until we found that separator
func (d *BodyReader) Read(b []byte) (n int, err error) {
	n, e := d.r.Read(b)
	if d.c {
		// Found sep, forward it all
		return n, e
	}
	if idx := bytes.Index(b[:n], SEP); idx != -1 {
		// Found sep, strip header
		t := b[idx+len(SEP):]
		for i := 0; i < n-idx; i++ {
			b[i] = t[i]
		}

		d.c = true
		return n-idx-len(SEP), e
	}

	// No end of header yet so return nothing
	return 0, e
}

func New(r io.Reader) *BodyReader {
	return &BodyReader{r: r}
}