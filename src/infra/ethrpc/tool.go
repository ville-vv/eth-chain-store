package ethrpc

import (
	"encoding/hex"
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common/hexutil"
	"golang.org/x/crypto/sha3"
	"math/big"
)

// GenMethodId 可以使这个方法生成 method id， 比如：GenMethodId("balanceOf(address)") 得到 0x70a08231
func GenMethodId(method string) string {
	return Keccak256Hash(method)[:10]
}

func Keccak256Hash(string string) string {
	kek := sha3.NewLegacyKeccak256()
	kek.Reset()
	kek.Write([]byte(string))
	return hexutil.Encode(kek.Sum(nil))
}

func has0xPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1]|32) == 'x'
}
func parseErc20NumericProperty(data string) *big.Int {
	if has0xPrefix(data) {
		data = data[2:]
	}
	if len(data) == 64 {
		var n big.Int
		_, ok := n.SetString(data, 16)
		if ok {
			return &n
		}
	}
	return nil
}

func parseErc20StringProperty(data string) string {
	if has0xPrefix(data) {
		data = data[2:]
	}
	if len(data) > 128 {
		n := parseErc20NumericProperty(data[64:128])
		if n != nil {
			l := n.Uint64()
			if 2*int(l) <= len(data)-128 {
				b, err := hex.DecodeString(data[128 : 128+2*l])
				if err == nil {
					return string(b)
				}
			}
		}
	}
	return ""
}
