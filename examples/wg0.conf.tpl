[Interface]
PrivateKey = ${private_key}

%{ for peer in peers ~}
[Peer]
PublicKey  = ${peer.public_key}
Endpoint   = ${peer.endpoint}
AllowedIPs = 0.0.0.0/0, ::/0
%{ endfor ~}
