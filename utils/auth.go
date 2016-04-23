package utils

import (
	"strings"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

const (
	API_Secret   string = "66b1b350-2fbc-4997-aa35-ac551a31f413"
	API_KEY      string = "b69f7b3a-f752-43c7-a42c-2d3c299d9fc3"
	COMPANY_CODE string = "90014"
)

type Error struct {
	ErrorCode    string
	ErrorMessage struct {
		Indonesian string
		English    string
	}
}

func GetSignature(method, relative, accesstoken, body, timestamp string) string {
	// remove all space.
	re := regexp.MustCompile(" ")
	body = re.ReplaceAllString(body, "")

	// hash sha256.
	hasher := sha256.New()
	hasher.Write([]byte(body))
	hash := hasher.Sum(nil)

	// to lower
	lower := strings.ToLower(hex.EncodeToString(hash))

	tosign := method + ":" + relative + ":" + accesstoken + ":" + lower + ":" + timestamp

	// hmac-sha256.
	h := hmac.New(sha256.New, []byte(API_Secret))
	h.Write([]byte(tosign))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}
