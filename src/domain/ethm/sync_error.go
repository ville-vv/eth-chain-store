package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
)

type SyncErrorManager struct {
	errorRepo repo.SyncErrorRepository
}

func NewSyncErrorManager(errorRepo repo.SyncErrorRepository) *SyncErrorManager {
	return &SyncErrorManager{errorRepo: errorRepo}
}
