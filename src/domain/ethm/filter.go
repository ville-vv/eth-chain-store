package ethm

type AddrConfigCache interface {
	Get(key string) (string, error)
}

type EthereumWriteFilter struct {
	addrCache AddrConfigCache
}

func NewEthereumWriteFilter(addrCache AddrConfigCache) *EthereumWriteFilter {
	return &EthereumWriteFilter{addrCache: addrCache}
}

func (sel *EthereumWriteFilter) Filter(addr string) (err error) {
	symbol, err := sel.addrCache.Get(addr)
	if err != nil {
		return err
	}
	if symbol == "" {
		//return errors.New("not allow to write it")
		return nil
	}
	return nil
}
