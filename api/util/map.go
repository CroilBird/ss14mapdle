package util

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path"
	"ss14mapdle/config"
	"strconv"
)

const ZoomMultiplier float32 = 1.2
const StartWidth = 100
const StartHeight = 100
const MaxZoomLevel = 6

func readImage(path string) (*image.Image, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		deferErr := fd.Close()
		if deferErr != nil {
			panic("this should not happen bye")
		}
	}()

	img, err := png.Decode(fd)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

// writeImage writes an Image back to the disk.
func writeImage(img image.Image, path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		deferErr := fd.Close()
		if deferErr != nil {
			panic("this should not happen bye")
		}
	}()

	return png.Encode(fd, img)
}

func cropImage(img image.Image, crop image.Rectangle) (*image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	subImg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	newImage := subImg.SubImage(crop)

	return &newImage, nil
}

func ReadMap(mapPath string, index int) (*image.Image, error) {
	imageFullPath := mapPath + strconv.Itoa(index) + ".png"

	mapBasePath, err := config.GetConfig(config.EnvMapBasePath)

	if err != nil {
		return nil, err
	}

	fullPath := path.Join(mapBasePath, imageFullPath)

	return readImage(fullPath)
}

func GetMapPathAtLevel(mapName string, mapPath string, index int, x int, y int, z int) (*string, error) {
	imageFullPath := mapPath + strconv.Itoa(index) + ".png"
	imageAtLevelName := fmt.Sprintf("%s-%d-%d-%d.png", mapName, x, y, z)

	// check if in cache
	mapBasePath, err := config.GetConfig(config.EnvMapBasePath)

	if err != nil {
		return nil, err
	}

	fullPath := path.Join(mapBasePath, imageFullPath)

	if z == -1 {
		z = 10
	}

	cacheBase := path.Join(mapBasePath, "cache")
	cachePath := path.Join(cacheBase, imageAtLevelName)

	_, err = os.Stat(cachePath)

	if os.IsExist(err) {
		return &cachePath, nil
	}

	img, err := readImage(fullPath)

	if err != nil {
		return nil, err
	}

	width := StartWidth * int(ZoomMultiplier*float32(z))
	height := StartHeight * int(ZoomMultiplier*float32(z))

	r := image.Rectangle{
		Min: image.Point{
			X: x - width/2,
			Y: y - height/2,
		},
		Max: image.Point{
			X: x + width/2,
			Y: y + height/2,
		},
	}

	croppedImage, err := cropImage(*img, r)

	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cacheBase, os.ModePerm)

	if err != nil {
		return nil, err
	}

	err = writeImage(*croppedImage, cachePath)

	if err != nil {
		return nil, err
	}

	return &cachePath, nil
}
