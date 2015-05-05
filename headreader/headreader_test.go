package headreader

import (
	"testing"
	"strings"
	"io"
	"bytes"
)

func TestHead(t *testing.T) {
	head := "Key: value\r\nKey2: value"
	body := "Body article\r\nblablabla"
	
	in := New(strings.NewReader(head + "\r\n\r\n" + body))
	buf := new(bytes.Buffer)

	if _, e := io.Copy(buf, in); e != nil {
		t.Error(e)
	}

	if buf.String() != head {
		t.Errorf("Invalid head received. Received=" + buf.String())
	}
}