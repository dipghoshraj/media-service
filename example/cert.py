import ipaddress
from pathlib import Path
from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import rsa
import datetime

# File names
CERT_FILE = Path("agni.local.internal.pem")
KEY_FILE = Path("agni.local.internal-key.pem")

# Generate key
key = rsa.generate_private_key(
    public_exponent=65537,
    key_size=2048,
)

# Generate certificate
subject = issuer = x509.Name([
    x509.NameAttribute(NameOID.COMMON_NAME, "agni.local.internal")
])

cert = (
    x509.CertificateBuilder()
    .subject_name(subject)
    .issuer_name(issuer)
    .public_key(key.public_key())
    .serial_number(x509.random_serial_number())
    .not_valid_before(datetime.datetime.utcnow())
    .not_valid_after(datetime.datetime.utcnow() + datetime.timedelta(days=365))
    .add_extension(
        x509.SubjectAlternativeName([
            x509.DNSName("agni.local.internal"),
            x509.DNSName("localhost"),
            x509.IPAddress(ipaddress.IPv4Address("127.0.0.1")),  # ✅ Correct usage
        ]),
        critical=False,
    )
    .sign(key, hashes.SHA256())
)

# Write cert and key to files
with open(CERT_FILE, "wb") as f:
    f.write(cert.public_bytes(serialization.Encoding.PEM))

with open(KEY_FILE, "wb") as f:
    f.write(key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.TraditionalOpenSSL,
        encryption_algorithm=serialization.NoEncryption()
    ))

print(f"Generated cert: {CERT_FILE} and key: {KEY_FILE}")
