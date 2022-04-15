package Tool

import (
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"strings"
)

// Convert int64 value to string, using RawURLEncoding base64 encoding
func Int64ToString(i int64) string {
	b := []byte{byte(0xff & i), byte(0xff & (i >> 8)), byte(0xff & (i >> 16)), byte(0xff & (i >> 24)),
		byte(0xff & (i >> 32)), byte(0xff & (i >> 40)), byte(0xff & (i >> 48)), byte(0xff & (i >> 56))}
	sEnc := base64.RawURLEncoding.EncodeToString(b)
	return sEnc
}

//Convert string to int64, decode with RawURLEncoding base64 encoding
func StringToInt64(s string) int64 {
	if len(s) == 0 {
		return int64(0)
	}
	if len(s) > 13 {
		val, errint := strconv.ParseInt(s, 10, 64)
		if errint == nil {
			return val
		}
	}
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil || len(bytes) > 8 {
		return 0
	}
	value := binary.LittleEndian.Uint64(bytes)
	return int64(value)
}

func StringToFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

//split string to array, using "," as separator
func StringSlit(s string) []string {
	splitFn := func(c rune) bool {
		return c == ','
	}
	return strings.FieldsFunc(s, splitFn)
}
