// Copyright (c) Pigeon AS
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &WireGuardProvider{}
var _ provider.ProviderWithEphemeralResources = &WireGuardProvider{}

// WireGuardProvider implements the wireguard provider.
type WireGuardProvider struct {
	version string
}

func (p *WireGuardProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "wireguard"
	resp.Version = p.version
}

func (p *WireGuardProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The WireGuard provider generates X25519 key pairs for WireGuard mesh networks. " +
			"Private keys are ephemeral (never stored in state). Public keys are persisted in state " +
			"for peer list assembly. Key rotation is controlled via `private_key_wo_version`.",
	}
}

func (p *WireGuardProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	// No provider configuration needed. Key generation is self-contained.
}

func (p *WireGuardProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPublicKeyResource,
	}
}

func (p *WireGuardProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewPrivateKeyEphemeral,
	}
}

func (p *WireGuardProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WireGuardProvider{
			version: version,
		}
	}
}
