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
