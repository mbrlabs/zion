package hodor

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
)

func GenerateRandomString(length int, alpabet string) string {
	alpabetLen := byte(len(alpabet))

	// make generate random byte array
	id := make([]byte, length)
	rand.Read(id)

	// replace rand num with char from alphabet
	for i, b := range id {
		id[i] = alpabet[b%alpabetLen]
	}

	return string(id)
}

func HashSha512(str string, salt string) string {
	hasher := sha512.New()
	hasher.Write([]byte(salt + str))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}
