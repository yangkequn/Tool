package Tool

import (
	"bytes"
	"encoding/gob"
)

// return string encoded by Gob
func GobEncode(e interface{}) string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		return ""
	}
	return buf.String()
}

// return interface decoded by Gob
func GobDecode(str string, d interface{}) error {
	buf := bytes.NewBufferString(str)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(d)
	if err != nil {
		return err
	}
	return nil
}

var MaxTimeDiff int64 = 365 * 24 * 60 * 60 * 1000 //1 years = 365 days = 365 * 24 * 60 * 60 * 1000
func UnixTimeStringToArray_ms(timeString string) (unixTime []int64) {
	timeArray := StringSlit(timeString)
	unixTime = make([]int64, len(timeArray))
	for i, v := range timeArray {
		unixTime[i] = StringToInt64(v)
	}
	// adjust unixTime to unixTime_ms
	for i, v := range unixTime {
		if i == 0 {
			continue
		}
		if v < MaxTimeDiff {
			unixTime[i] = v + unixTime[i-1]
		}
	}
	return unixTime
}
