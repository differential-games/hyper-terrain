package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/differential-games/hyper-terrain/pkg/noise"

	"github.com/spf13/cobra"
)

const (
	width  = 2560
	height = 1440
)

const (
	maxGrey          = float64(1 << 16)
	contour          = 0.1
	contourThreshold = 0.001

	invLargestScale = 1 / 305.1
	offset          = -5

	waterThreshold = 0.4
)

func run(_ *cobra.Command, _ []string) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := noise.Fractal{}
	n.Fill(r)

	img := image.NewNRGBA64(image.Rect(0, 0, width, height))

	values := make([]float64, width*height)
	minValue := 0.0
	maxValue := 0.0

	for x := 0; x < width; x++ {
		px := float64(x)*invLargestScale + offset

		for y := 0; y < height; y++ {
			py := float64(y)*invLargestScale + offset
			v := n.Cubic(px, py)
			minValue = math.Min(v, minValue)
			maxValue = math.Max(v, maxValue)
			values[x*height+y] = v
		}
	}

	for i, v := range values {
		grey := (v - minValue) / (maxValue - minValue)
		if grey < waterThreshold {
			grey = minValue
		} else if math.Mod(grey, contour) < contourThreshold {
			grey = minValue
		}

		grey *= maxGrey

		img.Set(i/height, i%height, color.Gray16{Y: uint16(grey)})
	}

	out, err := os.Create("out.png")
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			fmt.Fprintln(os.Stderr, closeErr)
		}
	}()

	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hyper-terrain",
		Short: "hyper-terrain is a fast random terrain generator",
		RunE:  run,
	}
}

func Execute() {
	if err := rootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
