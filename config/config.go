package config

import "github.com/tkanos/gonfig"

type ConfigMYSQL struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}
type ConfigRedis struct {
	DB_HOST     string
	DB_PASSWORD string
	DB_NAME     int
}

func GetConfigMysql() ConfigMYSQL {
	conf := ConfigMYSQL{}
	gonfig.GetConf("config/confmysql.json", &conf)
	return conf
}

func GetConfigRedis() ConfigRedis {
	conf := ConfigRedis{}
	gonfig.GetConf("config/confredis.json", &conf)
	return conf
}
