package netconf

import "encoding/xml"

type Session struct {
	ID                 int
	ServerCapabilities []Capability
	Transport          Transport
	ErrOnWarnings      bool
}

func (s *Session) Close() error {
	return s.Transport.Close()
}

func NewSession(t Transport) (*Session, error) {

	// recieve hello message from remote
	remoteHello, err := t.ReceiveHello()
	if err != nil {
		return nil, err
	}

	// send hello message with default capabilities
	capabilities := []Capability{CapabilityNetconfBase}
	if err = t.SendHello(&TransportHelloMessage{Capabilities: capabilities}); err != nil {
		return nil, err
	}

	return &Session{
		ID:                 remoteHello.SessionID,
		ServerCapabilities: remoteHello.Capabilities,
		Transport:          t,
	}, nil
}

func (s *Session) Exec(methods ...RPCMethod) (*RPCReply, error) {
	// prepare RPC message
	req, err := xml.Marshal(NewRPCMessage(methods))
	if err != nil {
		return nil, err
	}

	if err = s.Transport.Send(req); err != nil {
		return nil, err
	}

	res, err := s.Transport.Recieve()
	if err != nil {
		return nil, err
	}

	// unmarshal the reply
	reply := &RPCReply{RawReply: string(res)}
	if err := xml.Unmarshal(res, reply); err != nil {
		return nil, err
	}

	// handle RPC errors
	if reply.Errors != nil {

		// only act upon non-warning errors unless we have been instructed to fail on warnings
		for _, rpcErr := range reply.Errors {
			if rpcErr.Severity == "error" || s.ErrOnWarnings {
				return reply, &rpcErr
			}
		}
	}

	return reply, nil
}
