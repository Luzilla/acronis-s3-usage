package ostor

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

// this is a very confusing thing to do, especially because the "query string" used always
// needs to use the first bit only — ostor-usage, ostor-users, ostor-limits etc.
func createSignature(httpMethod, awsSecretKey, query string) (signature string, date string, err error) {
	date = time.Now().Format(time.RFC1123Z)
	data := fmt.Sprintf("%s\n\n\n%s\n/%s", httpMethod, date, "?"+query)

	h := hmac.New(sha1.New, []byte(awsSecretKey))
	_, err = h.Write([]byte(data))
	if err != nil {
		return
	}

	signature = base64.StdEncoding.EncodeToString(h.Sum(nil))

	// fmt.Println("data      : " + data)
	// fmt.Println("secret    : " + awsSecretKey)
	// fmt.Println("date      : " + date)
	// fmt.Println("query     : " + queryString)
	// fmt.Println("signature : " + signature)
	return
}

func authHeader(keyID, signature string) string {
	return fmt.Sprintf("AWS %s:%s", keyID, signature)
}
