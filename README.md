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