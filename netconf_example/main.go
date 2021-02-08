package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudhousetech/crypto/ssh"
	"github.com/cloudhousetech/netconf"
)

func main() {
	var username, password string
	var port int

	flag.StringVar(&username, "username", "", "login username")
	flag.StringVar(&password, "password", "", "login password")
	flag.IntVar(&port, "port", 830, "netconf port")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] targets...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", flag.Arg(0), port), netconf.SSHConfigPassword(username, password))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("SSH CLIENT: %x\n", string(client.SessionID()))

	// set up netconf client
	ncClient, err := netconf.NewClient(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer ncClient.Close()

	fmt.Printf("NETCONF CLIENT: %v\n", ncClient)

	conf, err := ncClient.GetConfig("running")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("XML DATA: %s\n", conf)
}
