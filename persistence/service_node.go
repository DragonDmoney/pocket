package persistence

import (
	"encoding/hex"
	"log"

	"github.com/pokt-network/pocket/persistence/schema"
	"github.com/pokt-network/pocket/shared/types"
)

func (p PostgresContext) GetServiceNodeExists(address []byte, height int64) (exists bool, err error) {
	return p.GetExists(schema.ServiceNodeActor, address, height)
}

func (p PostgresContext) GetServiceNode(address []byte, height int64) (operator, publicKey, stakedTokens, serviceURL, outputAddress string, pausedHeight, unstakingHeight int64, chains []string, err error) {
	actor, err := p.GetActor(schema.ServiceNodeActor, address, height)
	operator = actor.Address
	publicKey = actor.PublicKey
	stakedTokens = actor.StakedTokens
	serviceURL = actor.ActorSpecificParam
	outputAddress = actor.OutputAddress
	pausedHeight = actor.PausedHeight
	unstakingHeight = actor.UnstakingHeight
	chains = actor.Chains
	return
}

func (p PostgresContext) InsertServiceNode(address []byte, publicKey []byte, output []byte, _ bool, _ int, serviceURL string, stakedTokens string, chains []string, pausedHeight int64, unstakingHeight int64) error {
	return p.InsertActor(schema.ServiceNodeActor, schema.BaseActor{
		Address:            hex.EncodeToString(address),
		PublicKey:          hex.EncodeToString(publicKey),
		StakedTokens:       stakedTokens,
		ActorSpecificParam: serviceURL,
		OutputAddress:      hex.EncodeToString(output),
		PausedHeight:       pausedHeight,
		UnstakingHeight:    unstakingHeight,
		Chains:             chains,
	})
}

func (p PostgresContext) UpdateServiceNode(address []byte, serviceURL string, stakedTokens string, chains []string) error {
	return p.UpdateActor(schema.ServiceNodeActor, schema.BaseActor{
		Address:            hex.EncodeToString(address),
		StakedTokens:       stakedTokens,
		ActorSpecificParam: serviceURL,
		Chains:             chains,
	})
}

func (p PostgresContext) DeleteServiceNode(address []byte) error {
	log.Println("[DEBUG] DeleteServiceNode is a NOOP")
	return nil
}

func (p PostgresContext) GetServiceNodeCount(chain string, height int64) (int, error) {
	panic("GetServiceNodeCount not implemented")
}

func (p PostgresContext) GetServiceNodesReadyToUnstake(height int64, _ int) ([]*types.UnstakingActor, error) {
	return p.GetActorsReadyToUnstake(schema.ServiceNodeActor, height)
}

func (p PostgresContext) GetServiceNodeStatus(address []byte, height int64) (int, error) {
	return p.GetActorStatus(schema.ServiceNodeActor, address, height)
}

func (p PostgresContext) SetServiceNodeUnstakingHeightAndStatus(address []byte, unstakingHeight int64, _ int) error {
	return p.SetActorUnstakingHeightAndStatus(schema.ServiceNodeActor, address, unstakingHeight)
}

func (p PostgresContext) GetServiceNodePauseHeightIfExists(address []byte, height int64) (int64, error) {
	return p.GetActorPauseHeightIfExists(schema.ServiceNodeActor, address, height)
}

func (p PostgresContext) SetServiceNodeStatusAndUnstakingHeightIfPausedBefore(pausedBeforeHeight, unstakingHeight int64, _ int) error {
	return p.SetActorStatusAndUnstakingHeightIfPausedBefore(schema.ServiceNodeActor, pausedBeforeHeight, unstakingHeight)
}

func (p PostgresContext) SetServiceNodePauseHeight(address []byte, height int64) error {
	return p.SetActorPauseHeight(schema.ServiceNodeActor, address, height)
}

func (p PostgresContext) GetServiceNodeOutputAddress(operator []byte, height int64) (output []byte, err error) {
	return p.GetActorOutputAddress(schema.ServiceNodeActor, operator, height)
}
