package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO this all might looks more cleaner tbh
func makeJsonResp(w http.ResponseWriter, httpStatus int, payload string) ([]byte, error) {
	// not really clean shadow headers mutate here but ok for our narrow prototype
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)

	if httpStatus != 200 {
		msg := QErrorOut{payload}
		serializedResp, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		return serializedResp, nil
	}

	msg := QSchemaOut{payload}
	serializedResp, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return serializedResp, nil
}

func handleJsonMarshallingErr(w http.ResponseWriter, err error) {
	http.Error(
		w,
		fmt.Sprintf("something really strange on json hardcoded marshalling: %s", err.Error()),
		http.StatusInternalServerError,
	)
}
