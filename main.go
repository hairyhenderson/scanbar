package main

import (
	"errors"
	"fmt"
	"image/color"
	"os"

	"gopkg.in/bieber/barcode.v0"

	"gocv.io/x/gocv"
)

func main() {
	err := dothething()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func dothething() error {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return err
	}
	defer webcam.Close()

	window := gocv.NewWindow("I'M HUNGRY FOR BARCODES!")
	defer window.Close()

	mat := gocv.NewMat()
	defer mat.Close()

	scanner := barcode.NewScanner().SetEnabledAll(true)
	count := 0
	for {
		count++
		if ok := webcam.Read(&mat); !ok {
			return errors.New("cannot read from webcam")
		}
		if mat.Empty() {
			continue
		}
		// only scan once every 50 frames so it's not crazy choppy
		if count >= 5 {
			count = 0
			src, err := mat.ToImage()
			if err != nil {
				return err
			}
			img := barcode.NewImage(src)
			symbols, err := scanner.ScanImage(img)
			if err != nil {
				return err
			}
			if len(symbols) > 0 {
				for _, s := range symbols {
					fmt.Printf("Found a %s bar code! Data is: %s\n", s.Type.Name(), s.Data)
					c := color.RGBA{0, 0, 255, 0}
					if len(s.Boundary) > 0 {
						rect := gocv.BoundingRect(s.Boundary)
						gocv.Rectangle(&mat, rect, c, 3)
					}
				}
			}
		}
		window.IMShow(mat)
		window.WaitKey(10)

	}
	return nil
}
