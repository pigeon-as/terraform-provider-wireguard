---
page_title: "wireguard_private_key Ephemeral Resource"
description: |-
  Generates a random WireGuard X25519 key pair. Never stored in state.
---

# wireguard_private_key (Ephemeral Resource)

Generates a random WireGuard X25519 key pair. Both keys are ephemeral and never stored in state or plan. Pass the private key into a `wireguard_public_key` resource via its write-only argument and into provisioners for wg0.conf deployment.

## Example Usage

```terraform
ephemeral "wireguard_private_key" "node" {}

# Use the private key in a write-only argument (not stored in state):
resource "wireguard_public_key" "node" {
  private_key_wo         = ephemeral.wireguard_private_key.node.private_key
  private_key_wo_version = 1
}

# Use the private key in a provisioner (not stored in state):
resource "null_resource" "wg_config" {
  provisioner "file" {
    content     = "PrivateKey = ${ephemeral.wireguard_private_key.node.private_key}"
    destination = "/etc/wireguard/wg0.conf"
  }
}
```

## Schema

### Read-Only

- `private_key` (String, Sensitive) Base64-encoded WireGuard private key.
- `public_key` (String) Base64-encoded WireGuard public key.
