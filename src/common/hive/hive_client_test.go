package hive

import (
	"context"
	"fmt"
	"github.com/beltran/gohive"
	"log"
	"reflect"
	"testing"
)

type MockHiveConfigOption struct {
}

func (m *MockHiveConfigOption) GetHost() string {
	return "172.16.16.155"
}

func (m *MockHiveConfigOption) GetPort() int {
	return 10000
}

func (m *MockHiveConfigOption) GetDBName() string {
	return "etherum_orc"
}

func (m *MockHiveConfigOption) GetAuthMode() string {
	return "NOSASL"
}

func (m *MockHiveConfigOption) GetUserName() string {
	return "hadoop"
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
	err = hCli.Find("select * from top1000_erc20_token_orc limit 1000", &res)
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

func TestExecHive(t *testing.T) {
	ctx := context.Background()
	configuration := gohive.NewConnectConfiguration()
	configuration.Service = "hive"
	//configuration.FetchSize = 2000
	configuration.Database = "eth_test"
	configuration.Username = "hadoop"
	// Previously kinit should have done: kinit -kt ./secret.keytab hive/hs2.example.com@EXAMPLE.COM
	connection, errConn := gohive.Connect("172.16.16.155", 10000, "NOSASL", configuration)
	if errConn != nil {
		log.Fatal(errConn)
	}
	cursor := connection.Cursor()
	//cursor.Exec(ctx, "use eth_test")
	//
	//if cursor.Err != nil {
	//	t.Fatal(cursor.Err)
	//}

	cursor.Exec(ctx, "CREATE TABLE IF NOT EXISTS myTable (a INT, b STRING)")
	if cursor.Err != nil {
		t.Fatal(cursor.Err)
	}

	cursor.Exec(ctx, "INSERT INTO myTable VALUES(1, '1'), (2, '2'), (3, '3'), (4, '4')")
	if cursor.Err != nil {
		t.Fatal(cursor.Err)
	}

	cursor.Exec(ctx, "SELECT * FROM myTable")
	if cursor.Err != nil {
		t.Fatal(cursor.Err)
	}

	var i int32
	var s string
	for cursor.HasMore(ctx) {
		if cursor.Err != nil {
			t.Fatal(cursor.Err)
		}
		cursor.FetchOne(ctx, &i, &s)
		if cursor.Err != nil {
			t.Fatal(cursor.Err)
		}
		t.Log(i, s)
	}

	str := `insert into transaction_records values ('1','2021-10-11 16:51:09','10761307','0x5c8d113c5d485d8206f0207accff950c3c6b7c44139bfedf38bb6d41d6550d93','0x17b4357d5998954201928594007207c25da1bfb6ab4afb8c3e27c41163497b55','2020-08-30 17:06:33','','0xea674fdde714fd979de3edf0f56aa9716b898ec8','0x7732c77b42938cb1daa7b9b6fa7b93cf3211a615','1000000000','50072620449963621','',''),('2','2021-10-11 16:51:09','10761307','0x5c8d113c5d485d8206f0207accff950c3c6b7c44139bfedf38bb6d41d6550d93','0x62fcb7ee7d3cddf35499e9946c5862ccbf039a0026bc80c7e1eb5ddb0668601b','2020-08-30 17:06:33','','0xea674fdde714fd979de3edf0f56aa9716b898ec8','0xb47b9c701b668f87ce2eecf3a9a2b30b3af916d1','1000000000','100077333569603821','','')`
	cursor.Close()

	for i := 0; i < 10; i++ {
		cursor = connection.Cursor()

		cursor.Execute(ctx, str, false)
		if cursor.Err != nil {
			t.Fatal(cursor.Err)
		}
		cursor.Close()
	}

	connection.Close()
}
