package ecdh

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"fmt"
)

type ECDH struct {
	curve        ecdh.Curve
	privateKey   *ecdh.PrivateKey
	publicKey    *ecdh.PublicKey
	sharedSecret []byte
}

func NewX25519() (*ECDH, error) {
	return NewECDH(ecdh.X25519())
}

func NewECDH(curve ecdh.Curve) (*ECDH, error) {
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECDH{
		curve:        curve,
		privateKey:   priv,
		publicKey:    priv.PublicKey(),
		sharedSecret: nil,
	}, nil
}

func (e *ECDH) PublicKey() []byte {
	return e.publicKey.Bytes()
}

func (e *ECDH) SharedSecret() ([]byte, error) {
	if e.sharedSecret == nil {
		return nil, errors.New("shared secret was not yet generated")
	}
	return e.sharedSecret, nil
}

func (e *ECDH) GenerateSharedSecret(peerPubBytes []byte) error {
	if e.sharedSecret != nil {
		return errors.New("shared secret was already generated")
	}

	peerPubKey, err := e.curve.NewPublicKey(peerPubBytes)
	if err != nil {
		return fmt.Errorf("invalid peer public key: %w", err)
	}

	secret, err := e.privateKey.ECDH(peerPubKey)
	if err != nil {
		return err
	}

	e.sharedSecret = secret

	return nil
}
