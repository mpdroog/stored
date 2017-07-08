package db

import (
	"testing"
)

func TestLookupPath(t *testing.T) {
	// lookupPath(msgHash [32]byte) string {
	h := hash("4K2dWjrk3y0Y3QvZ8xmcaQlT6qBRkJ81.206-261`@`JBinDown.local")
	res := lookupPath(h)

	if res != "d5d17901/0f2324c1/a9297fdc/ca9736a7/1e8b5c6e49cf35dfd9834ecc0fc4bb57.txt" {
		t.Errorf("Hash mismatch")
	}
}