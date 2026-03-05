---
page_title: "wireguard_private_key Ephemeral Resource"
description: |-
  Generates a random WireGuard X25519 key pair. Never stored in state.
---

# wireguard_private_key (Ephemeral Resource)

Generates a random WireGuard X25519 key pair. Both keys are ephemeral and never stored in state or plan. Pass the private key into a `wireguard_public_key` resource via its write-only argument and into provisioners for wg0.conf deployment.

## Example Usage

```terraform
ephemeral "wireguard_private_key" "example" {}

resource "wireguard_public_key" "example" {
  private_key_wo         = ephemeral.wireguard_private_key.example.private_key
  private_key_wo_version = 1
}

resource "terraform_data" "example" {
  triggers_replace = wireguard_public_key.example.public_key

  provisioner "file" {
    content = templatestring(
      <<-WG
      [Interface]
      PrivateKey = $${private_key}

      [Peer]
      PublicKey  = $${peer_public_key}
      AllowedIPs = 0.0.0.0/0
      WG
      , {
        private_key     = ephemeral.wireguard_private_key.example.private_key
        peer_public_key = wireguard_public_key.example.public_key
      }
    )
    destination = "/etc/wireguard/wg0.conf"
  }
}
```

## Schema

### Read-Only

- `private_key` (String, Sensitive) Base64-encoded WireGuard private key.
- `public_key` (String) Base64-encoded WireGuard public key.
