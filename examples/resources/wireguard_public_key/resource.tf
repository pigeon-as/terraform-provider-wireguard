ephemeral "wireguard_private_key" "example" {}

resource "wireguard_public_key" "example" {
  private_key_wo         = ephemeral.wireguard_private_key.example.private_key
  private_key_wo_version = 1
}
