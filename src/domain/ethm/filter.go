package ethm

import "errors"

type AddrConfigCache interface {
	Get(key string) (string, error)
}

type EthereumWriteFilter struct {
	addrCache AddrConfigCache
}

func (sel *EthereumWriteFilter) Filter(addr string) (err error) {
	symbol, err := sel.addrCache.Get(addr)
	if err != nil {
		return err
	}
	if symbol == "" {
		return errors.New("not allow to write it")
	}
	return nil
}
