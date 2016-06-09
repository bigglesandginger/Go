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

/*

json format

{
  credentiols: {
    user:
    cert:
    pwd:
  },
  machines:[
    {
      ip:
    }
  ],
  commands:{

  }
}

*/


func main() {
  flag.Parse()

  if len(*sshKey) == 0 || len(*sshServer) == 0 || len(*user) == 0 || len(*pwd) == 0 {
    flag.PrintDefaults()
  } else {
    config := utils.GetConfig(*user, *pwd, *sshKey)
    conn := utils.CreateConnection(config, *sshServer)
    defer conn.Close()
    var output string = ""
    output += utils.RunCommand(conn, "ls -l")
    utils.CopyToServer(conn, "foobar", "There once was a man from Nantucket")
    output += utils.RunCommand(conn, "ls -l")
    output += utils.RunCommand(conn, "ls -l testdir")

    fmt.Println(output)
  }
}
