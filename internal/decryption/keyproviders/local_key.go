// Copyright 2022 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keyproviders

import (
	"encoding/base64"
	"os"

	"github.com/mongodb/mongocli/internal/decryption/aes"
)

// LocalKeyIdentifier config for the localKey used to encrypt the Log Encryption Key (LEK).
type LocalKeyIdentifier struct {
	KeyStoreIdentifier
	Filename string
}

// DecryptKey decrypts LEK using KMIP get or decrypt methods.
func (ki *LocalKeyIdentifier) DecryptKey(encryptedKey []byte) ([]byte, error) {
	encodedKEK, err := os.ReadFile(ki.Filename)
	if err != nil {
		return nil, err
	}

	kek, err := base64.StdEncoding.DecodeString(string(encodedKEK))
	if err != nil {
		return nil, err
	}

	iv := encryptedKey[:16]
	encryptedLEK := encryptedKey[16:48]

	cbc := &aes.CBCInput{
		Key: kek,
		IV:  iv,
	}
	return cbc.Decrypt(encryptedLEK)
}
