// Copyright (c) Pigeon AS
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/base64"
	"testing"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func TestGenerateKeyPair(t *testing.T) {
	t.Run("returns valid base64 keys", func(t *testing.T) {
		priv, pub, err := generateKeyPair()
		if err != nil {
			t.Fatalf("generateKeyPair: %v", err)
		}

		privBytes, err := base64.StdEncoding.DecodeString(priv)
		if err != nil {
			t.Fatalf("decoding private key: %v", err)
		}
		if len(privBytes) != wgtypes.KeyLen {
			t.Errorf("private key size: got %d, want %d", len(privBytes), wgtypes.KeyLen)
		}

		pubBytes, err := base64.StdEncoding.DecodeString(pub)
		if err != nil {
			t.Fatalf("decoding public key: %v", err)
		}
		if len(pubBytes) != wgtypes.KeyLen {
			t.Errorf("public key size: got %d, want %d", len(pubBytes), wgtypes.KeyLen)
		}
	})

	t.Run("public key matches private key", func(t *testing.T) {
		priv, pub, err := generateKeyPair()
		if err != nil {
			t.Fatal(err)
		}

		// Independently derive public key from private key via wgtypes.
		privBytes, _ := base64.StdEncoding.DecodeString(priv)
		key, err := wgtypes.NewKey(privBytes)
		if err != nil {
			t.Fatalf("wgtypes.NewKey: %v", err)
		}

		if pub != key.PublicKey().String() {
			t.Errorf("public key mismatch: got %s, derived %s", pub, key.PublicKey().String())
		}
	})

	t.Run("generates unique keys", func(t *testing.T) {
		_, pub1, _ := generateKeyPair()
		_, pub2, _ := generateKeyPair()
		if pub1 == pub2 {
			t.Error("two calls produced the same public key")
		}
	})
}

func TestPublicKeyFromPrivate(t *testing.T) {
	t.Run("round-trip with generateKeyPair", func(t *testing.T) {
		priv, expectedPub, err := generateKeyPair()
		if err != nil {
			t.Fatal(err)
		}

		pub, err := publicKeyFromPrivate(priv)
		if err != nil {
			t.Fatalf("publicKeyFromPrivate: %v", err)
		}
		if pub != expectedPub {
			t.Errorf("public key mismatch: got %s, want %s", pub, expectedPub)
		}
	})

	t.Run("invalid base64", func(t *testing.T) {
		_, err := publicKeyFromPrivate("not-valid-base64!!!")
		if err == nil {
			t.Error("expected error for invalid base64")
		}
	})

	t.Run("wrong key size", func(t *testing.T) {
		short := base64.StdEncoding.EncodeToString([]byte("too-short"))
		_, err := publicKeyFromPrivate(short)
		if err == nil {
			t.Error("expected error for wrong key size")
		}
	})
}
