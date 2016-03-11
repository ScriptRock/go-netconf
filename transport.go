package netconf

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"io"
)

const (
	TransportMessageSeparator = "]]>]]>"
)

type TransportHelloMessage struct {
	XMLName      xml.Name     `xml:"hello"`
	Capabilities []Capability `xml:"capabilities>capability"`
	SessionID    int          `xml:"session-id,omitempty"`
}

type Transport interface {
	Send(data []byte) error
	Recieve() ([]byte, error)
	Close() error
	SendHello(*TransportHelloMessage) error
	ReceiveHello() (*TransportHelloMessage, error)
}

// a transport wait func allows us to control the way we read input using a lambda function
type TransportWaitFunc func(buf []byte) (int, error)

// transport object base -- implements the basic I/O operations required by the transport interface
type transportBase struct {
	io.ReadWriteCloser
	chunkedFraming bool
}

func (t *transportBase) WaitForFunc(fn TransportWaitFunc) ([]byte, error) {
	var out bytes.Buffer
	buf := make([]byte, 4096)

	for {
		// read from transport into buffer
		n, err := t.Read(buf)
		if n == 0 {
			return nil, fmt.Errorf("transport base read no data")
		}

		// handle errors from the read
		if err != nil && err == io.EOF {
			return nil, err
		} else if err != nil {
			break
		}

		// call into the given wait func
		end, err := fn(buf)
		if err != nil {
			return nil, err
		}

		// if we encounter the end of input, write and finish
		if end > -1 {
			out.Write(buf[0:end])
			return out.Bytes(), nil
		}

		// otherwise, write what we've read already and continue reading
		out.Write(buf[0:n])
	}
	return nil, fmt.Errorf("transport base wait-function failed")
}

// wait for byte subsequence
func (t *transportBase) WaitForBytes(b []byte) ([]byte, error) {
	return t.WaitForFunc(func(buf []byte) (int, error) {
		return bytes.Index(buf, b), nil
	})
}

// wait for string -- extends wait for bytes
func (t *transportBase) WaitForString(s string) (string, error) {
	out, err := t.WaitForBytes([]byte(s))
	if err != nil {
		return "", err
	}
	if out == nil {
		out = []byte{}
	}
	return string(out), nil
}

// TODO: wait for regular expression

// Send a well formatted RPC message, including framing delimiters where applicable
func (t *transportBase) Send(data []byte) error {
	if _, err := t.Write(data); err != nil {
		return err
	}
	if _, err := t.Write([]byte(TransportMessageSeparator)); err != nil {
		return err
	}
	if _, err := t.Write([]byte{'\n'}); err != nil {
		return err
	}
	return nil
}

func (t *transportBase) Recieve() ([]byte, error) {
	return t.WaitForBytes([]byte(TransportMessageSeparator))
}

func (t *transportBase) SendHello(msg *TransportHelloMessage) error {
	req, err := xml.Marshal(msg)
	if err != nil {
		return err
	}
	if err = t.Send(req); err != nil {
		return err
	}
	return nil
}

func (t *transportBase) ReceiveHello() (*TransportHelloMessage, error) {
	msg := &TransportHelloMessage{}

	res, err := t.Recieve()
	if err != nil {
		return nil, err
	}
	if err = xml.Unmarshal(res, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// we implement ReadWriteCloser here so we can use the result of things like SSH dialing as io.ReadWriteCloser
type TransportReadWriteCloser struct {
	io.Reader
	io.WriteCloser
}

func NewTransportReadWriteCloser(r io.Reader, w io.WriteCloser) *TransportReadWriteCloser {
	return &TransportReadWriteCloser{r, w}
}
