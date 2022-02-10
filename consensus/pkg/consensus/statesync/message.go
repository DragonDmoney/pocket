package statesync

import (
	"bytes"
	"encoding/gob"

	consensus_types "pocket/consensus/pkg/consensus/types"
	"pocket/shared"
)

type StateSyncMessageType string

const (
	StateSyncBlockRequest  StateSyncMessageType = "StateSyncBlockRequest"
	StateSyncBlockResponse StateSyncMessageType = "StateSyncBlockResponse"
)

type StateSyncMessage struct {
	consensus_types.GenericConsensusMessage

	MessageType StateSyncMessageType
	MessageData []byte
}

func (m *StateSyncMessage) GetType() consensus_types.ConsensusMessageType {
	return consensus_types.StateSyncConsensusMessage
}

func (m *StateSyncMessage) Encode() ([]byte, error) {
	bytes, err := shared.GobEncode(m)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (m *StateSyncMessage) Decode(data []byte) error {
	err := shared.GobDecode(data, m)
	if err != nil {
		return err
	}
	return nil
}

func StateSyncMessageFromBytes(data []byte) *StateSyncMessage {
	buff := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buff)
	m := &StateSyncMessage{}
	dec.Decode(m)
	return m
}
