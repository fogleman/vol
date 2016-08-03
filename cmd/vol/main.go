package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"path"
	"runtime"

	"github.com/fogleman/vol"
)

const (
	width   = 512
	height  = 512
	samples = 16
	zScale  = 0.625
	fovy    = 65
	zNear   = 0.1
	zFar    = 10
)

var (
	eye     = vol.V(1.5, 1.5, -0.25)
	center  = vol.V(0, 0, -0.1)
	up      = vol.V(0, 0, -1)
	windows = []vol.Window{
		vol.Window{0.10, 0.11, color.RGBA{255, 0, 0, 255}},
		vol.Window{0.30, 0.31, color.RGBA{0, 255, 0, 255}},
		vol.Window{0.51, 0.51, color.RGBA{0, 0, 255, 255}},
	}
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: vol DIRECTORY")
		return
	}
	dirname := args[0]
	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	var images []image.Image
	for _, info := range infos {
		filename := path.Join(dirname, info.Name())
		im, err := vol.LoadImage(filename)
		if err != nil {
			panic(err)
		}
		images = append(images, im)
	}
	aspect := float64(width) / float64(height)
	jobs := make(chan int, 1024)
	done := make(chan int)
	ncpu := runtime.NumCPU()
	for i := 0; i < ncpu; i++ {
		go func() {
			for a := range jobs {
				fmt.Println(a)
				x := math.Cos(float64(a)*math.Pi/180) * 1.5
				y := math.Sin(float64(a)*math.Pi/180) * 1.5
				z := -0.25
				eye = vol.V(x, y, z)
				m := vol.LookAt(eye, center, up)
				m = m.Perspective(fovy, aspect, zNear, zFar)
				im := vol.Render(images, windows, m, width, height, samples, zScale)
				vol.SavePNG(fmt.Sprintf("out%03d.png", a), im)
				done <- 1
			}
		}()
	}
	for i := 0; i < 360; i += 5 {
		jobs <- i
	}
	for i := 0; i < 360; i += 5 {
		<-done
	}
}
