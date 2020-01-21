package main

import (
	"fmt"

	arg "github.com/alexflint/go-arg"
	validator "github.com/go-playground/validator/v10"
)

type Config struct {
	Endpoint string `arg:"--endpoint" default:"127.0.0.1:6379" help:"redis endpoint" validate:"tcp_addr"`
	Stream   string `arg:"--stream,required" help:"key name of redis stream" validate:"required"`
	Mode     string `arg:"--mode,required" help:"mode: produce or consume" validate:"oneof=produce consume"`

	// produce options
	SourceFormat string `arg:"--fmt" default:"csv" help:"produce mode: source format (only csv supported)" validate:"oneof=csv blob"`
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

	switch cfg.Mode {
	case "produce":
		if cfg.SourceFormat != "csv" && cfg.SourceFormat != "blob" {
			return nil, fmt.Errorf("--fmt can be set to only 'csv' or 'blob'")
		}
	}

	return &cfg, nil
}
