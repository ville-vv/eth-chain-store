// Copyright 2019 The go-eth Authors
// This file is part of the go-eth library.
//
// The go-eth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-eth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-eth library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/core"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/core/forkid"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/p2p/enode"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/rlp"
)

// ethEntry is the "eth" ENR entry which advertises eth protocol
// on the discovery network.
type ethEntry struct {
	ForkID forkid.ID // Fork identifier per EIP-2124

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e ethEntry) ENRKey() string {
	return "eth"
}

// startEthEntryUpdate starts the ENR updater loop.
func (eth *Ethereum) startEthEntryUpdate(ln *enode.LocalNode) {
	var newHead = make(chan core.ChainHeadEvent, 10)
	sub := eth.blockchain.SubscribeChainHeadEvent(newHead)

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-newHead:
				ln.Set(eth.currentEthEntry())
			case <-sub.Err():
				// Would be nice to sync with eth.Stop, but there is no
				// good way to do that.
				return
			}
		}
	}()
}

func (eth *Ethereum) currentEthEntry() *ethEntry {
	return &ethEntry{ForkID: forkid.NewID(eth.blockchain.Config(), eth.blockchain.Genesis().Hash(),
		eth.blockchain.CurrentHeader().Number.Uint64())}
}
