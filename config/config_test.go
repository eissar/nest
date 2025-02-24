package config

import (
	"fmt"
	"testing"
)

func TestDebugGetConfig(t *testing.T) {
	fmt.Printf("%v", GetConfig())
}

func TestDebugGetRecentLibraries(t *testing.T) {
	fmt.Printf("%v", getRecentLibraries())
}
