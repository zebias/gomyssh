package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// SSH客户端
func main() {
	config := &ssh.ClientConfig{
		User: "*****", // 登录名
		Auth: []ssh.AuthMethod{
			ssh.Password("****"), // 密码
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "0.0.0.0:22", config) //IP:PORT
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()
	fd := int(os.Stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	defer term.Restore(fd, state)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed:", err)
	}

	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	fmt.Println(session.Wait())
}
