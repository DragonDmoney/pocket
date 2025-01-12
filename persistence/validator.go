package persistence

import (
	"encoding/hex"
	"log"

	"github.com/pokt-network/pocket/persistence/schema"
	"github.com/pokt-network/pocket/shared/types"
)

func (p PostgresContext) GetValidatorExists(address []byte, height int64) (exists bool, err error) {
	return p.GetExists(schema.ValidatorActor, address, height)
}

func (p PostgresContext) GetValidator(address []byte, height int64) (operator, publicKey, stakedTokens, serviceURL, outputAddress string, pausedHeight, unstakingHeight int64, err error) {
	actor, err := p.GetActor(schema.ValidatorActor, address, height)
	operator = actor.Address
	publicKey = actor.PublicKey
	stakedTokens = actor.StakedTokens
	serviceURL = actor.ActorSpecificParam
	outputAddress = actor.OutputAddress
	pausedHeight = actor.PausedHeight
	unstakingHeight = actor.UnstakingHeight
	return
}

func (p PostgresContext) InsertValidator(address []byte, publicKey []byte, output []byte, _ bool, _ int, serviceURL string, stakedTokens string, pausedHeight int64, unstakingHeight int64) error {
	return p.InsertActor(schema.ValidatorActor, schema.BaseActor{
		Address:            hex.EncodeToString(address),
		PublicKey:          hex.EncodeToString(publicKey),
		StakedTokens:       stakedTokens,
		ActorSpecificParam: serviceURL,
		OutputAddress:      hex.EncodeToString(output),
		PausedHeight:       pausedHeight,
		UnstakingHeight:    unstakingHeight,
	})
}

func (p PostgresContext) UpdateValidator(address []byte, serviceURL string, stakedTokens string) error {
	return p.UpdateActor(schema.ValidatorActor, schema.BaseActor{
		Address:            hex.EncodeToString(address),
		StakedTokens:       stakedTokens,
		ActorSpecificParam: serviceURL,
	})
}

func (p PostgresContext) DeleteValidator(address []byte) error {
	log.Println("[DEBUG] DeleteValidator is a NOOP")
	return nil
}

func (p PostgresContext) GetValidatorsReadyToUnstake(height int64, _ int) ([]*types.UnstakingActor, error) {
	return p.GetActorsReadyToUnstake(schema.ValidatorActor, height)
}

func (p PostgresContext) GetValidatorStatus(address []byte, height int64) (int, error) {
	return p.GetActorStatus(schema.ValidatorActor, address, height)
}

func (p PostgresContext) SetValidatorUnstakingHeightAndStatus(address []byte, unstakingHeight int64, _ int) error {
	return p.SetActorUnstakingHeightAndStatus(schema.ValidatorActor, address, unstakingHeight)
}

func (p PostgresContext) GetValidatorPauseHeightIfExists(address []byte, height int64) (int64, error) {
	return p.GetActorPauseHeightIfExists(schema.ValidatorActor, address, height)
}

func (p PostgresContext) SetValidatorsStatusAndUnstakingHeightIfPausedBefore(pausedBeforeHeight, unstakingHeight int64, _ int) error {
	return p.SetActorStatusAndUnstakingHeightIfPausedBefore(schema.ValidatorActor, pausedBeforeHeight, unstakingHeight)
}

func (p PostgresContext) SetValidatorPauseHeight(address []byte, height int64) error {
	return p.SetActorPauseHeight(schema.ValidatorActor, address, height)
}

// TODO(team): The Get & Update operations need to be made atomic
// TODO(team): Deprecate this functiona altogether and use UpdateValidator where applicable
func (p PostgresContext) SetValidatorStakedTokens(address []byte, tokens string) error { //
	height, err := p.GetHeight()
	if err != nil {
		return err
	}
	operator, _, _, serviceURL, _, _, _, err := p.GetValidator(address, height)
	if err != nil {
		return err
	}
	addr, err := hex.DecodeString(operator)
	if err != nil {
		return err
	}
	return p.UpdateValidator(addr, serviceURL, tokens)
}

func (p PostgresContext) GetValidatorStakedTokens(address []byte, height int64) (tokens string, err error) {
	_, _, tokens, _, _, _, _, err = p.GetValidator(address, height)
	return
}

func (p PostgresContext) GetValidatorOutputAddress(operator []byte, height int64) (output []byte, err error) {
	return p.GetActorOutputAddress(schema.ValidatorActor, operator, height)
}

// TODO(team): implement missed blocks
func (p PostgresContext) SetValidatorPauseHeightAndMissedBlocks(address []byte, pausedHeight int64, missedBlocks int) error {
	return nil
}

// TODO(team): implement missed blocks
func (p PostgresContext) SetValidatorMissedBlocks(address []byte, missedBlocks int) error {
	return nil
}

// TODO(team): implement missed blocks
func (p PostgresContext) GetValidatorMissedBlocks(address []byte, height int64) (int, error) {
	return 0, nil
}
