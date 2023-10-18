package encrypt

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func copyBlock(src []byte, srcOffset int, dst []uint32, dstOffset, count int) {
	dstBytes := make([]byte, 4)
	for i := 0; i < count; {
		dstIdx := (i + dstOffset) / 4
		dstVal := dst[dstIdx]
		byteIdx := (i + dstOffset) % 4
		binary.LittleEndian.PutUint32(dstBytes, dstVal)
		var j int
		for j = 0; j < 4-byteIdx && i+j < count; j++ {
			dstBytes[j+byteIdx] = src[i+j+srcOffset]
		}
		dst[dstIdx] = binary.LittleEndian.Uint32(dstBytes)
		i += j
	}
}

func MD5HashProc(md5 *MD5, data []byte) {
	md5.TransformBlock(data, 0, len(data))
}

func MD5HashProcFinish(md5 *MD5) string {
	res := md5.TransformFinalBlock()
	buf := bytes.NewBufferString("")
	for _, b := range res {
		buf.WriteString(fmt.Sprintf("%02x", b))
	}
	return buf.String()
}
