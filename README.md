# Terraform Provider: WireGuard

Generates WireGuard X25519 key pairs with **no secrets stored in Terraform state** — private keys are ephemeral and flow only into write-only arguments and provisioners.

## Design

```
ephemeral "wireguard_private_key" ─── generates random X25519 key pair (never in state)
         │
         ├── private_key ──→ resource "wireguard_public_key".private_key_wo  (write-only, never in state)
         │                   └── derives public_key (stored in state for peer list assembly)
         │
         └── private_key ──→ provisioner "file" { content = wg0.conf }    (deployed to host)
```

**Key rotation**: bump `private_key_wo_version` → new ephemeral key generated → public key updated → provisioner re-deploys config.

**Steady-state applies**: version unchanged → no diff → no key churn → no mesh disruption.

## Requirements

- Terraform >= 1.11 (write-only arguments)
- Go >= 1.24 (building from source)

## Usage

```hcl
provider "wireguard" {}

ephemeral "wireguard_private_key" "example" {
  for_each = var.nodes
}

resource "wireguard_public_key" "example" {
  for_each = var.nodes

  private_key_wo         = ephemeral.wireguard_private_key.example[each.key].private_key
  private_key_wo_version = var.wireguard_key_version
}

resource "terraform_data" "example" {
  for_each = var.nodes

  triggers_replace = wireguard_public_key.example[each.key].public_key

  provisioner "file" {
    content     = templatefile("${path.module}/wg0.conf.tpl", {
      private_key = ephemeral.wireguard_private_key.example[each.key].private_key
      peers       = [for k, v in wireguard_public_key.example : { public_key = v.public_key, endpoint = var.nodes[k].endpoint } if k != each.key]
    })
    destination = "/etc/wireguard/wg0.conf"
  }
}
```

## Building

```shell
go install
# or
make build
```

## Testing

```shell
make test
```
