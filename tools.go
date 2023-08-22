package GoUtils

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
)

func PrettyPrint(str []byte) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, str, "", " "); err != nil {
		return string(str), err
	}
	return strings.TrimSuffix(buf.String(), "\n"), nil
}

func GetCwdPath() string {
	dir, _ := os.Getwd()
	return dir
}

// Map2JsonString map转string
func Map2JsonString(param map[string]interface{}) string {
	data, _ := json.Marshal(param)
	return string(data)
}

func Json2Map(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}
