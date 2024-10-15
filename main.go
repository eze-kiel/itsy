package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/eze-kiel/itsy/config"
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
	flag.StringVar(&url, "img-url", "", "url of the image to download periodically")
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

	if err := downloadImg(url); err != nil {
		logrus.Fatal(err)
	}

	// aka what is the percentage of pixels that match our definition of
	// "this pixel is snow"
	confidence, err := analyzeImg(imgFilePath)
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

func downloadImg(u string) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(imgFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func analyzeImg(f string) (float64, error) {
	fd, err := os.Open(f)
	if err != nil {
		return -1, fmt.Errorf("cannot open the image '%s': %s", f, err)
	}
	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return -1, fmt.Errorf("cannot decode the image: %s", err)
	}

	totalPixels, snowPixels := 0, 0

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			totalPixels++
			pixel := img.At(x, y)

			if isSnowy(pixel) {
				snowPixels++
			}
		}
	}

	percent := (float64(snowPixels) / float64(totalPixels)) * 100
	return percent, nil
}

func isSnowy(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	// 200 out of 255
	threshold := uint32(200 << 8)

	return r > threshold && g > threshold && b > threshold
}
