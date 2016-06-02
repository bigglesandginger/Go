package utils

import (
  "crypto/x509"
  "encoding/pem"
  "golang.org/x/crypto/ssh"
  //"io"
  "io/ioutil"
  //"os"
  "log"
  "bytes"
)

func SshToServer(user string, pwd string, key string, server string) string {
  pemKey, err := ioutil.ReadFile(key)
  if err != nil {
      log.Fatalf("unable to read private key: %v", err)
  }

  block, _ := pem.Decode([]byte(pemKey))

  derKey, err := x509.DecryptPEMBlock(block, []byte(pwd))
  if err != nil {
      log.Fatalf("unable to decrypt private key: %v", err)
  }

  privKey, err := x509.ParsePKCS1PrivateKey(derKey)
  if err != nil {
      log.Fatalf("unable to decrypt pkcs1 private key: %v", err)
  }

  signer, err := ssh.NewSignerFromKey(privKey)
  if err != nil {
      log.Fatalf("unable to parse private key: %v", err)
  }

  config := &ssh.ClientConfig{
  	User: user,
  	Auth: []ssh.AuthMethod{
      ssh.PublicKeys(signer),
    },
  }

  conn, err := ssh.Dial("tcp", server, config)
  if err != nil {
    log.Fatalf("unable to connect: %s", err)
  }
  defer conn.Close()

  session, err := conn.NewSession()
  if err != nil {
    log.Fatalf("unable to create session: %s", err)
  }
  defer session.Close()

  var stdoutBuf bytes.Buffer
  session.Stdout = &stdoutBuf
  session.Run("ls -l")

  return stdoutBuf.String()
}
