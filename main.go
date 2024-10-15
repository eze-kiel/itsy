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
	var snowOnly bool
	flag.StringVar(&nkind, "notifier", notifier.Terminal, "select notifier to use (term, ntfy)")
	flag.StringVar(&url, "img-url", "", "url of the image to download (mandatory)")
	flag.StringVar(&cname, "name", "snow monitor", "name of the monitor")
	flag.Float64Var(&threshold, "threshold", 25, "confidence threshold, in percent (100 = absolutely sure)")
	flag.BoolVar(&snowOnly, "snow-only", false, "send notification only if snow has been detected")

	// notifier-dedicated flags
	var nftyTopic string
	var nftyCallback string
	var nftyEmbedImg bool
	flag.StringVar(&nftyTopic, "nfty-topic", "", "nfty topic to send notifications when using nfty notifier")
	flag.StringVar(&nftyCallback, "nfty-callback-address", "", "if set, you'll be redirected to this address when opening the notification")
	flag.BoolVar(&nftyEmbedImg, "nfty-embed-image", false, "if set, it will embed the downloaded image to the notification (if size < 2Mo)")
	flag.Parse()

	if url == "" {
		logrus.Fatal("-img-url cannot be empty")
	}

	conf := config.Config{
		ImgLink:          url,
		NftyTopic:        nftyTopic,
		NftyCallbackAddr: nftyCallback,
		NftyEmbedImg:     nftyEmbedImg,
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
	} else if !snowOnly {
		isSnow = false
	} else {
		return
	}

	if err := n.Send(cname, isSnow); err != nil {
		logrus.Fatal(err)
	}
}
