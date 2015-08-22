#!/usr/bin/sh
EXPIRE=30
SHA= #-aes256

# generate $1.*
openssl genrsa $SHA -out $1.key.pem 4096
chmod og-rwx $1.key.pem
openssl req -new -sha256 -key $1.key.pem -out $1.csr
openssl x509 -req -days $EXPIRE -extfile /etc/ssl/openssl.cnf -extensions v3_req -CA intermediate.cert.pem -CAkey intermediate.key.pem -set_serial 01 -in $1.csr -out $1.cert.pem
#openssl x509 -noout -text -in $1.cert.pem

rm $1.csr
cat $1.cert.pem intermediate.cert.pem ca.cert.pem > $1.cert.bundle.pem
