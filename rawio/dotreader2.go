package rawio

import (
	"io"
	"bufio"
)

// End Of Stream
var END = []byte("\r\n.\r\n")
//var END_SHORT = []byte(".\r\n")

type DotReader struct {
	r        *bufio.Reader
	done     bool
	state    int
}

func New(r *bufio.Reader, shortEnd bool) *DotReader {
	return &DotReader{r: r, done: false}
}

// Just keep on reading until we found that END.
func (d *DotReader) Read(b []byte) (n int, err error) {
	// Run data through a simple state machine to
  	// elide leading dots, rewrite trailing \r\n into \n,
  	// and detect ending .\r\n line.
  	const (
  		stateBeginLine = iota // beginning of line; initial state; must be zero
  		stateDot              // read . at beginning of line
  		stateDotCR            // read .\r at beginning of line
  		stateCR               // read \r (possibly at end of line)
  		stateData             // reading data in middle of line
  		stateEOF              // reached .\r\n end marker line
  	)
  	br := d.r
  	for n < len(b) && d.state != stateEOF {
  		var c byte
  		c, err = br.ReadByte()
  		if err != nil {
  			if err == io.EOF {
  				err = io.ErrUnexpectedEOF
  			}
  			break
  		}
		b[n] = c
  		n++

  		switch d.state {
  		case stateBeginLine:
  			if c == '.' {
  				d.state = stateDot
  				continue
  			}
  			if c == '\r' {
  				d.state = stateCR
  				continue
  			}
  			d.state = stateData
  
  		case stateDot:
  			if c == '\r' {
  				d.state = stateDotCR
  				continue
  			}
  			if c == '\n' {
  				d.state = stateEOF
  				continue
  			}
  			d.state = stateData
  
  		case stateDotCR:
  			if c == '\n' {
  				d.state = stateEOF
  				continue
  			}
  			// Not part of .\r\n.
  			// Consume leading dot and emit saved \r.
  			br.UnreadByte()
  			c = '\r'
  			d.state = stateData
  
  		case stateCR:
  			if c == '\n' {
  				d.state = stateBeginLine
  				break
  			}
  			// Not part of \r\n. Emit saved \r
  			br.UnreadByte()
  			c = '\r'
  			d.state = stateData
  
  		case stateData:
  			if c == '\r' {
  				d.state = stateCR
  				continue
  			}
  			if c == '\n' {
  				d.state = stateBeginLine
  			}
  		}
  	}
  	if err == nil && d.state == stateEOF {
  		err = io.EOF
  	}
  	/*if err != nil && d.r.dot == d {
  		d.r.dot = nil
  	}*/
  	return
}