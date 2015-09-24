#!/usr/bin/sh
EXPIRE=400
SUBJ=C=US/ST=CA/O=Clouds
CADIR=${1:-/srv/cert/self}

# generate ca.key.pem and ca.cert.pem
#openssl genrsa -aes256 -passout file:ca.passwd -out ca.key.pem 4096
#chmod og-rwx ca.key.pem
#openssl req -new -x509 -sha256 -days $EXPIRE -extensions v3_ca -subj "/CN=Clouds CA/$SUBJ/" -passin file:ca.passwd -key ca.key.pem -out ca.cert.pem
#chmod o-rwx ca.cert.pem

if ! [ -f intermediate.passwd ]; then
  echo Generating intermediate.passwd
  cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1 > intermediate.passwd
  chmod 400 intermediate.passwd
  chown clouds:cert intermediate.passwd
fi

# generate intermediate.*
openssl genrsa -aes256 -passout file:intermediate.passwd -out intermediate.key.pem 4096
chmod og-rwx intermediate.key.pem
openssl req -new -sha256 -key intermediate.key.pem -subj "/CN=Clouds Internet Authority (Temp)/$SUBJ/" -passin file:intermediate.passwd -out intermediate.csr
openssl x509 -req -days $EXPIRE -extfile ext.cnf -extensions v3_ca -passin "file:$CADIR/ca.passwd" -CA "$CADIR/ca.cert.pem" -CAkey "$CADIR/ca.key.pem" -set_serial 01 -in intermediate.csr -out intermediate.cert.pem
chmod o-rwx intermediate.cert.pem intermediate.csr
chown clouds:cert intermediate.*.pem intermediate.csr

#rm intermediate.csr
