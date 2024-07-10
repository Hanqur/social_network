package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database   DatabaseCnf   `yaml:"database"`
	HTTPServer HTTPServerCnf `yaml:"http_server"`
}

type HTTPServerCnf struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	ReadTimeout     int    `yaml:"read_timeout"`
	WriteTimeout    int    `yaml:"write_timeout"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"`
}

func (h HTTPServerCnf) Address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

type DatabaseCnf struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database_name"`
}

func New() *Config {
	f, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var cnf Config

	err = yaml.Unmarshal(f, &cnf)
	if err != nil {
		log.Fatal(err)
	}

	return &cnf
}
