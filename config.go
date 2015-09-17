package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Config struct {
	Port    int
	Limit   int
	Workers int
}

func (c *Config) Load(f string) error {
	configFile, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(configFile, &c); err != nil {
		return err
	}
	if c.Port < 1 || c.Limit < 1 || c.Workers < 2 {
		err = errors.New("Invalid config")
	}
	return err
}
