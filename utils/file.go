package utils

import (
	"fmt"
	"jtrans/constants"
	"math"
	"strings"
	"unicode/utf8"
)

func ComputeChunkCount(size, chunkSize int64) int64 {
	chunkCount := size / chunkSize
	if size%chunkSize > 0 {
		chunkCount += 1
	}
	return chunkCount
}

func ComputeCurrentChunkSize(chunkNo, chunkCount, size int64) int64 {
	chunkSize := constants.ChunkSize
	if chunkNo == chunkCount {
		chunkSize = size - (chunkCount-1)*chunkSize
	}
	return chunkSize
}

func GetFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func FormatBytes(bytes float64) string {
	suffixes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	base := 1024.0

	if bytes < base {
		return fmt.Sprintf("%.0f %s", bytes, suffixes[0])
	}

	exp := int(math.Log(bytes) / math.Log(base))
	index := int(math.Min(float64(exp), float64(len(suffixes)-1)))

	return fmt.Sprintf("%.2f %s", bytes/math.Pow(base, float64(exp)), suffixes[index])
}

func Utf8Substr(str string, begin, size int) string {
	str = str[begin:]
	output := ""
	for i := 0; i < size && len(str) > 0; i++ {
		r, size := utf8.DecodeRuneInString(str)
		output += string(r)
		str = str[size:]
	}
	return output
}
