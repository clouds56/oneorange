#!/usr/bin/sh
EXPIRE=30

# generate ca.key.pem and ca.cert.pem
openssl genrsa -aes256 -out ca.key.pem 4096
chmod og-rwx ca.key.pem
openssl req -new -x509 -sha256 -days $EXPIRE -extensions v3_ca -key ca.key.pem -out ca.cert.pem

# generate intermediate.*
openssl genrsa -aes256 -out intermediate.key.pem 4096
chmod og-rwx intermediate.key.pem
openssl req -new -sha256 -key intermediate.key.pem -out intermediate.csr
openssl x509 -req -days $EXPIRE -extfile /etc/ssl/openssl.cnf -extensions v3_ca -CA ca.cert.pem -CAkey ca.key.pem -set_serial 01 -in intermediate.csr -out intermediate.cert.pem

rm intermediate.csr
