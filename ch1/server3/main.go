package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/lissajous", handlerLissajous)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handlerLissajous(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	parameters := strings.Split(url, "?")
	if len(parameters) > 1 {
		for _, p := range strings.Split(parameters[1], "&") {
			k, v := strings.Split(p, "=")[0], strings.Split(p, "=")[1]
			if k == "cycles" {
				cycles, _ := strconv.Atoi(v)
				lissajous(w, cycles)
			}
		}
	}
	lissajous(w, 5)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

var palette = []color.Color{color.White, color.Black, color.RGBA{G: 255, A: 255}}

const (
	whiteIndex = 0
	blackIndex = 1
	greenIndex = 2
)

func lissajous(out io.Writer, cycles int) {
	const (
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)

			x = math.Sin(t + 0.1)
			y = math.Sin((t+0.1)*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), greenIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: 인코딩 오류 무시
}
