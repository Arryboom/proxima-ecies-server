package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/O2-Czech-Republic/smartbox-ecies-go"
	"github.com/gorilla/mux"
)

func DecryptMessage(w http.ResponseWriter, req *http.Request) {

	message, err := base64.StdEncoding.DecodeString(req.URL.Query().Get("message"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}

	ownPrivateKey, err := LoadPrivateKey("./server.key.ec.pem")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}

	var sharedParam1 []byte

	if req.URL.Query().Get("iv") != "" {
		sharedParam1, _ = base64.StdEncoding.DecodeString(req.URL.Query().Get("iv"))
	} else {
		sharedParam1 = nil
	}

	ownPrivateKeyEC := ecies.ImportECDSA(ownPrivateKey)
	decrypted, err := ownPrivateKeyEC.Decrypt(bytes.NewReader(sharedParam1), message, nil, nil)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}

	if w != nil {
		w.Header().Set("Content-type", "text/plain")
		_, _ = w.Write([]byte(decrypted))
	}

}

func EncryptMessage(w http.ResponseWriter, req *http.Request) {

	var message string

	if req != nil {
		message = req.URL.Query().Get("message")
	} else {
		message = "Testing testing"
	}

	messageBytes := []byte(message)
	peerPublicKey, err := LoadPublicKey("./server.pub.pem")
	peerPublicKeyEC := ecies.ImportECDSAPublic(peerPublicKey)

	if err != nil {
		fmt.Println(err)
	}

	encrypted, err := ecies.Encrypt(rand.Reader, peerPublicKeyEC, messageBytes, nil, nil)
	baseSF := base64.StdEncoding.EncodeToString(encrypted)

	if w != nil {
		w.Header().Set("Content-type", "text/plain")
		_, _ = w.Write([]byte("Plain Base64:\n"))
		_, _ = w.Write([]byte("Ciphertext: " + baseSF + "\n"))
		//_, _ = w.Write([]byte("Shared Param (IV): " + base64.StdEncoding.EncodeToString(sharedParam1) + "\n"))
		_, _ = w.Write([]byte("\n\nURL-encoded:\n"))
		_, _ = w.Write([]byte("Ciphertext: " + url.QueryEscape(baseSF) + "\n"))
		//_, _ = w.Write([]byte("Shared Param (IV): " + url.QueryEscape(base64.StdEncoding.EncodeToString(sharedParam1)) + "\n"))
	}

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/eciesEncrypt", EncryptMessage).Methods("GET")
	router.HandleFunc("/eciesDecrypt", DecryptMessage).Methods("GET")
	log.Fatal(http.ListenAndServe(":2015", router))
}
