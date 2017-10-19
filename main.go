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
	sharedParam1, err := base64.StdEncoding.DecodeString(req.URL.Query().Get("s1"))
	sharedParam2, err := base64.StdEncoding.DecodeString(req.URL.Query().Get("s2"))

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
	decrypted, err := ownPrivateKeyEC.Decrypt(rand.Reader, message, sharedParam1, sharedParam2)

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

	sharedParam1 := make([]byte, 16)
	sharedParam2 := make([]byte, 16)

	rand.Read(sharedParam1)
	rand.Read(sharedParam2)

	peerPublicKey := ImportECDSAPublic(&peerPrivateKey.PublicKey)
	encrypted, err := Encrypt(rand.Reader, peerPublicKey, messageBytes, sharedParam1, sharedParam2)
	baseSF := base64.StdEncoding.EncodeToString(encrypted)

	if w != nil {
		w.Header().Set("Content-type", "text/plain")
		_, _ = w.Write([]byte("Plain Base64:\n"))
		_, _ = w.Write([]byte("Ciphertext: " + baseSF + "\n"))
		_, _ = w.Write([]byte("Shared Param 1: " + base64.StdEncoding.EncodeToString(sharedParam1) + "\n"))
		_, _ = w.Write([]byte("Shared Param 2: " + base64.StdEncoding.EncodeToString(sharedParam2) + "\n"))
		_, _ = w.Write([]byte("\n\nURL-encoded:\n"))
		_, _ = w.Write([]byte("Ciphertext: " + url.QueryEscape(baseSF) + "\n"))
		_, _ = w.Write([]byte("Shared Param 1: " + url.QueryEscape(base64.StdEncoding.EncodeToString(sharedParam1)) + "\n"))
		_, _ = w.Write([]byte("Shared Param 2: " + url.QueryEscape(base64.StdEncoding.EncodeToString(sharedParam2)) + "\n"))
	}

	privateKey := ImportECDSA(peerPrivateKey)
	decrypted, _ := privateKey.Decrypt(rand.Reader, encrypted, nil, nil)
	fmt.Println(decrypted)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/eciesEncrypt", EncryptMessage).Methods("GET")
	router.HandleFunc("/eciesDecrypt", DecryptMessage).Methods("GET")
	log.Fatal(http.ListenAndServe(":2015", router))
}
