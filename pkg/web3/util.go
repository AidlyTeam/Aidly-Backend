package solana_service

import (
	"crypto/ed25519"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

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

func IsValidSolanaAddress(address string) bool {
	decoded, err := base58.Decode(address)
	if err != nil {
		return false
	}

	return len(decoded) == 32
}
