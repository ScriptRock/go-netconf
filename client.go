package netconf

import (
	"io"

	"github.com/cloudhousetech/crypto/ssh"
)

// TODO: we may want to do some of the fancy concurrency stuff our SFTP client does, but we should probably refactor
// the code behind that so we can reuse it here instead of reinventing the wheel.

// Client works quite a lot like SFTP, which is also a subsystem of SSH
type Client struct {
	session *Session
}

// NewClient is the constructor for Client. It creates a Client with stdin and stdout as the Reader and WriteCloser
func NewClient(conn *ssh.Client) (*Client, error) {
	// establish ssh session
	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	// request netconf subsystem from remote
	if err = session.RequestSubsystem(SSHNetconfSubsystem); err != nil {
		return nil, err
	}

	// set up stdin and stdout pipes
	pw, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	pr, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	return NewClientPipe(conn, session, pr, pw)
}

// NewClientPipe creates a Client using the given Reader and WriteCloser
func NewClientPipe(sshClient *ssh.Client, sshSession *ssh.Session, r io.Reader, w io.WriteCloser) (*Client, error) {
	// set up transport using the given reader and writecloser
	transport := TransportSSHFromSSHClient(sshClient, sshSession)
	transport.ReadWriteCloser = NewTransportReadWriteCloser(r, w)

	ncSession, err := NewSession(transport)
	if err != nil {
		return nil, err
	}

	// set up the client with a new netconf session
	netconf := &Client{
		session: ncSession,
	}
	return netconf, nil
}

// Close closes the netconf session to which the Client is connected
func (c *Client) Close() error {
	return c.session.Close()
}

// GetConfig issues a request to netconf for configuration data.
func (c *Client) GetConfig(typ string) ([]byte, error) {
	reply, err := c.session.Exec(MethodGetConfig(typ))
	if err != nil {
		return nil, err
	}
	return []byte(reply.Data), nil
}
