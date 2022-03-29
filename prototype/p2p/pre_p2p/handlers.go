package pre_p2p

import (
	"io/ioutil"
	"log"
	"net"
	"pocket/p2p/pre_p2p/types"

	"google.golang.org/protobuf/proto"
)

func (m *networkModule) handleNetworkMessage(conn net.Conn) {
	defer conn.Close()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Println("Error reading from conn: ", err)
		return
	}

	networkMessage := types.P2PMessage{}
	if err := proto.Unmarshal(data, &networkMessage); err != nil {
		panic(err) // TODO remove and handle
	}
	// networkMessage, err := DecodeNetworkMessage(data)
	if err != nil {
		log.Println("Error decoding network message: ", err)
		return
	}

	// temporarily convert

	log.Println(conn.LocalAddr().String(), "Received a network message of topic:", networkMessage.Topic)
	event := types.Event{
		SourceModule: types.P2P,
		PocketTopic:  networkMessage.Topic,
		MessageData:  networkMessage.Data,
	}

	m.GetBus().PublishEventToBus(&event)
}

func (m *networkModule) respondToTelemetryMessage(conn net.Conn) {
	// TODO: quick hack. not running `defer conn.Close()` since the connection is passed
	// to Consensus node for debugging purposes.
	log.Println("Responding to telemetry request...")

	event := types.Event{
		SourceModule: types.P2P,
		PocketTopic:  string(types.CONSENSUS_TELEMETRY_MESSAGE),

		NetworkConnection: conn,
	}
	m.GetBus().PublishEventToBus(&event)
}