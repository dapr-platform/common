package common

import (
	"encoding/json"
	"net/http"
)

func HttpError(w http.ResponseWriter, response *Response, code int) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response.ToBytes())
}

func HttpResult(w http.ResponseWriter, response *Response) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if response != nil {
		w.Write(response.ToBytes())
	} else {
		w.Write(OK.ToBytes())
	}
}

func HttpSuccess(w http.ResponseWriter, response *Response) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if response != nil {
		w.Write(response.ToBytes())
	} else {
		w.Write(OK.ToBytes())
	}
}

func ReadRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}
func ReadResponseBody(r *http.Response, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func MarshalWithRemoveKey(v interface{}, key string) ([]byte, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var tmp map[string]interface{}
	err = json.Unmarshal(buf, &tmp)
	if err != nil {
		return nil, err
	}
	delete(tmp, "key")
	return json.Marshal(tmp)
}
