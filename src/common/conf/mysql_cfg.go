package conf

import "fmt"

type MysqlConf struct {
	DBEnv       string
	Username    string
	Password    string
	Host        string
	Port        string
	DbName      string
	MaxIdleConn int
	MaxOpenConn int
	DbCharset   string
}

func NewMysqlConf() *MysqlConf {
	m := &MysqlConf{
		DBEnv:       "local",
		Host:        "127.0.0.1",
		Port:        "3306",
		MaxIdleConn: 10,
		MaxOpenConn: 1000,
		DbCharset:   "utf8",
	}
	m.Dcv()
	return m
}

func (sel *MysqlConf) GetUserName() string {
	return sel.Username
}

func (sel *MysqlConf) GetPassword() string {
	return sel.Password
}

func (sel *MysqlConf) GetHost() string {
	return sel.Host
}

func (sel *MysqlConf) GetPort() string {
	return sel.Port
}

func (sel *MysqlConf) GetDbName() string {
	return sel.DbName
}

func (sel *MysqlConf) GetMaxIdleConn() int {
	return sel.MaxIdleConn
}

func (sel *MysqlConf) GetMaxOpenConn() int {
	return sel.MaxOpenConn
}

func (sel *MysqlConf) GetCharset() string {
	return sel.DbCharset
}

func (sel *MysqlConf) Dcv() {
	sel.discoverFromEnv()
	sel.discoverFromFlag()
	return
}

func (sel *MysqlConf) discoverFromEnv() {
	ReadEnv(&sel.DBEnv, "DB_ENV")
	ReadEnv(&sel.Username, "MYSQL_USER_NAME")
	ReadEnv(&sel.Password, "MYSQL_PASSWORD")
	ReadEnv(&sel.Host, "MYSQL_HOST")
	ReadEnv(&sel.Port, "MYSQL_PORT")
	sel.MaxIdleConn = 10
	sel.MaxOpenConn = 1000
}

func (sel *MysqlConf) discoverFromFlag() {
	ReadFlag(&sel.DBEnv, "db_env")
	ReadFlag(&sel.Username, "db_user")
	ReadFlag(&sel.Password, "db_passwd")
	ReadFlag(&sel.Host, "db_host")
	ReadFlag(&sel.Port, "db_port")
}

func (sel *MysqlConf) ReSetDbName(name string) {
	sel.DbName = fmt.Sprintf("eth_%s_%s", name, sel.DBEnv)
}

func GetEthBusinessDbConfig() *MysqlConf {
	mysqlConfig := NewMysqlConf()
	mysqlConfig.ReSetDbName("business")
	return mysqlConfig
}

func GetEthereumDbConfig() *MysqlConf {
	mysqlConfig := NewMysqlConf()
	mysqlConfig.ReSetDbName("ethereum")
	return mysqlConfig
}

func GetEthContractDbConfig() *MysqlConf {
	mysqlConfig := NewMysqlConf()
	mysqlConfig.ReSetDbName("contract")
	return mysqlConfig
}

func GetEthContractTransactionDbConfig() *MysqlConf {
	mysqlConfig := NewMysqlConf()
	mysqlConfig.ReSetDbName("contract_transaction")
	return mysqlConfig
}

func GetEthTransactionDbConfig() *MysqlConf {
	mysqlConfig := NewMysqlConf()
	mysqlConfig.ReSetDbName("transaction")
	return mysqlConfig
}

func GetEthereumHiveMapDbConfig() *MysqlConf {
	mysqlConfig := NewMysqlConf()
	mysqlConfig.ReSetDbName("ethereum_hive_map")
	return mysqlConfig
}
