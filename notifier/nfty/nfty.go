package nfty

import (
	"errors"
	"net/http"
	"strings"

	"github.com/eze-kiel/itsy/config"
	"github.com/eze-kiel/itsy/shared"
)

const (
	nftyDomain = "https://ntfy.sh/"
	nftyTopic  = "is-there-snow-yet"
)

type NftyNotifier struct {
	Topic        string
	CallbackAddr string
	EmbedImg     bool
	EmbedImgLink string
}

func New(c config.Config) (NftyNotifier, error) {
	if c.NftyTopic == "" {
		return NftyNotifier{}, errors.New("cannot use empty topic for the nfty notifier")
	}

	return NftyNotifier{
		Topic:        c.NftyTopic,
		CallbackAddr: c.NftyCallbackAddr,
		EmbedImg:     c.NftyEmbedImg,
		EmbedImgLink: c.ImgLink,
	}, nil
}

func (n NftyNotifier) Send(name string, isSnow bool) error {
	var msg string
	if isSnow {
		msg = shared.SnowMessage
	} else {
		msg = shared.NoSnowMessage
	}

	req, err := http.NewRequest("POST", nftyDomain+n.Topic, strings.NewReader(msg))
	if err != nil {
		return err
	}
	req.Header.Set("Title", "Monitor: "+name)

	if n.CallbackAddr != "" {
		req.Header.Set("Click", n.CallbackAddr)
	}

	if n.EmbedImg {
		req.Header.Set("Attach", n.EmbedImgLink)
	}

	_, err = http.DefaultClient.Do(req)
	return err
}
