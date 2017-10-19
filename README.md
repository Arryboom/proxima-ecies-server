# pexeso-ecies-server
A go server for ECIES testing

Run ./app and navigate to localhost:2015

## To encrypt
GET localhost:2015/eciesEncrypt

parameters:
  - *message*: message to be encrypted
  - *iv* (optional): if set, the server will also generate and return an initialization vector (iv)
  
  
returns:
  - base64-encoded encrypted message
  - base64-encoded iv (if requested)
  - url-encoded version of the base64-encoded encrypted message
  - url-encoded version of the base64-encoded iv (if requested)


## To decrypt
GET localhost:2015/eciesDecrypt

parameters:
  - *message*: the ciphertext to be decrypted
  - *iv* (optional): if set, the server will use this as the initialization vector (iv) for decryption
  
  
returns:
  - plaintext decoded message
