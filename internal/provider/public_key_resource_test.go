// Copyright (c) Pigeon AS
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func TestPublicKeyResource(t *testing.T) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testPublicKeyResourceConfig(key.String(), 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wireguard_public_key.test", "public_key", key.PublicKey().String()),
					resource.TestCheckResourceAttr("wireguard_public_key.test", "private_key_wo_version", "1"),
					resource.TestCheckNoResourceAttr("wireguard_public_key.test", "private_key_wo"),
				),
			},
		},
	})
}

func TestPublicKeyResource_PublicKeyFormat(t *testing.T) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testPublicKeyResourceConfig(key.String(), 1),
				Check: resource.TestMatchResourceAttr("wireguard_public_key.test", "public_key",
					regexp.MustCompile(`^[A-Za-z0-9+/]{43}=$`)),
			},
		},
	})
}

func TestPublicKeyResource_KeyRotation(t *testing.T) {
	key1, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	key2, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testPublicKeyResourceConfig(key1.String(), 1),
				Check:  resource.TestCheckResourceAttr("wireguard_public_key.test", "public_key", key1.PublicKey().String()),
			},
			{
				Config: testPublicKeyResourceConfig(key2.String(), 2),
				Check:  resource.TestCheckResourceAttr("wireguard_public_key.test", "public_key", key2.PublicKey().String()),
			},
		},
	})
}

func TestPublicKeyResource_InvalidPrivateKey(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      testPublicKeyResourceConfig("not-valid-base64!!!", 1),
				ExpectError: regexp.MustCompile(`Key Derivation Failed`),
			},
		},
	})
}

func testPublicKeyResourceConfig(privateKey string, version int) string {
	return fmt.Sprintf(`resource "wireguard_public_key" "test" {
  private_key_wo         = %q
  private_key_wo_version = %d
}`, privateKey, version)
}
