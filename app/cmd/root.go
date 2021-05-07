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

	"github.com/willbeason/hyper-terrain/pkg/noise"

	"github.com/spf13/cobra"
)

const (
	Width  = 2560
	Height = 1440
)

const (
	MaxGrey          = 1 << 16
	Contour          = 0.1
	ContourThreshold = 0.001

	InvLargestScale = 1 / 305.1
	Offset          = -5

	WaterThreshold = 0.4
)

var rootCmd = &cobra.Command{
	Use:   "hyper-terrain",
	Short: "hyper-terrain is a fast random terrain generator",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		n := noise.Fractal{}
		n.Fill(r)

		img := image.NewNRGBA64(image.Rect(0, 0, Width, Height))

		values := make([]float64, Width*Height)
		minValue := 0.0
		maxValue := 0.0

		for x := 0; x < Width; x++ {
			px := float64(x)*InvLargestScale + Offset
			for y := 0; y < Height; y++ {
				py := float64(y)*InvLargestScale + Offset
				v := n.Cubic(px, py)
				minValue = math.Min(v, minValue)
				maxValue = math.Max(v, maxValue)
				values[x*Height+y] = v
			}
		}

		for i, v := range values {
			grey := (v - minValue) / (maxValue - minValue)
			if grey < WaterThreshold {
				grey = minValue
			} else if math.Mod(grey, Contour) < ContourThreshold {
				grey = minValue
			}
			grey = grey * float64(MaxGrey)
			img.Set(i/Height, i%Height, color.Gray16{Y: uint16(grey)})
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
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
