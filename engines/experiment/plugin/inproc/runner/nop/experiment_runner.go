package nop

import (
	"log"

	plugin "github.com/gojek/turing/engines/experiment/plugin/inproc/runner"
	"github.com/gojek/turing/engines/experiment/runner/nop"
)

// init ensures this runner is registered when the package is imported.
func init() {
	err := plugin.Register("nop", nop.NewExperimentRunner)
	if err != nil {
		log.Fatal(err)
	}
}
