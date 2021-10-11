package conf

import "strconv"

type HiveConfig struct {
	DBEnv     string
	Host      string
	Port      string
	DbName    string
	AuthModel string
	UserName  string
	Password  string
}

func NewHivConfig() *HiveConfig {
	h := &HiveConfig{
		DBEnv:     "local",
		Host:      "localhost",
		AuthModel: "NONE",
		Port:      "10000",
		DbName:    "default",
	}

	return h
}

func (sel *HiveConfig) GetHost() string {
	return sel.Host
}

func (sel *HiveConfig) GetPort() int {
	port, _ := strconv.Atoi(sel.Port)
	return port
}

func (sel *HiveConfig) GetDBName() string {
	return sel.DbName
}

func (sel *HiveConfig) GetAuthMode() string {
	return sel.AuthModel
}

func (sel *HiveConfig) GetUserName() string {
	return sel.UserName
}

func (sel *HiveConfig) GetPassword() string {
	return sel.Password
}

func (sel *HiveConfig) Dcv() {
	sel.discoverFromEnv()
	sel.discoverFromFlag()
}

func (sel *HiveConfig) discoverFromEnv() {
	ReadEnv(&sel.Host, "DB_ENV")
	ReadEnv(&sel.UserName, "HIVE2_USER_NAME")
	ReadEnv(&sel.Password, "HIVE2_PASSWORD")
	ReadEnv(&sel.Host, "HIVE2_HOST")
	ReadEnv(&sel.Port, "HIVE2_PORT")
	ReadEnv(&sel.AuthModel, "HIVE2_AUTH_MODEL")
	ReadEnv(&sel.DbName, "HIVE2_DB_NAME")
}

func (sel *HiveConfig) discoverFromFlag() {
	ReadFlag(&sel.DBEnv, "db_env")
	ReadFlag(&sel.UserName, "hive_db_user")
	ReadFlag(&sel.Password, "hive_db_passwd")
	ReadFlag(&sel.Host, "hive_db_host")
	ReadFlag(&sel.Port, "hive_db_port")
	ReadFlag(&sel.AuthModel, "hive_db_auth_model")
}

func (sel *HiveConfig) ReSetDbName(dbName string) {
	sel.DbName = dbName
}

func GetHiveEthereumDb() *HiveConfig {
	cfg := NewHivConfig()
	cfg.ReSetDbName("etherum_orc")
	cfg.ReSetDbName("eth_test")
	cfg.Dcv()
	return cfg
}
