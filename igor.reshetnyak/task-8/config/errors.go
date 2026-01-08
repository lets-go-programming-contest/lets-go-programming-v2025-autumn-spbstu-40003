package config

import "errors"

var (
	ErrProviderNotInitialized = errors.New("config provider not initialized")
	ErrDevConfigNotEmbedded   = errors.New("dev config not embedded")
	ErrProdConfigNotEmbedded  = errors.New("prod config not embedded")
	ErrEnvironmentRequired    = errors.New("environment is required")
	ErrLogLevelRequired       = errors.New("log_level is required")
)
