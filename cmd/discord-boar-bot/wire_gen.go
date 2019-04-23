// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/fsufitch/discord-boar-bot/common"
)

// Injectors from wire.go:

func InitializeCLIRuntime() (*CLIRuntime, error) {
	configuration, err := common.ConfigurationFromEnvironment()
	if err != nil {
		return nil, err
	}
	loggerModule := common.NewLoggerModule()
	cliLogModule := common.CreateCLILogModule(configuration, loggerModule)
	cliRuntime := NewCLIRuntime(configuration, loggerModule, cliLogModule)
	return cliRuntime, nil
}
