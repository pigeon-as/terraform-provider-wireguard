// Copyright (c) Pigeon AS
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestPrivateKeyEphemeral(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `
ephemeral "wireguard_private_key" "test" {}

resource "wireguard_public_key" "test" {
  private_key_wo         = ephemeral.wireguard_private_key.test.private_key
  private_key_wo_version = 1
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("wireguard_public_key.test", "public_key"),
					resource.TestMatchResourceAttr("wireguard_public_key.test", "public_key",
						regexp.MustCompile(`^[A-Za-z0-9+/]{43}=$`)),
				),
			},
		},
	})
}
