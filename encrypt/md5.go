package encrypt

import (
	"bytes"
	md5 "crypto/md5"
	"encoding/binary"
	"encoding/hex"
)

func MD5Hash(data []byte) string {
	md := md5.New()
	md.Write(data)
	return hex.EncodeToString(md.Sum(nil))
}

var sineTable = []uint32{
	0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee, 0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
	0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be, 0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
	0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa, 0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
	0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed, 0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
	0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c, 0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
	0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05, 0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
	0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039, 0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
	0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1, 0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
}

var shiftTable = [4][4]int{
	{7, 12, 17, 22}, {5, 9, 14, 20}, {4, 11, 16, 23}, {6, 10, 15, 21},
}

var terminator = []byte{0x80}
var zeroBytes = [56]byte{}

type MD5 struct {
	mState     [4]uint32
	mBuffer    [16]uint32
	mBitCount  int64
	mBufferPos int
}

type MD5StateStorage struct {
	MdBuffer    [4]uint32  `json:"md_buffer"`
	BitCount    int64      `json:"bit_count"`
	BufCount    int        `json:"buf_count"`
	InputBuffer [16]uint32 `json:"input_buffer"`
}

func NewMD5() *MD5 {
	return &MD5{
		mState: [4]uint32{0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476},
	}
}

func NewMD5FromStorage(storage *MD5StateStorage) *MD5 {
	return &MD5{
		mState:     storage.MdBuffer,
		mBitCount:  storage.BitCount,
		mBufferPos: storage.BufCount,
		mBuffer:    storage.InputBuffer,
	}
}

func (md5 *MD5) init() {
	md5.mState = [4]uint32{0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476}
	md5.mBitCount = 0
	md5.mBufferPos = 0
}

func (md5 *MD5) rotL(v uint32, count int) uint32 {
	count &= 0x1F
	return (v << count) | (v >> (32 - count))
}

func (md5 *MD5) transform() {
	a := md5.mState[0]
	b := md5.mState[1]
	c := md5.mState[2]
	d := md5.mState[3]

	for i := 0; i < 64; i++ {
		var f uint32
		var g int
		if i < 16 {
			f = d ^ b&(c^d)
			g = i
		} else if i < 32 {
			f = c ^ d&(b^c)
			g = (5*i + 1) & 0xF
		} else if i < 48 {
			f = b ^ c ^ d
			g = (3*i + 5) & 0xF
		} else {
			f = c ^ (b | ^d)
			g = (7 * i) & 0xF
		}
		t := d
		d = c
		c = b
		b += md5.rotL(a+f+md5.mBuffer[g]+sineTable[i], shiftTable[i>>4][i&3])
		a = t
	}

	md5.mState[0] += a
	md5.mState[1] += b
	md5.mState[2] += c
	md5.mState[3] += d
}

func (md5 *MD5) TransformBlock(data []byte, pos, count int) {
	md5.mBitCount += int64(count << 3)
	if md5.mBufferPos != 0 {
		bufCount := 64 - md5.mBufferPos
		if count < bufCount {
			copyBlock(data, pos, md5.mBuffer[:], md5.mBufferPos, count)
			md5.mBufferPos += count
			return
		}
		copyBlock(data, pos, md5.mBuffer[:], md5.mBufferPos, bufCount)
		md5.transform()
		pos += bufCount
		count -= bufCount
		md5.mBufferPos = 0
	}
	for count >= 64 {
		copyBlock(data, pos, md5.mBuffer[:], 0, 64)
		md5.transform()
		pos += 64
		count -= 64
	}
	if count > 0 {
		copyBlock(data, pos, md5.mBuffer[:], 0, count)
		md5.mBufferPos += count
	}
}

func (md5 *MD5) TransformFinalBlock() []byte {
	copyBlock(terminator, 0, md5.mBuffer[:], md5.mBufferPos, 1)
	md5.mBufferPos++
	bufCount := 64 - md5.mBufferPos
	if bufCount < 8 {
		copyBlock(zeroBytes[:], 0, md5.mBuffer[:], md5.mBufferPos, bufCount)
		md5.transform()
		md5.mBufferPos = 0
		bufCount = 64
	}

	copyBlock(zeroBytes[:], 0, md5.mBuffer[:], md5.mBufferPos, bufCount-8)
	md5.mBuffer[14] = uint32(md5.mBitCount)
	md5.mBuffer[15] = uint32(md5.mBitCount >> 32)
	md5.transform()
	md5.mBufferPos = 0

	var buf bytes.Buffer
	for _, val := range md5.mState {
		_ = binary.Write(&buf, binary.LittleEndian, val)
	}

	return buf.Bytes()
}

func (md5 *MD5) GetValue() *MD5StateStorage {
	return &MD5StateStorage{
		MdBuffer:    md5.mState,
		BitCount:    md5.mBitCount,
		BufCount:    md5.mBufferPos,
		InputBuffer: md5.mBuffer,
	}
}
