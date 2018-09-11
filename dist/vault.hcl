listener "tcp" {
  address = "0.0.0.0:8200"
  tls_disable = true
}

storage "file" {
  path = "/home/ac/go/src/github.com/kubevault/hq/dist/vault/data"
}

disable_mlock= true

ui = true

