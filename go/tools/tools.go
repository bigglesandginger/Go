package main

import (
  "github.com/bigglesandginger/common/go/utils"
  "fmt"
  "flag"
)

var (
  sshKey = flag.String("k", "", "location of SSH key")
  sshServer = flag.String("h", "", "server to SSH to")
  user = flag.String("u", "", "username")
  pwd = flag.String("p", "", "password")
)

func main() {
  flag.Parse()

  if len(*sshKey) == 0 || len(*sshServer) == 0 || len(*user) == 0 || len(*pwd) == 0 {
    flag.PrintDefaults()
  } else {
    output := utils.SshToServer(*user, *pwd, *sshKey, *sshServer)
    fmt.Println(output)
  }
}
