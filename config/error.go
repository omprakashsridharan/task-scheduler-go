package config

import "errors"

var InvalidConfigFilePathError = errors.New("invalid config file path")
var FileLoadError = errors.New("error while loading config file")
var UnmarshalError = errors.New("config unmarshall error")
var ValidationError = errors.New("config validation error")
