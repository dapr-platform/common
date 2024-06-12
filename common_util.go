package common

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/jaevor/go-nanoid"
)

var idgen func() string

func init() {
	var err error
	idgen, err = nanoid.Standard(21)
	if err != nil {
		panic(err)
	}
}
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
func GetSha1Hash(text string) string {
	hasher := sha1.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
func NanoId() string {
	return idgen()
}
