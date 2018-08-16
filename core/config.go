package core

import (
	"io/ioutil"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type config struct {
	Token          string `yaml:"token"`
	CurrentGuild   string `yaml:"current_guild"`
	CurrentChannel string `yaml:"current_channel"`
}

var homeDirectory string

func init() {
	var err error
	homeDirectory, err = homedir.Dir()
	if err != nil {
		panic(err)
	}
}

func loadConfig() (cfg config, err error) {
	contents, err := ioutil.ReadFile(filepath.Join(homeDirectory, ".config", "cordless.yaml"))
	if err != nil {
		return
	}
	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return
	}
	return
}

func updateConfig(cfg config) (err error) {
	contents, err := yaml.Marshal(cfg)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filepath.Join(homeDirectory, ".config", "cordless.yaml"), contents, 0700)
	if err != nil {
		return
	}
	return
}
