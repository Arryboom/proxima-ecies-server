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

func LoadPublicKey(path string) (*ecdsa.PublicKey, error) {

	publicKey, err := ioutil.ReadFile(path)

	block, _ := pem.Decode([]byte(publicKey))

	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}

	return pub.(*ecdsa.PublicKey), nil

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
