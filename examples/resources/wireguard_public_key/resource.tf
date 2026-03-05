ephemeral "wireguard_private_key" "node" {}

resource "wireguard_public_key" "node" {
  private_key_wo         = ephemeral.wireguard_private_key.node.private_key
  private_key_wo_version = 1
}
