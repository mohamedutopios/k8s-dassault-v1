# Générer la clé CA
openssl genrsa -out ca.key 2048

# Créer le certificat CA
openssl req -x509 -new -nodes -key ca.key -subj "/CN=admission_ca" -days 10000 -out ca.crt

# Générer la clé du serveur
openssl genrsa -out server.key 2048

# Créer une demande de signature de certificat pour le serveur
openssl req -new -key server.key -subj "/CN=admission_webhook" -out server.csr

# Signer le certificat du serveur avec le CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 10000
