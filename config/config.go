package config

import (
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"time"
)

const (
	configFilePath        = "configFilePath"
	defaultConfigFilePath = "config/config.yaml"
)

var (
	configFilePathVar = flag.String(configFilePath, defaultConfigFilePath, "config file path")
)

func Parse() {
	configItems, err := readFromFile()
	if err == nil {
		flag.VisitAll(func(f *flag.Flag) {
			if configItems[f.Name] != nil {
				_ = flag.Set(f.Name, configItems[f.Name].Value.(string))
			}
		})
	}
	flag.Parse()
}

func String(name string, value string, usage string) *string {
	return flag.String(name, value, usage)
}
func Bool(name string, value bool, usage string) *bool {
	return flag.Bool(name, value, usage)
}

func Duration(name string, value time.Duration, usage string) *time.Duration {
	return flag.Duration(name, value, usage)
}

type ConfigItem struct {
	Name  string
	Usage string
	Value interface{}
}

func readFromFile() (map[string]*ConfigItem, error) {
	log.Printf("Using configuration file %s", *configFilePathVar)
	var args []*ConfigItem
	yamlData, err := ioutil.ReadFile(*configFilePathVar)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlData, &args)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*ConfigItem)
	for _, item := range args {
		result[item.Name] = item
	}
	return result, nil
}

func WriteToFile() error {
	var args []interface{}
	flag.VisitAll(func(flag *flag.Flag) {
		if flag.Name != configFilePath {
			args = append(args, convertToConfigItem(flag))
		}
	})
	yamlData, err := yaml.Marshal(args)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(*configFilePathVar, yamlData, 0755)
}

func convertToConfigItem(flag *flag.Flag) *ConfigItem {
	return &ConfigItem{
		Name:  flag.Name,
		Usage: flag.Usage,
		Value: flag.Value,
	}
}
