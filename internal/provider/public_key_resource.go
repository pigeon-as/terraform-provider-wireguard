// Copyright (c) Pigeon AS
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &PublicKeyResource{}

func NewPublicKeyResource() resource.Resource {
	return &PublicKeyResource{}
}

// PublicKeyResource stores a WireGuard public key in state, derived from a
// write-only private key. The private key is never persisted.
type PublicKeyResource struct{}

// PublicKeyResourceModel maps the resource schema.
// PrivateKey is write-only — available in config during Create/Update but
// always null in state.
type PublicKeyResourceModel struct {
	PrivateKeyWO        types.String `tfsdk:"private_key_wo"`
	PrivateKeyWOVersion types.Int64  `tfsdk:"private_key_wo_version"`
	PublicKey           types.String `tfsdk:"public_key"`
}

func (r *PublicKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_key"
}

func (r *PublicKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Stores a WireGuard public key in state, derived from a write-only private key. " +
			"The private key is accepted via a write-only argument and never stored in state or plan. " +
			"Bump `private_key_wo_version` to trigger key rotation.",

		Attributes: map[string]schema.Attribute{
			"private_key_wo": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded WireGuard private key (write-only). " +
					"Pass the value from `ephemeral.wireguard_private_key.<name>.private_key`. " +
					"Never stored in state or plan.",
				Required:  true,
				WriteOnly: true,
			},
			"private_key_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Version number for the private key. Bump this value to trigger " +
					"key rotation — the resource will update and recompute the public key " +
					"from the new ephemeral private key.",
				Required: true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Base64-encoded WireGuard public key, derived from the private key. " +
					"This is the only value stored in state. Use it for peer list assembly.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *PublicKeyResource) Configure(_ context.Context, _ resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	// No provider data needed — key derivation is self-contained.
}

func (r *PublicKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read write-only value from config (not plan — write-only is null in plan).
	var config PublicKeyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.PrivateKeyWO.IsNull() || config.PrivateKeyWO.IsUnknown() {
		resp.Diagnostics.AddError("Missing Private Key", "private_key_wo is required.")
		return
	}

	publicKey, err := publicKeyFromPrivate(config.PrivateKeyWO.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Key Derivation Failed", err.Error())
		return
	}

	// Read version from plan (normal attribute).
	var plan PublicKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := PublicKeyResourceModel{
		PublicKey:           types.StringValue(publicKey),
		PrivateKeyWOVersion: plan.PrivateKeyWOVersion,
		// PrivateKeyWO is write-only — framework nullifies it automatically.
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PublicKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PublicKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Public key is in state. No external system to reconcile with.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PublicKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Triggered by private_key_version change.
	// Read write-only value from config.
	var config PublicKeyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.PrivateKeyWO.IsNull() || config.PrivateKeyWO.IsUnknown() {
		resp.Diagnostics.AddError("Missing Private Key", "private_key_wo is required for key rotation.")
		return
	}

	publicKey, err := publicKeyFromPrivate(config.PrivateKeyWO.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Key Derivation Failed", err.Error())
		return
	}

	// Read version from plan.
	var plan PublicKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := PublicKeyResourceModel{
		PublicKey:           types.StringValue(publicKey),
		PrivateKeyWOVersion: plan.PrivateKeyWOVersion,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PublicKeyResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No external state to clean up.
}
