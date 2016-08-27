// Copyright (c) 2016. See AUTHORS file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hodor

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"strings"
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

func escapeHTML(s string) string {
	new := strings.Replace(s, " ", "&nbsp;", -1)
	new = strings.Replace(new, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;", -1)
	new = strings.Replace(new, "\n", "<br>", -1)
	return new
}
