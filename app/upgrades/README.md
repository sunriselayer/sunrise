# Sunrise Upgrades

## Testnet

### Genesis

<https://github.com/sunriselayer/sunrise/releases/tag/v0.1.4>
md5: e0e9775c509ff47b0e90d0b1c6310421

Use this binary version when synchronizing from height0. The upgrade is automatically applied when using Cosmovisor with the following settings.

/lib/systemd/system/cosmovisor.service
`Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=true"`

### v0.2.0-test

<https://github.com/sunriselayer/sunrise/releases/tag/v0.2.0>
md5: 5190cc61247aef4bf92ae51740388134

- Execute sendCoin vRISE to new validators
- Delete DAv1 (blob & blobstream module) and add DAv2 (da module)
