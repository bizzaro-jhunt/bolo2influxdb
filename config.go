package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/starkandwayne/metrics/influxdb"
)

type BoloConfig struct {
	Addr string `json:"ip"`
	Port string `json:"port"`
}

type Config struct {
	Bolo              BoloConfig      `json:"bolo"`
	Influx            influxdb.Config `json:"influx"`
	SkipSSLValidation bool            `json:"skip_ssl_validation"`
}

func LoadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	if cfg.SkipSSLValidation {
		cfg.Influx.InsecureSkipVerify = cfg.SkipSSLValidation
	}
	return cfg, nil
}
