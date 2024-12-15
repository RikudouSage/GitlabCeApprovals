package json

import (
	goJson "encoding/json"
	"io"
)

func ToJson(object any) (result []byte, err error) {
	result, err = goJson.Marshal(object)
	return
}

func ToJsonString(object any) (result string, err error) {
	raw, err := ToJson(object)
	if err != nil {
		return
	}

	result = string(raw)
	return
}

func ReadAsJson(reader io.Reader) ([]byte, error) {
	result, err := io.ReadAll(reader)
	return result, err
}

func Map[TOutput any](raw []byte, output *TOutput) error {
	return goJson.Unmarshal(raw, output)
}
