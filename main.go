package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func DecryptMessage(w http.ResponseWriter, req *http.Request) {

	message, err := base64.StdEncoding.DecodeString(req.URL.Query().Get("message"))
	sharedParam1, err := base64.StdEncoding.DecodeString(req.URL.Query().Get("iv"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}

	ownPrivateKey, err := LoadPrivateKey("./server.key.pem")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}

	ownPrivateKeyEC := ImportECDSA(ownPrivateKey)
	decrypted, err := ownPrivateKeyEC.Decrypt(rand.Reader, message, sharedParam1, nil)

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
	peerPrivateKey, err := LoadPrivateKey("./server.key.pem")

	if err != nil {
		fmt.Println(err)
	}

	var sharedParam1 []byte

	if req.URL.Query().Get("iv") != "" {
		sharedParam1 = make([]byte, 16)
		rand.Read(sharedParam1)
	} else {
		sharedParam1 = nil
	}

	peerPublicKey := ImportECDSAPublic(&peerPrivateKey.PublicKey)
	encrypted, err := Encrypt(rand.Reader, peerPublicKey, messageBytes, sharedParam1, nil)
	baseSF := base64.StdEncoding.EncodeToString(encrypted)

	if w != nil {
		w.Header().Set("Content-type", "text/plain")
		_, _ = w.Write([]byte("Plain Base64:\n"))
		_, _ = w.Write([]byte("Ciphertext: " + baseSF + "\n"))
		_, _ = w.Write([]byte("Shared Param (IV): " + base64.StdEncoding.EncodeToString(sharedParam1) + "\n"))
		_, _ = w.Write([]byte("\n\nURL-encoded:\n"))
		_, _ = w.Write([]byte("Ciphertext: " + url.QueryEscape(baseSF) + "\n"))
		_, _ = w.Write([]byte("Shared Param (IV): " + url.QueryEscape(base64.StdEncoding.EncodeToString(sharedParam1)) + "\n"))
	}

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/eciesEncrypt", EncryptMessage).Methods("GET")
	router.HandleFunc("/eciesDecrypt", DecryptMessage).Methods("GET")
	log.Fatal(http.ListenAndServe(":2015", router))
}
