package repository

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type ExportConfig struct {
	Table  string   `yaml:"table"`
	Fields []string `yaml:"fields"`
}

type Configuration struct {
	DbInstance string         `yaml:"dbInstance"`
	GcsBucket  string         `yaml:"gcsBucket"`
	GcsPrefix  string         `yaml:"gcsPrefix"`
	BqDataset  string         `yaml:"bqDataset"`
	Exports    []ExportConfig `yaml:"exports"`
}

func ParseConfiguration(fp string) (Configuration, error) {
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return Configuration{}, err
	}
	var config Configuration
	if err := yaml.Unmarshal(content, &config); err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func FindExportConfig(exports []ExportConfig, table string) (ExportConfig, error) {
	for _, e := range exports {
		if e.Table == table {
			return e, nil
		}
	}
	return ExportConfig{}, fmt.Errorf("not found export match table name \"%s\"", table)
}
