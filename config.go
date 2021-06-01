package main

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)
type Config struct {
	Port int `yaml:"port"`
	Proxy string `yaml:"proxy"`
	CacheDir string `yaml:"cacheDir"`
	LoongsonPKGS []string `yaml:"loongsonPKGS""`
}

func UnmarshalConfig(path string) Config {
	conf := Config{}
	yamlFile,err := ioutil.ReadFile(path)
	if err != nil{
		panic(err)
	}
	if err := yaml.Unmarshal(yamlFile,&conf);err != nil{
		panic(err)
	}
	return conf
}
