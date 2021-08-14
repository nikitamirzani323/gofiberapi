package config

import "github.com/tkanos/gonfig"

type ConfigMYSQL struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}

func GetConfigMysql() ConfigMYSQL {
	conf := ConfigMYSQL{}
	gonfig.GetConf("config/confmysql.json", &conf)
	return conf
}
