package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Yaml struct {
	MysqlConfig     Mysql     `yaml:"Mysql"`
	ModelFileConfig ModelFile `yaml:"ModelFile"`
}

type Mysql struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
}

type ModelFile struct {
	ModelPath   string   `yaml:"ModelPath"`
	VoPath      string   `yaml:"VoPath"`
	TypePath    string   `yaml:"TypePath"`
	PackageName string   `yaml:"PackageName"`
	IsCover     bool     `yaml:"IsCover"`
	TableName   []string `yaml:"TableName"`
}

var YamlConfig Yaml

func LoadConfig() error {
	content, err := os.ReadFile("./config/conf.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, &YamlConfig)
	if err != nil {
		return err
	}
	return nil
}
