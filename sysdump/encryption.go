// Copyright 2021 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sysdump

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

//go:embed pubkey.gpg
var rootPubKey []byte

func createEncryptedZipFile(pathToZipFile string, pathsToEncryptionKey []string) (string, error) {
	// Read the destination entity.
	var entityList []*openpgp.Entity
	for _, k := range pathsToEncryptionKey {
		e, err := readDestinationEntity(k)
		if err != nil {
			return "", fmt.Errorf("failed to read encryption key %q: %w", k, err)
		}
		entityList = append(entityList, e)
	}
	// Create the destination file.
	d := pathToZipFile + ".gpg"
	o, err := os.Create(d)
	if err != nil {
		return "", err
	}
	// Encrypt the zip file.
	h := &openpgp.FileHints{
		FileName: filepath.Base(pathToZipFile),
		IsBinary: true,
		ModTime:  time.Now(),
	}
	encOut, err := openpgp.Encrypt(o, entityList, nil, h, nil)
	if err != nil {
		return "", err
	}
	defer encOut.Close()
	i, err := os.Open(pathToZipFile)
	if err != nil {
		return "", err
	}
	r, err := io.ReadAll(i)
	if err != nil {
		return "", err
	}
	if _, err = encOut.Write(r); err != nil {
		return "", err
	}
	// Return the path to the encrypted file.
	return d, nil
}

func readDestinationEntity(pathToEncryptionKey string) (*openpgp.Entity, error) {
	var r io.Reader
	switch pathToEncryptionKey {
	case DefaultEncryptionKey:
		r = bytes.NewReader(rootPubKey)
	default:
		k, err := os.Open(pathToEncryptionKey)
		if err != nil {
			return nil, err
		}
		r = k
	}
	return openpgp.ReadEntity(packet.NewReader(r))
}
