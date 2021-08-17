package roc

import "github.com/hashicorp/go-hclog"

var LogLevel = hclog.Info

type SystemConfig struct {
	LogLevel string
	Protocol string
}
