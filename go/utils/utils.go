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
  "fmt"
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

  session1, err := conn.NewSession()
  if err != nil {
    log.Fatalf("unable to create session: %s", err)
  }
  defer session1.Close()

  var stdoutBuf1 bytes.Buffer
  session1.Stdout = &stdoutBuf1
  err = session1.Run("ls -l")
  if err != nil {
    log.Fatalf("Failed to get listing: %s", err)
  }

  session2, err := conn.NewSession()
  if err != nil {
    log.Fatalf("unable to create session: %s", err)
  }
  defer session2.Close()

  go func() {
		w, _ := session2.StdinPipe()
		defer w.Close()
		content := "123456789\n"
		fmt.Fprintln(w, "D0755", 0, "testdir") // mkdir
		fmt.Fprintln(w, "C0644", len(content), "testfile1")
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00") // transfer end with \x00
		fmt.Fprintln(w, "C0644", len(content), "testfile2")
		fmt.Fprint(w, content)
		fmt.Fprint(w, "\x00")
	}()
	if err := session2.Run("/usr/bin/scp -tr ./"); err != nil {
    log.Fatalf("unable to scp: %s", err)
		//panic("Failed to run: " + err.Error())
	}

  session3, err := conn.NewSession()
  if err != nil {
    log.Fatalf("unable to create session: %s", err)
  }
  defer session3.Close()

  var stdoutBuf2 bytes.Buffer
  session3.Stdout = &stdoutBuf2
  err = session3.Run("ls -l")
  if err != nil {
    log.Fatalf("Failed to get listing: %s", err)
  }

  session4, err := conn.NewSession()
  if err != nil {
    log.Fatalf("unable to create session: %s", err)
  }
  defer session4.Close()

  var stdoutBuf3 bytes.Buffer
  session4.Stdout = &stdoutBuf3
  err = session4.Run("ls -l testdir")
  if err != nil {
    log.Fatalf("Failed to get listing: %s", err)
  }

  return "Before:\n" + stdoutBuf1.String() + "\nAfter\n"+ stdoutBuf2.String() + "\nTestDir :\n" + stdoutBuf3.String()

}
