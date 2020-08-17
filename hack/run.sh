#!/bin/bash

set -o xtrace
set -e

# generate certs
FILE=ca.crt
if [ -f "$FILE" ]; then
  echo "skipping cert gen"
else
  openssl req -nodes -new -x509 -keyout ca.key -out ca.crt
  openssl req -out client.csr -new -newkey rsa:4096 -nodes -keyout client.key -subj "/CN=development/O=system:masters"
  openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt
  openssl pkcs12 -export -in ./client.crt -inkey ./client.key -out client.p12 -passout pass:password
fi

# generate config file
FILE=kubeconfig
if [ -f "$FILE" ]; then
  echo "skipping kubecoonfig gen"
else
echo "apiVersion: v1
clusters:
- cluster:
    certificate-authority: ./apiserver.local.config/certificates/apiserver.crt
    server: https://localhost:8443
  name: apiserver
contexts:
- context:
    cluster: apiserver
    user: apiserver
  name: apiserver
current-context: apiserver
kind: Config
preferences: {}
users:
- name: apiserver
  user:
    client-certificate: ./apiserver.local.config/certificates/apiserver.crt
    client-key: ./apiserver.local.config/certificates/apiserver.key
    tokenFile: ./token
" > kubeconfig
fi

go run . --secure-port 8443 --etcd-servers http://127.0.0.1:2379 --v=7 \
 --client-ca-file ca.crt \
 --kubeconfig ~/.kube/config \
 --authentication-kubeconfig ~/.kube/config \
 --authorization-kubeconfig ~/.kube/config --write-token-file

