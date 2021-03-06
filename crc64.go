package Tool

import (
	"math/rand"

	"github.com/cespare/xxhash/v2"
)

//Random stringID ,base64 encoding
func RandomBase64ID() string {
	id := int64(rand.Uint64())
	return Int64ToString(id)
}

//Hash to string
func HashToString(s string) string {
	v := int64(xxhash.Sum64String(s))
	return Int64ToString(v)
}
func HashToInt64(s string) int64 {
	return int64(xxhash.Sum64String(s))
}
func Sum64(b []byte) int64 {
	return int64(xxhash.Sum64(b))
}
