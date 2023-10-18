package encrypt

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256Hash(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
