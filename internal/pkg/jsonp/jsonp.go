package jsonp

import (
	"encoding/json"
	"io"
)

func ToStruct[Struct any](j interface{}) Struct {
	var s Struct

	ioj, isstream := j.(io.Reader)
	if isstream == true {
		json.NewDecoder(ioj).Decode(&s)
		return s
	}

	strj, _ := j.(string)
	json.Unmarshal([]byte(strj), &s)

	return s
}

func ToJson(s interface{}) string {
	j, _ := json.Marshal(s)

	return string(j)
}
