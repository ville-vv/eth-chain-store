package conf

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
}

func (sel *HiveConfig) discoverFromFlag() {
	ReadFlag(&sel.DBEnv, "db_env")
	ReadFlag(&sel.UserName, "db_user")
	ReadFlag(&sel.Password, "db_passwd")
	ReadFlag(&sel.Host, "db_host")
	ReadFlag(&sel.Port, "db_port")
}
