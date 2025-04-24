package ecdh

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
)

// ECDH represents one party in the ECDH key exchange.
type ECDH struct {
	curve      ecdh.Curve
	privateKey *ecdh.PrivateKey
	publicKey  *ecdh.PublicKey
}

func NewX25519() (*ECDH, error) {
	return NewECDH(ecdh.X25519())
}

// NewECDH creates a new ECDH instance using the given curve.
func NewECDH(curve ecdh.Curve) (*ECDH, error) {
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECDH{
		curve:      curve,
		privateKey: priv,
		publicKey:  priv.PublicKey(),
	}, nil
}

// PublicKey returns this party's public key bytes.
func (e *ECDH) PublicKey() []byte {
	return e.publicKey.Bytes()
}

// SharedSecret calculates the shared secret using the other party's public key bytes.
func (e *ECDH) SharedSecret(peerPubBytes []byte) ([]byte, error) {
	peerPubKey, err := e.curve.NewPublicKey(peerPubBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid peer public key: %w", err)
	}

	secret, err := e.privateKey.ECDH(peerPubKey)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
