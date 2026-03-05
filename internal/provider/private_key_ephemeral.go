// Copyright (c) Pigeon AS
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ ephemeral.EphemeralResource = &PrivateKeyEphemeral{}

func NewPrivateKeyEphemeral() ephemeral.EphemeralResource {
	return &PrivateKeyEphemeral{}
}

// PrivateKeyEphemeral generates a random WireGuard X25519 key pair.
// Both private and public keys are ephemeral — never stored in state or plan.
type PrivateKeyEphemeral struct{}

// PrivateKeyEphemeralModel maps the ephemeral resource schema.
type PrivateKeyEphemeralModel struct {
	PrivateKey types.String `tfsdk:"private_key"`
	PublicKey  types.String `tfsdk:"public_key"`
}

func (r *PrivateKeyEphemeral) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_key"
}

func (r *PrivateKeyEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a random WireGuard X25519 key pair. Both keys are ephemeral " +
			"and never stored in state or plan. Pass the private key into a " +
			"`wireguard_public_key` resource via its write-only argument and into " +
			"provisioners for wg0.conf deployment.",

		Attributes: map[string]schema.Attribute{
			"private_key": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded WireGuard private key.",
				Computed:            true,
				Sensitive:           true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded WireGuard public key.",
				Computed:            true,
			},
		},
	}
}

func (r *PrivateKeyEphemeral) Open(ctx context.Context, _ ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	privateKey, publicKey, err := generateKeyPair()
	if err != nil {
		resp.Diagnostics.AddError("Key Generation Failed", err.Error())
		return
	}

	data := PrivateKeyEphemeralModel{
		PrivateKey: types.StringValue(privateKey),
		PublicKey:  types.StringValue(publicKey),
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
