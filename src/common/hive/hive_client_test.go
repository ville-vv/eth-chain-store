package hive

import (
	"fmt"
	"reflect"
	"testing"
)

type MockHiveConfigOption struct {
}

func (m *MockHiveConfigOption) GetHost() string {
	return "localhost"
}

func (m *MockHiveConfigOption) GetPort() int {
	return 10000
}

func (m *MockHiveConfigOption) GetDBName() string {
	return "ethereum"
}

func (m *MockHiveConfigOption) GetAuthMode() string {
	return "NONE"
}

func (m *MockHiveConfigOption) GetUserName() string {
	return ""
}

func (m *MockHiveConfigOption) GetPassword() string {
	return ""
}

type ResultTst struct {
	Index        int32  `json:"top1000_erc20_token.index" column:"top1000_erc20_token.index"`
	ContractAddr string `json:"contract_addr" column:"top1000_erc20_token.token_contract_address"`
	Name         string `json:"name" column:"top1000_erc20_token.name"`
	Symbol       string `json:"symbol" column:"top1000_erc20_token.symbol"`
}

func TestHiveClient_New(t *testing.T) {
	hCli, err := New(&MockHiveConfigOption{})
	if err != nil {
		t.Error(err)
		return
	}
	res := make([]*ResultTst, 0, 0)
	err = hCli.Find("select * from top1000_erc20_token", &res)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(res)
}

func TestHiveClient_New2(t *testing.T) {
	hCli, err := New(&MockHiveConfigOption{})
	if err != nil {
		t.Error(err)
		return
	}
	res := make([]*ResultTst, 0, 0)
	err = hCli.Find("select * from top1000_erc20_token", &res)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i].ContractAddr)
	}
}

func TestHiveCLI_newObj(t *testing.T) {
	aaa := &ResultTst{}
	oType := reflect.TypeOf(aaa).Elem()
	res := mapToStru(oType, map[string]interface{}{
		"top1000_erc20_token.index":                  1234,
		"top1000_erc20_token.token_contract_address": "0x0000000000004946c0e9F43F4Dee607b0eF1fA1c"})
	fmt.Println(res.Interface())
}

func TestHiveCLI_newObj2(t *testing.T) {
	aaaa := make([]*ResultTst, 0)
	oType := reflect.TypeOf(&aaaa)
	fmt.Println(oType.Elem())
	fmt.Println(oType.Elem().Elem().Elem())
	res := mapToStru(oType.Elem().Elem().Elem(), map[string]interface{}{
		"top1000_erc20_token.index":                  int32(1234),
		"top1000_erc20_token.token_contract_address": "0x0000000000004946c0e9F43F4Dee607b0eF1fA1c"})
	fmt.Println(res.Interface())
}

func TestHiveCLI_Count(t *testing.T) {
	hCli, err := New(&MockHiveConfigOption{})
	defer hCli.Close()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(hCli.Count("transaction_records"))
}
