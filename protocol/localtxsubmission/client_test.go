// Copyright 2023 Blink Labs Software
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

package localtxsubmission_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/blinklabs-io/gouroboros/internal/test/ouroboros_mock"
	"github.com/blinklabs-io/gouroboros/ledger"
	"github.com/blinklabs-io/gouroboros/protocol"
	"github.com/blinklabs-io/gouroboros/protocol/localtxsubmission"
)

/*type clientTestDefinition struct {
	CborHex     string
	Message     protocol.Message
	MessageType uint
}*/

// Helper function to allow inline hex decoding without capturing the error
func hexDecode(data string) []byte {
	// Strip off any leading/trailing whitespace in hex string
	data = strings.TrimSpace(data)
	decoded, err := hex.DecodeString(data)
	if err != nil {
		panic(fmt.Sprintf("error decoding hex: %s", err))
	}
	return decoded
}

// Valid CBOR that serves as a placeholder for real TX content in the tests
// [h'DEADBEEF']
var placeholderTx = hexDecode("82008204d818468144DEADBEEF")

// Valid CBOR that serves as a placeholder for TX rejection errors
// [2, 4]
var placeholderRejectError = hexDecode("820204")

var ConversationEntrySubmitTxRequest = ouroboros_mock.ConversationEntry{
	Type:            ouroboros_mock.EntryTypeInput,
	ProtocolId:      localtxsubmission.ProtocolId,
	InputMessage:    localtxsubmission.NewMsgSubmitTx(ledger.EraIdAlonzo, placeholderTx),
	MsgFromCborFunc: localtxsubmission.NewMsgFromCbor,
}

var ConversationEntryAcceptTxResponse = ouroboros_mock.ConversationEntry{
	Type:       ouroboros_mock.EntryTypeOutput,
	ProtocolId: localtxsubmission.ProtocolId,
	IsResponse: true,
	OutputMessages: []protocol.Message{
		localtxsubmission.NewMsgAcceptTx(),
	},
}

var ConversationEntryRejectTxResponse = ouroboros_mock.ConversationEntry{
	Type:       ouroboros_mock.EntryTypeOutput,
	ProtocolId: localtxsubmission.ProtocolId,
	IsResponse: true,
	OutputMessages: []protocol.Message{
		localtxsubmission.NewMsgRejectTx(placeholderRejectError),
	},
}

func TestLocalTxSubmissionAccepted(t *testing.T) {
	mockConn := ouroboros_mock.NewConnection(
		ouroboros_mock.ProtocolRoleClient,
		[]ouroboros_mock.ConversationEntry{
			ouroboros_mock.ConversationEntryHandshakeRequestGeneric,
			ouroboros_mock.ConversationEntryHandshakeNtCResponse,
			ConversationEntrySubmitTxRequest,
			ConversationEntryAcceptTxResponse,
		},
	).(*ouroboros_mock.Connection)

	oConn, err := ouroboros.New(
		ouroboros.WithConnection(mockConn),
		ouroboros.WithNetworkMagic(ouroboros_mock.MockNetworkMagic),
	)
	if err != nil {
		t.Fatalf("unexpected error when creating Ouroboros object: %s", err)
	}

	go func() {
		err, ok := <-oConn.ErrorChan()
		if !ok {
			return
		}
		// We can't call t.Fatalf() from a different Goroutine, so we panic instead
		panic(fmt.Sprintf("unexpected Ouroboros connection error: %s", err))
	}()

	/*oConn.LocalTxSubmission().Client.Start()

	errOnSubmit := oConn.LocalTxSubmission().Client.SubmitTx(ledger.TxTypeAlonzo, placeholderTx)
	if errOnSubmit != nil {
		t.Fatalf("unexpected error file trying to submit a valid transation: %s", errOnSubmit)
	}

	errOnStop := oConn.LocalTxSubmission().Client.Stop()
	if errOnStop != nil {
		t.Fatalf("failed to stop LocalTxSubmission client: %s", errOnStop)
	}*/

	// Close Ouroboros connection
	if err := oConn.Close(); err != nil {
		t.Fatalf("unexpected error when closing Ouroboros object: %s", err)
	}
	// Wait for connection shutdown
	select {
	case <-oConn.ErrorChan():
	case <-time.After(10 * time.Second):
		t.Errorf("did not shutdown within timeout")
	}
}

func TestLocalTxSubmissionRejected(t *testing.T) {
	mockConn := ouroboros_mock.NewConnection(
		ouroboros_mock.ProtocolRoleClient,
		[]ouroboros_mock.ConversationEntry{
			ouroboros_mock.ConversationEntryHandshakeRequestGeneric,
			ouroboros_mock.ConversationEntryHandshakeNtCResponse,
			ConversationEntrySubmitTxRequest,
			ConversationEntryRejectTxResponse,
		},
	).(*ouroboros_mock.Connection)

	oConn, err := ouroboros.New(
		ouroboros.WithConnection(mockConn),
		ouroboros.WithNetworkMagic(ouroboros_mock.MockNetworkMagic),
	)
	if err != nil {
		t.Fatalf("unexpected error when creating Ouroboros object: %s", err)
	}

	go func() {
		err, ok := <-oConn.ErrorChan()
		if !ok {
			return
		}
		// We can't call t.Fatalf() from a different Goroutine, so we panic instead
		panic(fmt.Sprintf("unexpected Ouroboros connection error: %s", err))
	}()

	/*oConn.LocalTxSubmission().Client.Start()

	errOnSubmit := oConn.LocalTxSubmission().Client.SubmitTx(ledger.TxTypeAlonzo, placeholderTx)
	if errOnSubmit != nil {
		t.Fatalf("unexpected error file trying to submit a valid transation: %s", errOnSubmit)
	}

	errOnStop := oConn.LocalTxSubmission().Client.Stop()
	if errOnStop != nil {
		t.Fatalf("failed to stop LocalTxSubmission client: %s", errOnStop)
	}*/

	// Close Ouroboros connection
	if err := oConn.Close(); err != nil {
		t.Fatalf("unexpected error when closing Ouroboros object: %s", err)
	}
	// Wait for connection shutdown
	select {
	case <-oConn.ErrorChan():
	case <-time.After(10 * time.Second):
		t.Errorf("did not shutdown within timeout")
	}
}
