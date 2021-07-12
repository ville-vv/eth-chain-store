package eth_erc20

import (
	"github.com/ville-vv/eth-chain-store/src/infra/cache"
)

type Erc20TokenCache struct {
	st cache.ICache
}

func NewErc20TokenCache(st cache.ICache) *Erc20TokenCache {
	c := &Erc20TokenCache{st: st}
	return c
}
