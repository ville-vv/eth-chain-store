package hive

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/beltran/gohive"
	"github.com/ville-vv/eth-chain-store/src/common/hive/gen-go/hiveserver2"
)

var defaultCtx = context.Background()
var password = ""
var userName = ""

func Connect(addr string) (err error) {
	//thrift.NewTSSLServerSocket()
	var transport thrift.TTransport
	cfg := &thrift.TConfiguration{}
	transport, err = thrift.NewTSocketConf(addr, cfg)
	if err != nil {
		return
	}
	defer transport.Close()
	transport, err = thrift.NewTTransportFactory().GetTransport(transport)
	if err != nil {
		return
	}

	saslConfiguration := map[string]string{"username": "x", "password": "x"}
	transport, err = gohive.NewTSaslTransport(transport, addr, "PLAIN", saslConfiguration)
	if err != nil {
		return
	}

	hiveConnClient := hiveserver2.NewTCLIServiceClientFactory(transport, thrift.NewTBinaryProtocolFactoryConf(cfg))
	if err = transport.Open(); err != nil {
		return
	}

	argvalue0 := hiveserver2.NewTOpenSessionReq()
	resp, err := hiveConnClient.OpenSession(defaultCtx, argvalue0)
	if err != nil {
		return
	}
	if resp.Status.StatusCode != hiveserver2.TStatusCode_SUCCESS_STATUS {
		return err
	}
	fmt.Println(resp)

	getInfoReq := hiveserver2.NewTExecuteStatementReq()
	getInfoReq.SessionHandle = resp.SessionHandle
	getInfoReq.Statement = "USE test"
	getInfoReq.RunAsync = true
	fmt.Println(hiveConnClient.ExecuteStatement(defaultCtx, getInfoReq))

	getInfoReq.SessionHandle = resp.SessionHandle
	getInfoReq.Statement = "show databases"
	getInfoReq.RunAsync = true

	oprt, err := hiveConnClient.ExecuteStatement(defaultCtx, getInfoReq)
	fmt.Println(oprt, err)

	resultGet := hiveserver2.NewTGetResultSetMetadataReq()
	resultGet.OperationHandle = oprt.OperationHandle

	fmt.Println(hiveConnClient.GetResultSetMetadata(defaultCtx, resultGet))

	return
}
