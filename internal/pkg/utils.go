package pkg

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
