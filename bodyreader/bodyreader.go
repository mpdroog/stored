package bodyreader

import (
	"bytes"
	"fmt"
	"io"
)

var SEP = []byte("\r\n\r\n")

const (
	DefaultBufSize = 4096
)

type BodyReader struct {
	r    io.Reader
	done bool
	buf  []byte
	n    int
	pos  int
}

// Just keep on reading until we found that separator
func (d *BodyReader) Read(b []byte) (n int, err error) {
	max := len(b)

	if d.done {
		remain := d.n - d.pos
		if remain > 0 {
			if remain > max {
				remain = max
			}
			read := copy(b, d.buf[d.pos:d.pos+remain])
			d.pos += read
			return read, nil
		}

		return d.r.Read(b)
	}

	n, e := d.r.Read(d.buf[d.n:])
	d.n += n
	if len(d.buf)-10 <= d.n {
		// Grow
		d.buf = append(d.buf, make([]byte, DefaultBufSize)...)
	}

	if idx := bytes.Index(d.buf, SEP); idx != -1 {
		d.done = true
		from := idx + len(SEP)
		to := d.n - from
		if to > max {
			to = max
		}
		if len(d.buf) < d.n {
			panic(fmt.Errorf("Buf smaller than positions?? %d/%d (pos=%d)", len(d.buf), d.n, from))
		}

		read := copy(b, d.buf[from:from+to])
		d.pos = from + read
		return read, e
	}

	// No end of header yet so return nothing
	return 0, e
}

func New(r io.Reader, bufSize int) *BodyReader {
	return &BodyReader{r: r, buf: make([]byte, bufSize)}
}
