---
page_title: "wireguard Provider"
description: |-
  Generates WireGuard X25519 key pairs. Private keys are ephemeral; only public keys are stored in state.
---

# wireguard Provider

Generates WireGuard X25519 key pairs with no secrets stored in Terraform state.

- **`wireguard_private_key`** (ephemeral) — generates a random key pair per apply, never persisted.
- **`wireguard_public_key`** (resource) — accepts the ephemeral private key via write-only argument, stores only the derived public key.

Requires Terraform >= 1.11 (write-only arguments).

## Example Usage

```terraform
provider "wireguard" {}
```

## Schema

No configuration required.
