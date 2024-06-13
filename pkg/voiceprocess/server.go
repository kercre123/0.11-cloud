package voiceprocessing

import (
	"fmt"

	"github.com/kercre123/0.11-cloud/pkg/vars"
)

// Server stores the config
type Server struct{}

var LoadedIntents []vars.JsonIntent

// New returns a new server
func New() (*Server, error) {
	var err error
	LoadedIntents, err = vars.LoadIntents()
	if err != nil {
		fmt.Println("error loading intents: " + err.Error())
	}

	return &Server{}, nil
}
