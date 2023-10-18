package tests

import (
	"bytes"
	"fmt"
	"github.com/spf13/cast"
	"jtrans/encrypt"
	"testing"
)

func TestMD5(t *testing.T) {
	md5 := encrypt.NewMD5()
	md5.TransformBlock([]byte{49, 50, 51}, 0, 3)
	res := md5.TransformFinalBlock()
	buf := bytes.NewBufferString("")
	for _, b := range res {
		buf.WriteString(cast.ToString(b))
	}
	fmt.Printf(buf.String())
}
