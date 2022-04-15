package Tool

import (
	"github.com/cespare/xxhash/v2"
)

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
