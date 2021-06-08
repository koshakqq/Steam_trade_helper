package dmarket

import (
	"crypto/ed25519"
	"encoding/hex"
)



func GetPrivateKey(s string) (*[64]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var privateKey [64]byte
	copy(privateKey[:], b[:64])

	return &privateKey, nil
}

func sign(pk *[64]byte, msg []byte) string {
	return hex.EncodeToString(ed25519.Sign((*pk)[:], msg))
}

func (k Keys)Sign(msg string) (signature string, err error) {
	b, err := GetPrivateKey(k.Private)
	return sign(b, []byte(msg)), nil
}
