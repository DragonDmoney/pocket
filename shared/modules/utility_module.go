package modules

type UnstakingActor interface {
	GetAddress() []byte
	GetStakeAmount() string
	GetOutputAddress() []byte
}

// TODO(andrew): Look into a better way to structure this interface.
type UtilityContext interface {
	ReleaseContext()
	GetPersistenceContext() PersistenceContext
	CheckTransaction(tx []byte) error
	GetTransactionsForProposal(proposer []byte, maxTransactionBytes int, lastBlockByzantineValidators [][]byte) (transactions [][]byte, err error)
	ApplyBlock(Height int64, proposer []byte, transactions [][]byte, lastBlockByzantineValidators [][]byte) (appHash []byte, err error)
}

type UtilityModule interface {
	Module
	NewContext(height int64) (UtilityContext, error)
}
