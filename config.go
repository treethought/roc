package roc

import "github.com/hashicorp/go-hclog"

var LogLevel = hclog.Debug

type SystemConfig struct {
	LogLevel string
	Protocol string
}
