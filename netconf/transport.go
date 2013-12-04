package netconf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
)

const (
	MSG_SEPERATOR = "]]>]]>"
)

var DEFAULT_CAPABILITIES = []string{
	"urn:ietf:params:xml:ns:netconf:base:1.0",
}

type HelloMessage struct {
	XMLName      xml.Name `xml:"hello"`
	Capabilities []string `xml:"capabilities>capability"`
	SessionID    int      `xml:"session-id,omitempty"`
}

type Transport interface {
	Send([]byte) error
	Receive() ([]byte, error)
	Close() error
	ReceiveHello() (*HelloMessage, error)
	SendHello(*HelloMessage) error
}

type transportBasicIO struct {
	io.ReadWriteCloser
	chunkedFraming bool
}

func (t *transportBasicIO) Writeln(b []byte) (int, error) {
	t.Write(b)
	t.Write([]byte("\n"))
	return 0, nil
}

// Sends a well formated netconf rpc message as a slice of bytes adding on the
// nessisary framining messages.
func (t *transportBasicIO) Send(data []byte) error {
	t.Write(data)
	t.Write([]byte(MSG_SEPERATOR))
	t.Write([]byte("\n"))
	return nil // TODO: Implement error handling!
}

func (t *transportBasicIO) Receive() ([]byte, error) {
	return t.WaitForBytes([]byte(MSG_SEPERATOR))
}

func (t *transportBasicIO) WaitForBytes(m []byte) ([]byte, error) {
	var out bytes.Buffer
	buf := make([]byte, 4096)

	for {
		n, err := t.Read(buf)

		if n == 0 {
			return nil, fmt.Errorf("WaitForBytes read no data.")
		}

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		end := bytes.Index(buf, m)
		if end > -1 {
			out.Write(buf[0:end])
			return out.Bytes(), nil
		}
		out.Write(buf[0:n])
	}

	return nil, fmt.Errorf("WaitForBytes failed")
}

func (t *transportBasicIO) WaitForString(m string) (string, error) {
	out, err := t.WaitForBytes([]byte(m))
	if out != nil {
		return string(out), err
	}
	return "", err
}

func (t *transportBasicIO) WaitForRegexp(re *regexp.Regexp) ([]byte, [][]byte, error) {
	var out bytes.Buffer

	buf := make([]byte, 4096)
	for {
		n, err := t.Read(buf)

		if n == 0 {
			break // TODO: Handle Error
		}

		if err != nil {
			if err != io.EOF {
				return nil, nil, err
			}
			break
		}

		loc := re.FindSubmatchIndex(buf)
		if loc != nil {
			out.Write(buf[0:loc[1]])

			var matches [][]byte
			for i := 2; i < len(loc); i += 2 {
				matches = append(matches, buf[loc[i]:loc[i+1]])
			}

			return out.Bytes(), matches, nil
		}
		out.Write(buf[0:n])
	}
	return nil, nil, fmt.Errorf("WaitForRegexp failed")
}

func (t *transportBasicIO) SendHello(hello *HelloMessage) error {
	val, err := xml.MarshalIndent(hello, "  ", "    ")
	if err != nil {
		return err
	}

	err = t.Send(val)
	return err
}

func (t *transportBasicIO) ReceiveHello() (*HelloMessage, error) {
	hello := new(HelloMessage)

	val, err := t.Receive()
	if err != nil {
		return hello, err
	}

	err = xml.Unmarshal([]byte(val), hello)
	return hello, err
}

type ReadWriteCloser struct {
	io.Reader
	io.WriteCloser
}

func NewReadWriteCloser(r io.Reader, w io.WriteCloser) *ReadWriteCloser {
	return &ReadWriteCloser{r, w}
}
