// sha1
package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetSha1Hash(content string) string {
	h := sha1.New()

	h.Write([]byte(content))

	return hex.EncodeToString(h.Sum(nil))
}
