package config

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
)

var K = koanf.New(".")

func Load(config_path string) {

	f := file.Provider(config_path)

	if err := K.Load(f, toml.Parser()); err != nil {
		zap.L().Error(fmt.Sprintf("Watching config %s", config_path))
		os.Exit(1)
	}
	// Watch the file and get a callback on change. The callback can do whatever,
	// like re-load the configuration.
	// File provider always returns a nil `event`.
	f.Watch(func(event interface{}, err error) {
		if err != nil {
			zap.L().Error("Watch error occured")
			return
		}

		// Throw away the old config and load a fresh copy.
		zap.L().Info("Config changed, reloading")
		K = koanf.New(".")
		K.Load(f, toml.Parser())
		K.Print()
	})

	// Block forever (and manually make a change to mock/mock.json) to
	// reload the config.
	zap.L().Info(fmt.Sprintf("Watching config %s", config_path))
}

func GetAwsRegion() string {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("Unable to get region from config")
			os.Exit(1)
		}
	}()
	region := K.MustString("aws.region")
	return region
}

func GetAwsProfile() string {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("Unable to get profile from config")
		}
	}()
	return K.MustString("aws.profile")
}

func GetJumpHostTagName() string {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("Unable to get Jumphost Tagname from config")
		}
	}()
	return K.MustString("jumphost.tag_name")
}
func GetJumpHostTagValue() string {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("Unable to get Jumphost Tagvalue from config")
		}
	}()
	return K.MustString("jumphost.tag_value")
}

func GetJumpHostSecondaryName() string {
	return K.String("jumphost.secondary_tag_name")
}

func GetJumpHostSecondaryValue() string {
	return K.String("jumphost.secondary_tag_value")
}

func GetJumpHostUser() string {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("Unable to get Jumphost user from config")
		}
	}()
	return K.MustString("jumphost.username")
}

func GetJumphostVpc() string {
	return K.String("jumphost.vpc_name")
}

func GetPublicKey() string {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("Unable to get public key path from config")
		}
	}()
	loc := K.MustString("ssh.pub_key_loc")
	key, err := readFile(loc)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Unable to get public key (%s)", loc))
		zap.L().Error(err.Error())
		os.Exit(1)
	}
	return key
}

func GetPrivateKey() []byte {
	defer func() {	
		if err := recover(); err != nil {
		zap.L().Error("Unable to get private key from config")
		}
	}()
	loc := K.MustString("ssh.private_key_loc")
	key, err := readFileByte(loc)
	if err != nil {
		zap.L().Error("Unable to get private key")
		os.Exit(1)
	}
	return key
}

func readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	s := string(content)
	return s, nil
}

func readFileByte(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}
