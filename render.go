package vol

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

type Window struct {
	Lo, Hi float64
	Color  color.Color
}

func Render(images []image.Image, windows []Window, matrix Matrix, width, height, samples int, zScale float64) image.Image {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	w := images[0].Bounds().Size().X
	h := images[0].Bounds().Size().Y
	d := len(images)
	buffers := make([][]float64, len(windows))
	for i := range buffers {
		buffers[i] = make([]float64, width*height)
	}
	for z, im := range images {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c, _, _, _ := im.At(x, y).RGBA()
				if c == 0 {
					continue
				}
				for i, window := range windows {
					f := float64(c) / 65535
					f = (f - window.Lo) / (window.Hi - window.Lo)
					if f <= 0 {
						continue
					}
					if f > 1 {
						f = 1
					}
					for s := 0; s < samples; s++ {
						vx := (float64(x)+rnd.Float64())/float64(w-1)*2 - 1
						vy := (float64(y)+rnd.Float64())/float64(h-1)*2 - 1
						vz := (float64(z)+rnd.Float64())/float64(d-1)*2 - 1
						vz *= zScale
						v := Vector{vx, vy, vz}
						v = matrix.MulPositionW(v)
						px := int((v.X + 1) / 2 * float64(width))
						py := int((v.Y + 1) / 2 * float64(height))
						if px < 0 || py < 0 || px >= width || py >= height {
							continue
						}
						buffers[i][py*width+px] += f
					}
				}
			}
		}
	}
	for _, buffer := range buffers {
		hi := buffer[0]
		for _, x := range buffer {
			hi = math.Max(hi, x)
		}
		for i, x := range buffer {
			buffer[i] = x / hi
		}
	}
	im := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b float64
			for i, buffer := range buffers {
				wr, wg, wb, _ := windows[i].Color.RGBA()
				ww := buffer[y*width+x]
				r += ww * float64(wr) / 65535
				g += ww * float64(wg) / 65535
				b += ww * float64(wb) / 65535
			}
			if r > 1 {
				r = 1
			}
			if g > 1 {
				g = 1
			}
			if b > 1 {
				b = 1
			}
			// r /= float64(len(windows))
			// g /= float64(len(windows))
			// b /= float64(len(windows))
			c := color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
			im.SetRGBA(x, y, c)
		}
	}
	return im
}
