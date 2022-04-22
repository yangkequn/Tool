package Tool

import (
	"strconv"
	"strings"
)

func StringsJoint(stringArray []string) string {
	return strings.Join(stringArray, ",")
}
func StringToInt64Array(s string) []int64 {
	var items []string = StringSlit(s)
	return StringArrayToInt64Array(items)
}
func StringToFloat32Array(s string) []float32 {
	var items []string = StringSlit(s)
	return StringArrayToFloat32Array(items)
}
func Float32ArrayToString(array []float32) (result string) {
	if len(array) == 0 {
		return ""
	}
	for _, v := range array {
		// parameter -1 means no trailing zeros
		result += strconv.FormatFloat(float64(v), 'f', -1, 32) + ","
	}
	return result[:len(result)-1]
}
func Float64ArrayToString(array []float64) (result string) {
	if len(array) == 0 {
		return ""
	}
	for _, v := range array {
		result += strconv.FormatFloat(float64(v), 'f', -1, 64) + ","
	}
	return result[:len(result)-1]
}

//convert []string to []float32
func StringArrayToFloat32Array(array []string) (result []float32) {
	result = make([]float32, len(array))
	for _, v := range array {
		f32, err := StringToFloat32(v)
		if err != nil {
			return result
		}
		result = append(result, f32)
	}
	return result
}

//convert []string to []int64
func StringArrayToInt64Array(array []string) []int64 {
	var result []int64
	for _, v := range array {
		result = append(result, StringToInt64(v))
	}
	return result
}

// remove string item from string array, return true if item is removed
func RemoveItemFromStringArray(array *[]string, item string) bool {
	l, lOld := len(*array), len(*array)
	for i := 0; i < l; i++ {
		if (*array)[i] == item {
			for j := i; j < l-1; j++ {
				(*array)[j] = (*array)[j+1]
			}
			i--
		}
	}
	return lOld != l
}

//convert []int64 to []string
func Int64ArrayToBasr64StringArray(array []int64) []string {
	var result []string
	for _, v := range array {
		result = append(result, Int64ToString(v))
	}
	return result
}
func Int64ArrayToBase64String(unixTime []int64) string {
	var result []string = Int64ArrayToBasr64StringArray(unixTime)
	return StringsJoint(result)
}
func Int64ArrayToBase10String(unixTime []int64) (result string) {
	for _, v := range unixTime {
		result += strconv.FormatInt(v, 10) + ","
	}
	return result[:len(result)-1]
}
func Base10StringToInt64Array(s string) (result []int64) {
	var items []string = StringSlit(s)
	for _, v := range items {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return result
		}
		result = append(result, i)
	}
	return result
}

// return index of item in string array, -1 if not found
func IndexOfStringArray(array []string, item string) int {
	for i, v := range array {
		if v == item {
			return i
		}
	}
	return -1
}

// return index of item in int64 array, -1 if not found
func IndexOfInt64Array(array []int64, item int64) int {
	for i, v := range array {
		if v == item {
			return i
		}
	}
	return -1
}
