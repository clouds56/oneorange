#!/usr/bin/sh
EXPIRE=30
SHA= #-aes256
SUBJ=C=US/ST=CA/O=Orangez
OUT=${2:-$1}

# generate $1.*
openssl genrsa $SHA -out "$OUT.key.pem" 4096
chmod og-rwx "$OUT.key.pem"
openssl req -new -sha256 -key "$OUT.key.pem" -subj "/CN=$1/$SUBJ/" -out "$OUT.csr"
openssl x509 -req -days $EXPIRE -extfile ext.cnf -extensions v3_req -passin file:intermediate.passwd -CA intermediate.cert.pem -CAkey intermediate.key.pem -set_serial 01 -in "$OUT.csr" -out "$OUT.cert.pem"
chmod o-rwx "$OUT.cert.pem" "$OUT.csr"
#openssl x509 -noout -text -in "$OUT.cert.pem"

cat "$OUT.cert.pem" intermediate.cert.pem /srv/cert/self/ca.cert.pem > "$OUT.cert.bundle.pem"
chmod o-rwx "$OUT.cert.bundle.pem"
chown clouds:cert "$OUT".*.pem "$OUT".csr
#rm "$OUT.csr"
