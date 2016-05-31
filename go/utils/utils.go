package utils

import (
  "flag"
)

func Parse() {
  var sshKey string

  flag.StringVar(&sshKey, "k", "", "location of SSH key")
  flag.Parse()

  if len(sshKey) == 0 {
    flag.PrintDefaults()
  }
}
