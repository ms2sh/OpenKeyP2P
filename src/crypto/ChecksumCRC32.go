package crypto

import (
	"encoding/binary"
	"hash/crc32"
)

func ComputeChecksumCRC32(data []byte) []byte {
	cs := crc32.ChecksumIEEE(data)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, cs)
	return buf
}
