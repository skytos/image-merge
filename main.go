package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

func main() {
	imageURLs := []string{
		"https://i.natgeofe.com/n/46b07b5e-1264-42e1-ae4b-8a021226e2d0/domestic-cat_thumb_square.jpg",
		"https://cdn.britannica.com/91/81291-050-1CDF67EB/house-mouse.jpg",
		"https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/golden-retriever-royalty-free-image-506756303-1560962726.jpg",
	}
	i, err := merge(imageURLs)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, i, &jpeg.Options{Quality: 85})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	var f *os.File
	f, err = os.Create("out.jpg")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	_, err = f.Write(buff.Bytes())
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}

func merge(imageURLs []string) (result image.Image, err error) {
	totalWidth := 0
	totalHeight := 0
	var images []image.Image
	for _, url := range imageURLs {
		img, err := downloadImage(url)
		if err != nil {
			return nil, err
		}
		width := img.Bounds().Max.X - img.Bounds().Min.X
		height := img.Bounds().Max.Y - img.Bounds().Min.Y
		images = append(images, img)
		totalWidth += width
		if totalHeight < height {
			totalHeight = height
		}
	}
	resultx := image.NewRGBA(image.Rectangle{image.Point{}, image.Point{totalWidth, totalHeight}})
	x := 0
	for _, img := range images {
		width := img.Bounds().Max.X - img.Bounds().Min.X
		height := img.Bounds().Max.Y - img.Bounds().Min.Y
		draw.Draw(
			resultx,
			image.Rectangle{image.Point{x, 0}, image.Point{x + width, height}},
			img,
			image.ZP,
			0,
		)
		x += width
	}

	return resultx, nil
}

func downloadImage(url string) (image image.Image, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Could not get: %s", url)
	}
	image, err = jpeg.Decode(response.Body)
	if err != nil {
		log.Fatal("Image decode error: ", err)
		return
	}
	return
}
