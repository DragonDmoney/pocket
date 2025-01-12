package persistence

import (
	"encoding/hex"
	"log"

	"github.com/pokt-network/pocket/persistence/schema"
	"github.com/pokt-network/pocket/shared/types"
)

func (p PostgresContext) GetAppExists(address []byte, height int64) (exists bool, err error) {
	return p.GetExists(schema.ApplicationActor, address, height)
}

func (p PostgresContext) GetApp(address []byte, height int64) (operator, publicKey, stakedTokens, maxRelays, outputAddress string, pauseHeight, unstakingHeight int64, chains []string, err error) {
	actor, err := p.GetActor(schema.ApplicationActor, address, height)
	operator = actor.Address
	publicKey = actor.PublicKey
	stakedTokens = actor.StakedTokens
	maxRelays = actor.ActorSpecificParam
	outputAddress = actor.OutputAddress
	pauseHeight = actor.PausedHeight
	unstakingHeight = actor.UnstakingHeight
	chains = actor.Chains
	return
}

func (p PostgresContext) InsertApp(address []byte, publicKey []byte, output []byte, _ bool, _ int, maxRelays string, stakedTokens string, chains []string, pausedHeight int64, unstakingHeight int64) error {
	return p.InsertActor(schema.ApplicationActor, schema.BaseActor{
		Address:            hex.EncodeToString(address),
		PublicKey:          hex.EncodeToString(publicKey),
		StakedTokens:       stakedTokens,
		ActorSpecificParam: maxRelays,
		OutputAddress:      hex.EncodeToString(output),
		PausedHeight:       pausedHeight,
		UnstakingHeight:    unstakingHeight,
		Chains:             chains,
	})
}

func (p PostgresContext) UpdateApp(address []byte, maxRelays string, stakedTokens string, chains []string) error {
	return p.UpdateActor(schema.ApplicationActor, schema.BaseActor{
		Address:            hex.EncodeToString(address),
		StakedTokens:       stakedTokens,
		ActorSpecificParam: maxRelays,
		Chains:             chains,
	})
}

func (p PostgresContext) DeleteApp(_ []byte) error {
	log.Println("[DEBUG] DeleteApp is a NOOP")
	return nil
}

func (p PostgresContext) GetAppsReadyToUnstake(height int64, _ int) ([]*types.UnstakingActor, error) {
	return p.GetActorsReadyToUnstake(schema.ApplicationActor, height)
}

func (p PostgresContext) GetAppStatus(address []byte, height int64) (int, error) {
	return p.GetActorStatus(schema.ApplicationActor, address, height)
}

func (p PostgresContext) SetAppUnstakingHeightAndStatus(address []byte, unstakingHeight int64, _ int) error {
	return p.SetActorUnstakingHeightAndStatus(schema.ApplicationActor, address, unstakingHeight)
}

func (p PostgresContext) GetAppPauseHeightIfExists(address []byte, height int64) (int64, error) {
	return p.GetActorPauseHeightIfExists(schema.ApplicationActor, address, height)
}

func (p PostgresContext) SetAppStatusAndUnstakingHeightIfPausedBefore(pausedBeforeHeight, unstakingHeight int64, _ int) error {
	return p.SetActorStatusAndUnstakingHeightIfPausedBefore(schema.ApplicationActor, pausedBeforeHeight, unstakingHeight)
}

func (p PostgresContext) SetAppPauseHeight(address []byte, height int64) error {
	return p.SetActorPauseHeight(schema.ApplicationActor, address, height)
}

func (p PostgresContext) GetAppOutputAddress(operator []byte, height int64) ([]byte, error) {
	return p.GetActorOutputAddress(schema.ApplicationActor, operator, height)
}
