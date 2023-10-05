package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InArrayInt64(needle int64, haystack []int64) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
