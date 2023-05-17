// Copyright 2023 Blink Labs, LLC.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package localtxsubmission

import (
	"fmt"
	"github.com/blinklabs-io/gouroboros/protocol"
)

// Server implements the LocalTxSubmission server
type Server struct {
	*protocol.Protocol
	config *Config
}

// NewServer returns a new Server object
func NewServer(protoOptions protocol.ProtocolOptions, cfg *Config) *Server {
	s := &Server{
		config: cfg,
	}
	protoConfig := protocol.ProtocolConfig{
		Name:                ProtocolName,
		ProtocolId:          ProtocolId,
		Muxer:               protoOptions.Muxer,
		ErrorChan:           protoOptions.ErrorChan,
		Mode:                protoOptions.Mode,
		Role:                protocol.ProtocolRoleServer,
		MessageHandlerFunc:  s.messageHandler,
		MessageFromCborFunc: NewMsgFromCbor,
		StateMap:            StateMap,
		InitialState:        stateIdle,
	}
	s.Protocol = protocol.New(protoConfig)
	return s
}

func (s *Server) messageHandler(msg protocol.Message, isResponse bool) error {
	var err error
	switch msg.Type() {
	case MessageTypeSubmitTx:
		err = s.handleSubmitTx(msg)
	case MessageTypeDone:
		err = s.handleDone()
	default:
		err = fmt.Errorf("%s: received unexpected message type %d", ProtocolName, msg.Type())
	}
	return err
}

func (s *Server) handleSubmitTx(msgGeneric protocol.Message) error {
	if s.config.SubmitTxFunc == nil {
		return fmt.Errorf("received local-tx-submission SubmitTx message but no callback function is defined")
	}
	msg := msgGeneric.(*MsgSubmitTx)
	// Call the user callback function
	return s.config.SubmitTxFunc(msg.Transaction)
}

func (s *Server) handleDone() error {
	return nil
}
