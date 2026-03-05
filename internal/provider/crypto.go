// Copyright (c) Pigeon AS
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/base64"
	"fmt"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// generateKeyPair generates a random WireGuard X25519 key pair.
// Returns base64-encoded private key and public key.
func generateKeyPair() (privateKey, publicKey string, err error) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return "", "", fmt.Errorf("generating private key: %w", err)
	}

	return key.String(), key.PublicKey().String(), nil
}

// publicKeyFromPrivate derives the WireGuard public key from a base64-encoded
// private key.
func publicKeyFromPrivate(privKeyB64 string) (string, error) {
	privBytes, err := base64.StdEncoding.DecodeString(privKeyB64)
	if err != nil {
		return "", fmt.Errorf("decoding private key: %w", err)
	}
	if len(privBytes) != wgtypes.KeyLen {
		return "", fmt.Errorf("invalid private key size: got %d, want %d", len(privBytes), wgtypes.KeyLen)
	}

	key, err := wgtypes.NewKey(privBytes)
	if err != nil {
		return "", fmt.Errorf("parsing private key: %w", err)
	}

	return key.PublicKey().String(), nil
}
