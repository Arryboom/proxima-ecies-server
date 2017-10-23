# smartbox-ecies-server
A proof-of-concept go server for ECIES testing

Run `./app` and navigate to `http://localhost:2015`

## Encryption
GET localhost:2015/eciesEncrypt

#### Parameters:
  - *message*: message to be encrypted
  - *iv* (optional): if set, the server will also generate and return an initialization vector (iv)
  
  
#### Returns:
  - base64-encoded encrypted message
  - base64-encoded iv (if requested)
  - url-encoded version of the base64-encoded encrypted message
  - url-encoded version of the base64-encoded iv (if requested)
  
#### Example:
http://localhost:2015/eciesEncrypt?message=test&iv=1
```
Plain Base64:
Ciphertext: BBC6+URDzcz7M96/nRJqBPd7/QgSzATvBtNq8HVpQtBfrxF6aT8sW28KVudggI6iSphW0VJ3RaLYepJ3j7l5ff6pXeJ+GvBOAkDlLFjJr+4iBHP3jlb4KjRZFVrWYyxFGG/Z4v/K/m0aXDzgBK4S/B7F6ZTY
Shared Param (IV): hSbGzg5g7u1D/fLTK0yotw==


URL-encoded:
Ciphertext: BBC6%2BURDzcz7M96%2FnRJqBPd7%2FQgSzATvBtNq8HVpQtBfrxF6aT8sW28KVudggI6iSphW0VJ3RaLYepJ3j7l5ff6pXeJ%2BGvBOAkDlLFjJr%2B4iBHP3jlb4KjRZFVrWYyxFGG%2FZ4v%2FK%2Fm0aXDzgBK4S%2FB7F6ZTY
Shared Param (IV): hSbGzg5g7u1D%2FfLTK0yotw%3D%3D
```


## Decryption
GET localhost:2015/eciesDecrypt

#### Parameters:
  - *message*: the ciphertext to be decrypted
  - *iv* (optional): if set, the server will use this as the initialization vector (iv) for decryption
  
  
#### Returns:
  - plaintext decoded message

#### Example:
http://localhost:2015/eciesDecrypt?message=BBC6%2BURDzcz7M96%2FnRJqBPd7%2FQgSzATvBtNq8HVpQtBfrxF6aT8sW28KVudggI6iSphW0VJ3RaLYepJ3j7l5ff6pXeJ%2BGvBOAkDlLFjJr%2B4iBHP3jlb4KjRZFVrWYyxFGG%2FZ4v%2FK%2Fm0aXDzgBK4S%2FB7F6ZTY&iv=hSbGzg5g7u1D%2FfLTK0yotw%3D%3D
```
test
```
