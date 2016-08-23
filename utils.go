package hodor

import (
	"crypto/rand"
)

const (
	alphabetAlphaNum     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphabetAlphaNumPlus = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$%&()*+-./<=>?@[]^_{|}~"
)

func generateRandomString(length int, alpabet string) string {
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
