
set -e

CERT_DIR="./certs"
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"

mkdir -p "$CERT_DIR"

echo "Generating self-signed certificate..."

openssl req -x509 -nodes -days 365 \
  -newkey rsa:2048 \
  -keyout "$KEY_FILE" \
  -out "$CERT_FILE" \
  -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=localhost"

echo "Certificate generated:"
echo " - Certificate: $CERT_FILE"
echo " - Private Key: $KEY_FILE"
