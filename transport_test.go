package netconf

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

// example RPC packets
const (
	deviceHelloStr = `<!-- user bbennett, class j-super-user -->
<hello>
  <capabilities>
    <capability>urn:ietf:params:xml:ns:netconf:base:1.0</capability>
    <capability>urn:ietf:params:xml:ns:netconf:capability:candidate:1.0</capability>
    <capability>urn:ietf:params:xml:ns:netconf:capability:confirmed-commit:1.0</capability>
    <capability>urn:ietf:params:xml:ns:netconf:capability:validate:1.0</capability>
    <capability>urn:ietf:params:xml:ns:netconf:capability:url:1.0?protocol=http,ftp,file</capability>
    <capability>http://xml.juniper.net/netconf/junos/1.0</capability>
    <capability>http://xml.juniper.net/dmi/system/1.0</capability>
  </capabilities>
  <session-id>19313</session-id>
</hello>
]]>]]>`
	loginStr = `SRX240 (ttyp2)
Password:
--- JUNOS 12.1X45-D15.5 built 2013-09-19 07:42:15 UTC
bbennett@SRX240>
`
	waitForStringResponseStr = `SRX240 (ttyp2)`
)

// construct a test packet to test WaitFor<X> methods
func waitTestStr(input []byte) []byte {
	// the message needs to be at least 4096 for some reason
	prefix := make([]byte, 4096)
	for i := 0; i < len(prefix); i++ {
		prefix[i] = 'x'
	}
	prefix[4095] = '\n'

	// add prefix to the actual information
	return append(prefix, input...)
}

// bogus structs for testing purposes
type testTransport struct {
	transportBase
}

type testReadCloser struct {
	io.Reader
	io.Writer
}

func (c *testReadCloser) Close() error {
	return nil
}

func newTestTransport(input []byte) (*testTransport, *bytes.Buffer) {
	var transport testTransport
	r := bytes.NewReader(input)
	w := new(bytes.Buffer)

	transport.ReadWriteCloser = &testReadCloser{r, w}
	return &transport, w
}

type transportTestItem struct {
	Name     string
	TestFunc func(*TransportHelloMessage) interface{}
	Expected interface{}
}

// verify that the hello messages being sent to the remote are formatted correctly
func TestSendHello(t *testing.T) {

	// define tests
	tests := []transportTestItem{
		{
			Name: "SessionID Nil",
			TestFunc: func(h *TransportHelloMessage) interface{} {
				return h.SessionID
			},
			Expected: 0,
		},
		{
			Name: "Capability length",
			TestFunc: func(h *TransportHelloMessage) interface{} {
				return len(h.Capabilities)
			},
			Expected: 1,
		},
		{
			Name: "Capability #0",
			TestFunc: func(h *TransportHelloMessage) interface{} {
				return h.Capabilities[0]
			},
			Expected: Capability("urn:ietf:params:xml:ns:netconf:base:1.0"),
		},
	}

	transport, out := newTestTransport([]byte(""))
	transport.SendHello(&TransportHelloMessage{Capabilities: []Capability{CapabilityNetconfBase}})

	sent := out.String()
	out.Reset()

	hello := &TransportHelloMessage{}
	if err := xml.Unmarshal([]byte(sent), hello); err != nil {
		t.Fatal(err)
	}

	for _, item := range tests {
		if result := item.TestFunc(hello); result != item.Expected {
			t.Fatalf("unexpected test result: %v vs %v", result, item.Expected)
		}
	}
}

// test that correctly formatted RPC responses from the remote are unmarshalled correctly
func TestRecieveHello(t *testing.T) {
	tests := []transportTestItem{
		{
			Name: "SessionID match",
			TestFunc: func(h *TransportHelloMessage) interface{} {
				return h.SessionID
			},
			Expected: 19313,
		},
		{
			Name: "Capability length",
			TestFunc: func(h *TransportHelloMessage) interface{} {
				return len(h.Capabilities)
			},
			Expected: 7,
		},
		{
			Name: "Capability #0",
			TestFunc: func(h *TransportHelloMessage) interface{} {
				return h.Capabilities[0]
			},
			Expected: Capability("urn:ietf:params:xml:ns:netconf:base:1.0"),
		},
	}

	transport, _ := newTestTransport([]byte(deviceHelloStr))

	hello, err := transport.ReceiveHello()
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range tests {
		if result := item.TestFunc(hello); result != item.Expected {
			t.Fatalf("unexpected test result: %v vs %v", result, item.Expected)
		}
	}
}

// test transportBase's WaitForString function
func TestWaitForString(t *testing.T) {
	transport, _ := newTestTransport(waitTestStr([]byte(loginStr)))
	out, err := transport.WaitForString("Password:")
	if err != nil {
		t.Fatal(err)
	}

	expected := string(waitTestStr([]byte(waitForStringResponseStr)))
	if strings.Trim(out, "x\n\t ") != strings.Trim(expected, "x\n\t ") {
		t.Fatal("incorrect output")
	}
}
