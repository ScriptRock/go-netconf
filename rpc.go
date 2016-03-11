package netconf

import (
	"bytes"
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"io"
)

type RPCMethod interface {
	MarshalRPCMethod() string
}

// a raw RPC method is just a plain string
type RPCMethodRaw string

func (m RPCMethodRaw) MarshalRPCMethod() string {
	return string(m)
}

func RPCMethodLock(target string) RPCMethod {
	return RPCMethodRaw(fmt.Sprintf("<lock><target><%s/></target></lock>", target))
}
func MethodUnlock(target string) RPCMethod {
	return RPCMethodRaw(fmt.Sprintf("<unlock><target><%s/></target></unlock>", target))
}
func MethodGetConfig(source string) RPCMethod {
	return RPCMethodRaw(fmt.Sprintf("<get-config><source><%s/></source></get-config>", source))
}

// a message to send to the remote
type RPCMessage struct {
	ID      string
	Methods []RPCMethod
}

func NewRPCMessage(methods []RPCMethod) *RPCMessage {
	return &RPCMessage{
		ID:      rpcUUID(),
		Methods: methods,
	}
}

func (m *RPCMessage) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var buf bytes.Buffer
	for _, method := range m.Methods {
		buf.WriteString(method.MarshalRPCMethod())
	}

	data := struct {
		ID      string `xml:"message-id,attr"`
		Methods []byte `xml:",innerxml"`
	}{m.ID, buf.Bytes()}

	start.Name.Local = "rpc"
	return e.EncodeElement(data, start)
}

// the remote's response to valid RPC messages
type RPCReply struct {
	XMLName  xml.Name   `xml:"rpc-reply"`
	Errors   []RPCError `xml:"rpc-error,omitempty"`
	Data     string     `xml:",innerxml"`
	Ok       bool       `xml:",omitempty"`
	RawReply string     `xml:"-"`
}

// an error returned over the course of RPC communication
type RPCError struct {
	Type     string `xml:"error-type"`
	Tag      string `xml:"error-tag"`
	Severity string `xml:"error-severity"`
	Path     string `xml:"error-path"`
	Message  string `xml:"error-message"`
	Info     string `xml:",innerxml"`
}

func (re *RPCError) Error() string {
	return fmt.Sprintf("netconf rpc [%s] '%s'", re.Severity, re.Message)
}

// construct a simple UUID for rpc messages
func rpcUUID() string {
	b := make([]byte, 16)
	io.ReadFull(rand.Reader, b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
