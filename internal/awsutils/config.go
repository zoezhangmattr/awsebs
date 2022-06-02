package awsutils

import (
	"os"
	"path/filepath"
	"strings"

	logger "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func LoadProfileConfig() ([]ProfileData, error) {
	f, err := os.OpenFile(AWSDefaultConfigFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf("open aws config file failed, %v", err)
		return nil, err
	}
	defer f.Close()

	cfg, err := ini.Load(f.Name())
	if err != nil {
		logger.WithFields(logger.Fields{"filename": f.Name()}).Error("fail to read ini file: %v", err)
		return nil, err
	}
	pl := []ProfileData{}
	for _, sec := range cfg.Sections() {
		region := sec.Key("region").Value()
		name := sec.Name()
		name = strings.Replace(name, "profile ", "", 1)
		if region == "" {
			continue
		}
		pl = append(pl, ProfileData{
			Region: region,
			Name:   name,
		})
	}
	logger.WithFields(logger.Fields{
		"filename":      f.Name(),
		"profiles":      pl,
		"profile_count": len(pl),
	}).Info("loading aws config file")
	return pl, nil

}

func GetHomeDirectory() string {
	return os.Getenv("HOME")
}

func AWSDefaultConfigDir() string {
	return filepath.Join(GetHomeDirectory(), ".aws")
}
func AWSDefaultConfigFilePath() string {
	return filepath.Join(AWSDefaultConfigDir(), "config")
}
