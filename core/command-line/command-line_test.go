package commandline

import (
	"testing"

	"github.com/eissar/nest/config"
)

func TestAdd(t *testing.T) {
	var localFile = "./command-line.go"
	Add(config.GetConfig(), &localFile)
}

func TestList(t *testing.T) {
	var localLimit = 5
	List(config.GetConfig(), &localLimit)
}

func TestSwitch(t *testing.T) {
	cfg := config.GetConfig()
	Switch(cfg, "inspo")
}

// no real point in testing this one right now...
// func TestReveal(t *testing.T) {
// 	var localFile = "./command-line.go"
// 	Reveal(config.GetConfig(), &localFile)
// }
