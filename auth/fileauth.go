package auth

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Auth struct {
	UserPass map[string]string `yaml:"UserPass,omitempty"`
}

func (auth *Auth) IsPermitted(user string, pass string) bool {
	if auth != nil {
		if p, exist := auth.UserPass[user]; exist {
			if p == pass {
				return true
			}
		}
	}
	return false
}

func (auth *Auth) LoadUserPass(path string) error {
	sourceConfig, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithError(err).Error("read user-pass config failed")
		return err
	}
	if err = yaml.Unmarshal(sourceConfig, auth); err != nil {
		logrus.WithError(err).Error("load user-pass config failed")
		return err
	}
	return nil
}
