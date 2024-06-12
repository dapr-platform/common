package common

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strings"
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

type AccessInfo struct {
	Aud string `json:"aud"` //client_id
	Exp int    `json:"exp"` //expired time
	Sub string `json:"sub"` //user_name
}

// GetIP returns request real ip.
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}
	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "0.0.0.0"
	}
	if net.ParseIP(ip) != nil {
		return ip
	}
	return "0.0.0.0"
}

func ExtractJwtFromMap(m map[string]string) (*AccessInfo, error) {
	authVal, exist := m["Authorization"]
	if !exist {
		_, exist = m["authorization"]
		if !exist {
			return nil, errors.New("no authorization header")
		}
		authVal = m["authorization"]
	}
	return ExtractJwtFromString(authVal)
}
func ExtractJwtFromString(authVal string) (*AccessInfo, error) {

	as := authVal
	kv := strings.Split(as, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		return nil, errors.New("error authorization header")
	}
	access := kv[1]
	kv = strings.Split(access, ".")
	if len(kv) != 3 {
		return nil, errors.New("error access token, is not jwt")
	}
	b, err := base64.RawURLEncoding.DecodeString(kv[1])
	if err != nil {
		return nil, err
	}
	var info AccessInfo
	err = json.Unmarshal(b, &info)
	return &info, err
}
func ExtractJwt(r *http.Request) (*AccessInfo, error) {
	as := r.Header.Get("Authorization")
	if as == "" {
		as = r.Header.Get("authorization")
	}
	if as == "" {
		return nil, errors.New("can't find key authorization or Authorization")
	}
	return ExtractJwtFromString(as)
}

func ExtractUserSub(r *http.Request) (sub string, err error) {
	if val := r.Header.Get("X-User-Id"); val != "" { //can add by apisix, using for swagger
		return val, nil
	} else {
		accInfo, err := ExtractJwt(r)
		if err != nil {
			err = errors.Wrap(err, "extract jwt error")
			return "", err
		} else {
			sub = accInfo.Sub
			return sub, nil
		}

	}
}
