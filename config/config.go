package config

import "github.com/spf13/viper"

var (
	RunMode         string
	Port            string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int

	PostgresqlMaster           string
	PostgresSlave              string
	PostgresDbName             string
	PostgresUser               string
	PostgresPassword           string
	PostgresHost               string
	PostgresSshTunnelingEnable bool
	PostgresSshAddress         string
	PostgresSshUserName        string
	PostgresSshKeyfile         string

	RedisAddresses []string
)

func Initialize() {
	RunMode = viper.GetString("RUN_MODE")
	Port = viper.GetString("PORT")
	ReadTimeout = viper.GetInt("READ_TIMEOUT")
	WriteTimeout = viper.GetInt("WRITE_TIMEOUT")
	ShutdownTimeout = viper.GetInt("SHUTDOWN_TIMEOUT")

	PostgresqlMaster = viper.GetString("POSTGRES_MASTER")
	PostgresSlave = viper.GetString("POSTGRES_SLAVE")
	PostgresDbName = viper.GetString("POSTGRES_DB_NAME")
	PostgresUser = viper.GetString("POSTGRES_USER")
	PostgresPassword = viper.GetString("POSTGRES_PASSWORD")
	PostgresHost = viper.GetString("POSTGRES_HOST")
	PostgresSshTunnelingEnable = viper.GetBool("POSTGRES_SSH_TUNNELING_ENABLE")
	PostgresSshAddress = viper.GetString("POSTGRES_SSH_ADDRESS")
	PostgresSshUserName = viper.GetString("POSTGRES_SSH_USER_NAME")
	PostgresSshKeyfile = viper.GetString("POSTGRES_SSH_KEYFILE")

	RedisAddresses = viper.GetStringSlice("REDIS_ADDRESSES")
}
