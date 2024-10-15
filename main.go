package main

import (
	"flag"
	_ "image/jpeg"

	"github.com/eze-kiel/itsy/config"
	"github.com/eze-kiel/itsy/img"
	"github.com/eze-kiel/itsy/notifier"
	"github.com/sirupsen/logrus"
)

const imgFilePath = "/tmp/asdf.jpg"

func main() {
	var nkind string
	var url string
	var cname string
	var threshold float64
	var nftyTopic string

	flag.StringVar(&nkind, "notifier", notifier.Terminal, "select notifier to use (term, ntfy)")
	flag.StringVar(&url, "img-url", "", "url of the image to download (mandatory)")
	flag.StringVar(&cname, "name", "snow monitor", "name of the monitor")
	flag.Float64Var(&threshold, "threshold", 25, "confidence threshold, in percent (100 = absolutely sure)")

	// notifier-dedicated flags
	flag.StringVar(&nftyTopic, "nfty-topic", "", "nfty topic to send notifications when using nfty notifier")
	flag.Parse()

	if url == "" {
		logrus.Fatal("-img-url cannot be empty")
	}

	conf := config.Config{
		NftyTopic: nftyTopic,
	}
	n, err := notifier.GetNotifier(nkind, conf)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := img.Download(url, imgFilePath); err != nil {
		logrus.Fatal(err)
	}

	// aka what is the percentage of pixels that match our definition of
	// "this pixel is snow"
	confidence, err := img.Analyze(imgFilePath)
	if err != nil {
		logrus.Fatal(err)
	}

	var isSnow bool
	if confidence >= threshold {
		isSnow = true
	} else {
		isSnow = false
	}

	if err := n.Send(cname, isSnow); err != nil {
		logrus.Fatal(err)
	}
}
