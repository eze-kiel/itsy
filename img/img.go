package img

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
)

func Download(u, imgPath string) error {
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func Analyze(f string) (float64, error) {
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
