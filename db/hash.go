package db

import (
	"stored/config"
	"encoding/binary"
	"crypto/md5"
	"fmt"
)

func hash(msgid string) [16]byte {
	return md5.Sum([]byte(msgid))
}

// Use MD5(msgid)[last 4bytes] to determine the biggest chance where we
// can find the file from all disks.
// MD5(msgid) -> last uint32 -> %diskcount -> start position in disks
func lookupDisk(msgHash [16]byte) int {
	offset := 1 // hardcoded
	begin := 16-offset-4

	pos := binary.BigEndian.Uint32(msgHash[begin:])
	modulo := pos % uint32(len(config.C.Storage)) +1 // TODO: remove +1?

	return int(modulo)
}

// Construct path to msgid if the file exists on a disk.
// Hex(MD5msgid) -> /dir/subdir/file.txt
func lookupPath(msgHash [16]byte) string {
	hash := fmt.Sprintf("%x", msgHash)
	return fmt.Sprintf("%s/%s/%s/%s.txt", hash[0:4], hash[5:8], hash[9:12], hash[13:16])
}