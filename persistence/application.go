package persistence

import (
	"encoding/hex"
	"log"

	typesGenesis "github.com/pokt-network/pocket/shared/types/genesis"
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

//type App struct {
	//state         protoimpl.MessageState
	//sizeCache     protoimpl.SizeCache
	//unknownFields protoimpl.UnknownFields

	//Address         []byte   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	//PublicKey       []byte   `protobuf:"bytes,2,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	//Paused          bool     `protobuf:"varint,3,opt,name=paused,proto3" json:"paused,omitempty"`
	//Status          int32    `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
	//Chains          []string `protobuf:"bytes,5,rep,name=chains,proto3" json:"chains,omitempty"`
	//MaxRelays       string   `protobuf:"bytes,6,opt,name=max_relays,json=maxRelays,proto3" json:"max_relays,omitempty"`
	//StakedTokens    string   `protobuf:"bytes,7,opt,name=staked_tokens,json=stakedTokens,proto3" json:"staked_tokens,omitempty"`
	//PausedHeight    uint64   `protobuf:"varint,8,opt,name=paused_height,json=pausedHeight,proto3" json:"paused_height,omitempty"`
	//UnstakingHeight int64    `protobuf:"varint,9,opt,name=unstaking_height,json=unstakingHeight,proto3" json:"unstaking_height,omitempty"` // DISCUSS: Why is this int64 but the above is a uint64?
	//Output          []byte   `protobuf:"bytes,10,opt,name=output,proto3" json:"output,omitempty"`
//}

func (p PostgresContext) GetAllApps(height int64) (apps []*typesGenesis.App, err error) {
	ctx, conn, err := p.DB.GetCtxAndConnection()
	if err != nil {
		return
	}
	rows, err := conn.Query(ctx, schema.SelectAll(schema.AppTableName, height))
	
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var a *typesGenesis.App

		var height int64
		err := rows.Scan(&a.Address, &a.PublicKey, &a.StakedTokens, &a.MaxRelays, &a.Output, &a.PausedHeight, &a.UnstakingHeight, &height) 
		if err != nil {
			return nil, err
		}

		a.Paused = height >= int64(a.PausedHeight) 

		apps = append(apps, a)
	}

	return 
}
