package utils

func GetSignature(method, relative, accesstoken, body, timestamp string) string {

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
