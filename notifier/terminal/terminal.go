package terminal

import (
	"fmt"

	"github.com/eze-kiel/itsy/config"
	"github.com/eze-kiel/itsy/shared"
)

type TermNotifier struct{}

func New(c config.Config) TermNotifier {
	return TermNotifier{}
}

func (n TermNotifier) Send(name string, isSnow bool) error {
	var msg string
	if isSnow {
		msg = shared.SnowMessage
	} else {
		msg = shared.NoSnowMessage
	}

	fmt.Println(msg)
	return nil
}
