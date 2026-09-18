package gifs

import (
	"EmojiPicker/app/gif/resources"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"uuid"

	_ "golang.org/x/image/webp"
)

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad status code: %s", resp.Status)
	}

	out := bytes.NewBuffer(nil)

	io.Copy(out, resp.Body)
	return out.Bytes(), nil
}

type Gif struct {
	Id       string
	Url      string
	Keywords []string
	Img      *image.RGBA

	Delay     []int
	LoopCount int
}

var AllGifs []Gif

func wrapGif(gifImage *gif.GIF) (img *image.RGBA) {
	width := gifImage.Config.Width
	height := gifImage.Config.Height
	w := len(gifImage.Image) * width
	h := height
	img = image.NewRGBA(image.Rect(0, 0, w, h))

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	var prevCanvas *image.RGBA

	for i, frame := range gifImage.Image {
		bounds := frame.Bounds()
		disposal := gifImage.Disposal[i]

		if disposal == gif.DisposalPrevious {
			prevCanvas = image.NewRGBA(canvas.Bounds())
			draw.Draw(prevCanvas, canvas.Bounds(), canvas, image.Point{}, draw.Src)
		}

		draw.Draw(canvas, bounds, frame, bounds.Min, draw.Over)

		frameSnapshot := image.NewRGBA(canvas.Bounds())
		draw.Draw(frameSnapshot, canvas.Bounds(), canvas, image.Point{}, draw.Src)

		offset := i * width
		dstBounds := image.Rect(offset, 0, offset+width, height)
		draw.Draw(img, dstBounds, frameSnapshot, img.Bounds().Min, draw.Src)

		switch disposal {
		case gif.DisposalBackground:
			draw.Draw(canvas, bounds, image.Transparent, image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			if prevCanvas != nil {
				canvas = prevCanvas
			}
		}
	}
	return
}

func Init() {
	resources.Init()
	for _, g := range resources.AllGifsData {
		if _, err := gif.DecodeConfig(bytes.NewReader(g.Data)); err != nil {
			fmt.Println(err)
			continue
		}
		gifImage, err := gif.DecodeAll(bytes.NewReader(g.Data))
		if err != nil {
			fmt.Println(err)
			continue
		}
		var g2 Gif
		g2.Id = g.Id
		g2.Url = g.Url
		g2.Img = wrapGif(gifImage)
		g2.Delay = gifImage.Delay
		g2.LoopCount = len(gifImage.Image)
		AllGifs = append(AllGifs, g2)
	}
}

func GetAllGifs() []Gif {
	return AllGifs
}

func getFileName(url string) (name string) {
	name = filepath.Base(url)
	splited := strings.Split(name, "?")
	name = splited[0]
	return
}

func SaveGif(url string) (g Gif, finalerr error) {
	data, err := downloadFile(url)
	if err != nil {
		finalerr = err
		return
	}

	if _, err := gif.DecodeConfig(bytes.NewReader(data)); err != nil {
		finalerr = err
		return
	}
	gifImage, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		finalerr = err
		return
	}

	id := uuid.NewV7().String()
	resources.StoreGif(id+".gif", url, data)
	g.Id = id
	g.Delay = gifImage.Delay
	g.LoopCount = len(gifImage.Image)
	g.Img = wrapGif(gifImage)
	g.Url = url

	AllGifs = append(AllGifs, g)
	return
}
