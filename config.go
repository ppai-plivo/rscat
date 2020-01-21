package main

import (
	"fmt"

	arg "github.com/alexflint/go-arg"
	validator "github.com/go-playground/validator/v10"
)

type Config struct {
	Endpoint string `arg:"--endpoint" default:"127.0.0.1:6379" help:"redis endpoint" validate:"tcp_addr"`
	Stream   string `arg:"--stream" help:"key name of redis stream"`
	Mode     string `arg:"--mode" help:"mode: produce or consume"`
	SetKey   string `arg:"--set-key" help:"simple key value set"`

	// produce options
	SourceFormat string `arg:"--fmt" help:"produce mode: source format (csv or blob)"`
	Id           string `arg:"--id" default:"auto" help:"produce mode: possible values - auto or linenum"`
	FieldName    string `arg:"--field-name" default:"blob" help:"produce mode (blob): name for field"`
	Silent       bool   `arg:"--silent" help:"produce mode: suppress printing id"`

	// consume options
	Block uint16 `arg:"--block" default:"1" help:"consume mode: time in seconds to block for"`
}

func initConfig() (*Config, error) {

	var cfg Config
	arg.MustParse(&cfg)

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %s", err.Error())
	}

	return &cfg, nil
}
