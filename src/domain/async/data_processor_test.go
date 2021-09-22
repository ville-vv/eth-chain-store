package async

import (
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"testing"
)

type EthAccountRepositoryMock struct {
}

func (e *EthAccountRepositoryMock) GetTxData() ([]*model.TransactionRecord, error) {
	return nil, nil
}

func (e *EthAccountRepositoryMock) Finish() error {
	return nil
}

func (e *EthAccountRepositoryMock) IsAccountExist(addr string) (bool, error) {
	return false, nil
}

func (e *EthAccountRepositoryMock) QueryAccountByAddr(addr string) (*model.EthereumAccount, error) {
	return nil, nil
}

func (e *EthAccountRepositoryMock) CreateAccount(account *model.EthereumAccount) error {
	return nil
}

func (e *EthAccountRepositoryMock) UpdateAccountBalance(addr string, balance string) error {
	return nil
}

func TestDataProcessorCtl_Process(t *testing.T) {
	//ethAccountProcess := NewEthAccountService(ethrpc.NewClient("http://172.16.16.115:8545"), &EthAccountRepositoryMock{})
	//prc := NewDataProcessorCtl(ethAccountProcess, &EthAccountRepositoryMock{}, nil)
	//prc.Process()
}
