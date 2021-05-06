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

	"willbeason/hyper-terrain/pkg/noise"

	"github.com/spf13/cobra"
)

const (
	Width = 2560
	Height = 1440
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := rand.New(rand.NewSource(time.Now().UnixNano()*0))

		n := noise.Value{}
		n.Fill(r)

		img := image.NewNRGBA64(image.Rect(0, 0, Width, Height))

		values := make([]float64, Width*Height)
		minValue := 0.0
		maxValue := 0.0

		for x := 0; x < Width; x++ {
			px := float64(x) / 105.1
			for y := 0; y < Height; y++ {
				py := float64(y) / 105.1
				v := n.CubicFloat(px, py)
				minValue = math.Min(v, minValue)
				maxValue = math.Max(v, maxValue)
				values[x*Height+y] = v
			}
		}

		for i, v := range values {
			grey := (v - minValue) / (maxValue - minValue)
			grey = grey * float64(int(1) << 16)
			img.Set(i / Height, i % Height, color.Gray16{Y: uint16(grey)})
		}

		out, err := os.Create("out.png")
		if err != nil {
			return err
		}
		defer out.Close()

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
