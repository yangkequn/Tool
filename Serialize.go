package Tool

import (
	"bytes"
	"encoding/gob"
	"strconv"
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
// JSTimeSequenceStringToArray convert JSTimeSequenceString to []int64
// unit of JS time is millisecond
// if value of JS time larger than 1 year, it's absolute time, otherwise it's timespan
func JSTimeSequenceStringToArray(timeString string) (unixTime []int64) {
	timeArray := StringSlit(timeString)
	unixTime = make([]int64, len(timeArray))
	for i, v := range timeArray {
		unixTime[i], _ = strconv.ParseInt(v, 10, 64)
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
