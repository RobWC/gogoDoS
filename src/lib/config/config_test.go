package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	testConfig := NewConfig()
	t.Log(reflect.TypeOf(testConfig))
	if reflect.TypeOf(testConfig).String() != "*config.Config" {
		t.Error("NewConfig unable to return *Config")
	} else {
		t.Log("NewConfig was able to return *Config")
	}
}
