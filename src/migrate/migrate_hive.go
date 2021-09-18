package migrate

import (
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/hive"
)

type HiveMigrate struct {
	hiveCli *hive.HiveCLI
}

func (sel *HiveMigrate) Create() {
	dbs := []string{
		conf.GetHiveEthereumDb().DbName,
	}
	for _, dbName := range dbs {
		if err := sel.hiveCli.Exec("CREATE DATABASE  %s IF NOT EXISTS" + dbName); err != nil {
			panic(err)
		}
	}
}

func (sel *HiveMigrate) Migrate() {
	tables := []string{
		`
 CREATE TABLE transaction_records(                
   id bigint,                                     
   created_at string,                             
   block_number string,                           
   block_hash string,                             
   tx_hash string,                                
   tx_time string,                              
   contract_address string,                       
   from_addr string,                              
   to_addr string,                                
   gas_price string,                              
   value string,                                  
   from_addr_balance string,                      
   to_addr_balance string)
`,
		` CREATE TABLE contract_transaction_records(                
   id bigint,                                     
   created_at string,                             
   block_number string,                           
   block_hash string,                             
   tx_hash string,                                
   tx_time string,                              
   contract_address string,                       
   from_addr string,                              
   to_addr string,                                
   gas_price string,                              
   value string,                                  
   from_addr_balance string,                      
   to_addr_balance string)
`,
	}
	for _, tb := range tables {
		if err := sel.hiveCli.Exec("CREATE TABLES  IF NOT EXISTS  %s  (" + tb + ")"); err != nil {
			panic(err)
		}
	}

}
func (sel *HiveMigrate) Drop() {

}
