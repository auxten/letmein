package main

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/CovenantSQL/auth-aws-sg-dynamic-ip/auth"
	"github.com/CovenantSQL/auth-aws-sg-dynamic-ip/aws"
)

type Config struct {
	Auth  auth.Auth         `yaml:"Auth"`
	AwsSg aws.SecurityGroup `yaml:"AwsSg"`
}

func (conf *Config) LoadConfig(path string) error {
	sourceConfig, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithError(err).Error("read config failed")
		return err
	}
	if err = yaml.Unmarshal(sourceConfig, conf); err != nil {
		logrus.WithError(err).Error("load config failed")
		return err
	}

	return nil
}
