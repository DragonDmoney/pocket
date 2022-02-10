package modules

import (
	"pocket/shared/events"
)

type NetworkMessage struct {
	Topic events.PocketEventTopic
	Data  []byte
}

type NetworkModule interface {
	PocketModule

	Config(protocol, address, external string, peers []string)
	Init() error
	Listen() error
	Ready() <-chan uint
	Close()
	Done() <-chan uint

	Send(addr string, msg []byte, wrapped bool) error

	ConsensusBroadcast(message []byte) error
	Broadcast(message []byte, isroot bool) error
	Handle()

	Request(addr string, msg []byte, wrapped bool) ([]byte, error)
	Respond(nonce uint32, iserroreof bool, addr string, msg []byte, wrapped bool) error

	Pong(message []byte) error
	Ping(addr string) (bool, error)

	//Broadcast(*context.PocketContext, *p2p_types.NetworkMessage) error
	//Send(*context.PocketContext, *p2p_types.NetworkMessage, types.NodeId) error
	//GetNetwork() p2p_types.Network
}
