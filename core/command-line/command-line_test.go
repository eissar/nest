package commandline

import (
	"testing"

	"github.com/eissar/nest/config"
)

var cfg = config.GetConfig()

func TestCommandLineSanityChecks(t *testing.T) {
	str := `D:\Dropbox\Pictures\3617502772_a_beautiful_and_elaborate_drawing_of_a_pug_lying_in_the_grass_in_the_style_of_Michelangelo.png`
	cnt := 10
	id := "M83M1XMMQPT3R"

	Add(cfg, &str)
	List(cfg, &cnt)
	Reveal(cfg, &id)
}
