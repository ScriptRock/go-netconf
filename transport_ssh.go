package netconf

import (
	"github.com/cloudhousetech/crypto/ssh"
)

const (
	SSHDefaultNetconfPort = 830
	SSHNetconfSubsystem   = "netconf"
)

type TransportSSH struct {
	transportBase
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
}

func TransportSSHFromSSHClient(client *ssh.Client, session *ssh.Session) *TransportSSH {
	return &TransportSSH{
		SSHClient:  client,
		SSHSession: session,
	}
}

func (t *TransportSSH) Close() error {
	// close any SSH struct which is defined
	if t.SSHSession != nil {
		if err := t.SSHSession.Close(); err != nil {
			return err
		}
	}
	if t.SSHClient != nil {
		if err := t.SSHClient.Close(); err != nil {
			return err
		}
	}
	return nil
}

// SSH config including authentication information
func SSHConfigPassword(username, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}
}
