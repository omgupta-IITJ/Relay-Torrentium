Features of relay libp2p node

protocol using => Circuit relay protocol
Features of Relay node

**NOISE PROTOCOL**
When two libp2p peers connect:

> libp2p checks which security protocols both support (Noise, TLS, etc.).
> They agree to use Noise.
> noise.New runs → performs a Noise handshake (like the XX pattern by default).
> They exchange keys → derive a shared session key.
> From now on, all traffic is encrypted & authenticated using that key.

**Why Noise here?**
> It’s lightweight compared to TLS.
> Provides forward secrecy and identity protection.
> Standard in many peer-to-peer systems (WireGuard(vpn SERVICE), libp2p, etc.).

**WHY USING WEBSOCKET TRANSPORT LAYER?**
> 

**CIRCUIT RELAY PROTOCOL**
>

**RSA as public-key cryptography algorithm**
To generate an RSA PEM(privacy-enhanced mail) key pair, we need to follow these steps:
> Generate a new RSA private key.
> Encode the private key to the PEM format.
> Extract the public key from the private key.
> Encode the public key to the PEM format.

**Main point:** libp2p uses its own crypto insterface hence we need to convert RSA private key to libp2p format
rsa.PrivateKey = pure cryptographic key.
crypto.PrivKey = same key, but in a wrapper that libp2p understands and can operate with.

**Marshal** = “marshal” means to convert a data structure into a format that can be stored or transmitted.

Marshal → DER → crypto.UnmarshalRsaPrivateKey is necessary to convert your standard RSA key into libp2p’s crypto.PrivKey

**🔑 1. RSA**

Type: Asymmetric crypto based on factoring large prime numbers.
Key Size: Typically 2048–4096 bits.
Operations:
Encryption/decryption
Digital signatures
Performance:
Relatively slow (big numbers, modular arithmetic).
Larger keys → more memory + network overhead.
Security level: 2048-bit RSA ≈ 112-bit security (modern minimum).
Libp2p impact: Generates Peer IDs starting with Qm... (legacy multihash style).

**2. Ed25519**

Type: Elliptic Curve Cryptography (ECC), based on Curve25519.
Key Size: 256 bits (much smaller).
Operations:
Only for signatures (not direct encryption).
Very fast signing + verification.
Performance:
Extremely fast and compact.
Tiny keys (32-byte private, 32-byte public).
Security level: ≈ 128-bit security (stronger than RSA-2048).
Libp2p impact: Generates Peer IDs like 12D3Koo... (default, modern format).

**using .bin file for ed25519_key**