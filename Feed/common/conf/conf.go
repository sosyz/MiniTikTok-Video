package conf

import "github.com/spf13/viper"

type Redis struct {
	Host     string
	Port     int
	Password string
	Db       int
}

type Pg struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	Db       string
	TimeZone string
}

type Neo4j struct {
	Host     string
	Port     int
	User     string
	Password string
	Realm    string
}

type Secret struct {
	SecretId  string
	SecretKey string
}

type S3 struct {
	Region    string
	Endpoint  string
	Bucket    string
	SecretId  string
	SecretKey string
}

type Listen struct {
	Host string
	Port int
}

type Node struct {
	Id int64
}

type Server Listen

type Run struct {
	Listen Listen
	Server Server
}

type Vol struct {
	Ak string
	Sk string
}

type ContainerConfig struct {
	AppName            string
	AppVersion         string
	AppRoot            string
	CFGAccessKey       string
	ENV                string
	RegionName         string
	SecAppToken        string
	ServiceDiscoverUrl string
	ImageTag           string
	ListenHost         string
	ListenPort         int
}

func ReadS3ConfigByEnv() *S3 {
	viper.SetEnvPrefix("DREAM_S3")
	viper.AutomaticEnv()
	return &S3{
		Region:    viper.GetString("REGION"),
		Endpoint:  viper.GetString("ENDPOINT"),
		SecretId:  viper.GetString("SECRETID"),
		SecretKey: viper.GetString("SECRETKEY"),
		Bucket:    viper.GetString("BUCKET"),
	}
}

func ReadRedisConfigByEnv() *Redis {
	viper.SetEnvPrefix("DREAM_REDIS")
	viper.AutomaticEnv()
	return &Redis{
		Host:     viper.GetString("HOST"),
		Port:     viper.GetInt("PORT"),
		Password: viper.GetString("PASSWORD"),
		Db:       viper.GetInt("DB"),
	}
}

func ReadSecretByEnv(prefix string) *Secret {
	viper.SetEnvPrefix("DREAM_" + prefix)
	viper.AutomaticEnv()
	return &Secret{
		SecretId:  viper.GetString("SECRETID"),
		SecretKey: viper.GetString("SECRETKEY"),
	}
}

func ReadNeo4jConfigByEnv() *Neo4j {
	viper.SetEnvPrefix("DREAM_NEO4J")
	viper.AutomaticEnv()
	return &Neo4j{
		Host:     viper.GetString("HOST"),
		Port:     viper.GetInt("PORT"),
		User:     viper.GetString("USER"),
		Password: viper.GetString("PASSWORD"),
		Realm:    viper.GetString("REALM"),
	}
}

func ReadNodeConfigByEnv() *Node {
	viper.SetEnvPrefix("DREAM_NODE")
	viper.AutomaticEnv()
	return &Node{
		Id: viper.GetInt64("ID"),
	}
}

func ReadContainerConfig() *ContainerConfig {
	viper.SetEnvPrefix("DREAM")
	viper.AutomaticEnv()
	return &ContainerConfig{
		AppName:            viper.GetString("APP_NAME"),
		AppVersion:         viper.GetString("APP_VERSION"),
		AppRoot:            viper.GetString("APP_ROOT"),
		CFGAccessKey:       viper.GetString("CFG_ACCESS_KEY"),
		ENV:                viper.GetString("ENV"),
		RegionName:         viper.GetString("REGION_NAME"),
		SecAppToken:        viper.GetString("SEC_APP_TOKEN"),
		ServiceDiscoverUrl: viper.GetString("SERVICE_DISCOVERY_URI"),
		ImageTag:           viper.GetString("IMAGE_TAG"),
		ListenHost:         viper.GetString("LISTEN_HOST"),
		ListenPort:         viper.GetInt("LISTEN_PORT"),
	}
}
