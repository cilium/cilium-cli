// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package certs

import (
	"encoding/base64"
)

type CertManager struct {
	params Parameters

	caCert []byte
	caKey  []byte
}

type Parameters struct {
	Namespace string
}

func NewCertManager(p Parameters) *CertManager {
	return &CertManager{params: p}
}

// CAKeyBytes return the CA private certificate bytes, or nil when it is not
// set.
func (c *CertManager) CAKeyBytes() []byte {
	// NOTE: return a copy just to avoid the caller modifiying our CA
	// certificate.
	crt := make([]byte, len(c.caKey))
	copy(crt, c.caKey)
	return crt
}

// EncodeCertBytes returns an encoded format of the certificate bytes.
func EncodeCertBytes(cert []byte) string {
	return base64.StdEncoding.EncodeToString(cert)
}
