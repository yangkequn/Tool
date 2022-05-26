package Tool

import (
	"bytes"
	"strconv"
	"strings"
)

func DecodeAccelero(data string) string {
	var buffer bytes.Buffer
	for i, ie := 0, len(data); i < ie; i++ {
		v := data[i]
		skipnum := 0
		if v == 'v' {
			skipnum, _ = strconv.Atoi(data[i+1 : i+2])
			i += 1
		} else if v == 'w' {
			skipnum, _ = strconv.Atoi(data[i+1 : i+3])
			i += 2
		} else if v == 'x' {
			skipnum, _ = strconv.Atoi(data[i+1 : i+4])
			i += 3
		} else if v == 'y' {
			skipnum, _ = strconv.Atoi(data[i+1 : i+5])
			i += 4
		} else if v == 'z' {
			skipnum, _ = strconv.Atoi(data[i+1 : i+6])
			i += 5
		} else {
			buffer.WriteString(string(v))
		}
		for j := 0; j < skipnum; j++ {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}
func EncodeAccelero(data []int64) (r string) {
	if len(data) == 0 {
		return ""
	}

	var resultArray []string = make([]string, 0, len(data))
	var counter int = 0
	var lastValue int64 = data[0] + 1
	for i, v := range data {
		if v != lastValue {
			if i != 0 {
				if counter == 1 {
					resultArray = append(resultArray, ",")
				} else if counter < 10 {
					resultArray = append(resultArray, "v"+strconv.Itoa(counter))
				} else if counter < 100 {
					resultArray = append(resultArray, "w"+strconv.Itoa(counter))
				} else if counter < 1000 {
					resultArray = append(resultArray, "x"+strconv.Itoa(counter))
				} else if counter < 10000 {
					resultArray = append(resultArray, "y"+strconv.Itoa(counter))
				} else {
					resultArray = append(resultArray, "z"+strconv.Itoa(counter))
				}
			}
			resultArray = append(resultArray, strconv.FormatInt(v, 10))
			lastValue = v
			counter = 1
		} else {
			counter += 1
		}
	}

	if counter -= 1; counter > 0 {
		if counter == 1 {
			resultArray = append(resultArray, ",")
		} else if counter < 10 {
			resultArray = append(resultArray, "v"+strconv.Itoa(counter))
		} else if counter < 100 {
			resultArray = append(resultArray, "w"+strconv.Itoa(counter))
		} else if counter < 1000 {
			resultArray = append(resultArray, "x"+strconv.Itoa(counter))
		} else if counter < 10000 {
			resultArray = append(resultArray, "y"+strconv.Itoa(counter))
		} else {
			resultArray = append(resultArray, "z"+strconv.Itoa(counter))
		}
	}
	return strings.Join(resultArray, "")
}
