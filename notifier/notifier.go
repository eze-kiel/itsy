package notifier

import (
	"fmt"

	"github.com/eze-kiel/itsy/config"
	"github.com/eze-kiel/itsy/notifier/nfty"
	"github.com/eze-kiel/itsy/notifier/terminal"
)

const (
	Nfty     = "nfty"
	Terminal = "term"
)

type Notifier interface {
	Send(name string, isSnow bool) error
}

func GetNotifier(k string, c config.Config) (Notifier, error) {
	switch k {
	case Nfty:
		return nfty.New(c)
	case Terminal:
		return terminal.New(c), nil
	default:
		return nil, fmt.Errorf("notifier '%s' is not supported", k)
	}
}
