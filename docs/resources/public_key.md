---
page_title: "wireguard_public_key Resource"
description: |-
  Stores a WireGuard public key in state, derived from a write-only private key.
---

# wireguard_public_key (Resource)

Stores a WireGuard public key in state, derived from a write-only private key. The private key is accepted via a write-only argument and never stored in state or plan.

Bump `private_key_wo_version` to trigger key rotation.

## Example Usage

```terraform
ephemeral "wireguard_private_key" "example" {}

resource "wireguard_public_key" "example" {
  private_key_wo         = ephemeral.wireguard_private_key.example.private_key
  private_key_wo_version = 1
}
```

## Schema

### Required

- `private_key_wo` (String, Write-Only) Base64-encoded WireGuard private key. Pass the value from `ephemeral.wireguard_private_key.<name>.private_key`. Never stored in state or plan.
- `private_key_wo_version` (Number) Version number for the private key. Bump this value to trigger key rotation.

### Read-Only

- `public_key` (String) Base64-encoded WireGuard public key, derived from the private key.
