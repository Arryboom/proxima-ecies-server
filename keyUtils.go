package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func LoadPrivateKey(path string) (*ecdsa.PrivateKey, error) {

	serverPrivateKey, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	privateKeyPEM, err := GetPEMBlock(serverPrivateKey, "EC PRIVATE KEY")

	if err != nil {
		return nil, err
	}

	return x509.ParseECPrivateKey(privateKeyPEM)

}

func GetPEMBlock(dat []byte, t string) ([]byte, error) {

	var block *pem.Block

	for len(dat) > 0 {

		block, dat = pem.Decode(dat)

		if block.Type == t {
			return block.Bytes, nil
		}

	}

	return nil, fmt.Errorf("PEM block %s not found in the provided data", t)

}
