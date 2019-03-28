package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)

func main() {

	var (
		auth []ssh.AuthMethod

		err error
		key string
	)

	//key = "/Users/luis/Documents/lin_key/linyong"
	key = "/home/lin/.ssh/id_rsa"

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	auth = make([]ssh.AuthMethod, 0)

	pemBytes, err := ioutil.ReadFile(key)
	if err != nil {
		panic(err)
	}

	var signer ssh.Signer

	signer, err = ssh.ParsePrivateKey(pemBytes)

	if err != nil {
		panic(err)
	}
	auth = append(auth, ssh.PublicKeys(signer))

	config := &ssh.ClientConfig{
		User:            "lintest",
		Auth:            auth,
		HostKeyCallback: hostKeyCallbk,
	}
	client, err := ssh.Dial("tcp", "192.168.0.87:22", config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		panic("Failed to run: " + err.Error())
	}

	fmt.Println(b.String())

}
