package conf

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
	readEnv(&sel.DBEnv, "DB_ENV")
	readEnv(&sel.Username, "MYSQL_USER_NAME")
	readEnv(&sel.Password, "MYSQL_PASSWORD")
	readEnv(&sel.Host, "MYSQL_HOST")
	readEnv(&sel.Port, "MYSQL_PORT")
	sel.MaxIdleConn = 10
	sel.MaxOpenConn = 1000
	//if sel.Host == "" {
	//	sel.Host = "127.0.0.1"
	//}
	//if sel.Port == "" {
	//	sel.Port = "3306"
	//}
}

func (sel *MysqlConf) discoverFromFlag() {
	readFlag(&sel.DBEnv, "db_env")
	readFlag(&sel.Username, "db_user")
	readFlag(&sel.Password, "db_passwd")
	readFlag(&sel.Host, "db_host")
	readFlag(&sel.Port, "db_port")
}