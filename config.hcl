"disable_mlock" = true

"listener" = {
  "tcp" = {
    "address" = "127.0.0.1:8200"
  }
}

"storage" = {
  "file" = {
    "path" = "/home/ac/go/src/github.com/kubevault/hq/dist/vault/data"
  }
}

"ui" = true