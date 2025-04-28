package hasherService

import (
	"crypto/ed25519"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/matthewhartstonge/argon2"
	"github.com/mr-tron/base58"
)

func HashPassword(password string) (hashedPassword string, err error) {
	argon := argon2.DefaultConfig()

	hash, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHashAndPassword(hashedPassword string, password string) (ok bool, err error) {
	return argon2.VerifyEncoded([]byte(password), []byte(hashedPassword))
}

// MD5Hash takes a string and returns its MD5 hash as a hexadecimal string.
func MD5Hash(msg string) string {
	hash := md5.New()

	hash.Write([]byte(msg))

	return hex.EncodeToString(hash.Sum(nil))
}

// VerifySignature verifies a Solana signature
func VerifySignature(pubKeyBase58 string, message string, signatureBase58 string) (bool, error) {
	// Public key'i çözümle
	publicKey, err := solana.PublicKeyFromBase58(pubKeyBase58)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %v", err)
	}

	// Signature'ı çözümle (Base58'den byte dizisine dönüştürme)
	signature, err := base58.Decode(signatureBase58)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %v", err)
	}

	// Signature uzunluğunu kontrol et
	if len(signature) != 64 {
		return false, fmt.Errorf("invalid signature length: got %d, want 64", len(signature))
	}

	// Public key uzunluğunu kontrol et
	if len(publicKey) != 32 {
		return false, fmt.Errorf("invalid public key length: got %d, want 32", len(publicKey))
	}

	// İmzanın geçerli olup olmadığını kontrol et
	valid := ed25519.Verify(publicKey[:], []byte(message), signature)

	return valid, nil
}
