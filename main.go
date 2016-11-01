package equitocube

//package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"image/color"
	"image/jpeg"

	"image/png"

	"encoding/base64"

	"bytes"
)

const (
	//ImageFileFormatPng exported png format
	ImageFileFormatPng string = ".png"

	//ImageFileFormatJpg exported jpg format
	ImageFileFormatJpg string = ".jpg"

	//ImageDataFormatPath return the image path
	ImageDataFormatPath string = "path"

	//ImageDataFormatBase64 return the image as base64 string
	ImageDataFormatBase64 string = "base64"
)

var (
	inWidth         float64               //image with
	inHeight        float64               //image height
	outputWidth     int                   //final image width (used only for cube calculations)
	outputHeight    int                   //final image height (used only for cube calculations)
	tileSize        = 2048                //size of exported images
	wg              sync.WaitGroup        //wait for conversion end
	images          map[string]string     //the 6 cube images in the imageDataFormat selected
	tmpDir          = "/tmp"              //temporal image directory
	imageFileFormat = ImageFileFormatJpg  //the exported image format (jpg, png)
	imageDataFormat = ImageDataFormatPath //the exported image data (path, base64)
)

/*
func main() {

	inputFilename := "/home/adolfo/Descargas/prueba_cubo/imagen5.jpg"

	imageDataFormat = ImageDataFormatBase64
	imageFileFormat = ImageFormatPng

	im := GetCubicImage("../img", "imagenprueba", inputFilename, 4096)
	fmt.Println(im["U"])
}*/

//Configuration set the init configuration of module
func Configuration(conf Config) {
	if ImageDataFormatPath != "" {
		imageDataFormat = ImageDataFormatPath
	}

	if conf.ImageFileFormat != "" {
		imageFileFormat = conf.ImageDataFormat
	}

	if conf.TempDir != "" {
		tmpDir = conf.TempDir
	}

	if !math.IsNaN(float64(conf.TileSize)) {
		tileSize = conf.TileSize
	}
}

//GetCubicImage transform equirectangular image to cubic and return the 6 image paths
func GetCubicImage(fileName, imageData string) map[string]string {
	images = make(map[string]string)

	outputWidth = tileSize * 2
	outputHeight = tileSize * 3

	cubemap, err := NewCubemap()
	if err != nil {
		fmt.Println("cubemap: ", err.Error())
	}

	w, h := cubemap.Resize(outputWidth, outputHeight)
	cubemap.TileSize.X = w
	cubemap.TileSize.Y = h
	reader, err := os.Open(imageData)

	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	im, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bo := im.Bounds()

	inWidth = float64(bo.Max.X)
	inHeight = float64(bo.Max.Y)
	ims := equirectToCubemap(im, *cubemap, fileName)
	return ims
}

//EquirectToCubemap convert an image equirectangular to cubic
func equirectToCubemap(equiImage image.Image, cubemap Cubemap, filename string) map[string]string {

	for face, faceOffset := range cubemap.FaceMap {
		wg.Add(1)
		go processCubeFace(equiImage, face, faceOffset, filename, cubemap)

	}

	wg.Wait()

	return images
}

//process a face an generate an image
func processCubeFace(equiImage image.Image, face string, faceOffset VectorArray3, name string, cubemap Cubemap) string {
	var colour = cubemap.GetFaceColor(face)

	var viewVector Vector3
	var latLong LatLong

	faceImg := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))

	for fy := 0; fy < tileSize; fy++ {
		for fx := 0; fx < tileSize; fx++ {
			var screenPos LatLong

			if fx >= faceOffset.X && fy >= faceOffset.Y {
				vx := float64(fx) / float64(cubemap.TileSize.X)
				vy := float64(fy) / float64(cubemap.TileSize.Y)
				viewVector = cubemap.ScreenToWorld(face, vx, vy)

				if viewVector.X != 0 {

					latLong = ViewToLatLon(viewVector)

					screenPos = GetScreenFromLatLong(latLong.X, latLong.Y, inWidth, inHeight)

					colour = readPixelClamped(equiImage, screenPos.X, screenPos.Y, inWidth, inHeight)
				}
				faceImg.Set(fx, fy, colour)
			}

		}
	}

	time := time.Now()
	filename := ""
	if imageDataFormat == ImageDataFormatBase64 {

		buf := new(bytes.Buffer)

		if imageFileFormat == ImageFileFormatJpg {
			var opt jpeg.Options
			opt.Quality = 80

			err := jpeg.Encode(buf, faceImg, &opt)
			if err != nil {
				return "error converting to base64"
			}
		} else {
			err := png.Encode(buf, faceImg)
			if err != nil {
				return "error converting to base64"
			}
		}

		filename = base64.URLEncoding.EncodeToString([]byte(buf.Bytes()))

	} else {

		filename = tmpDir + "/" + name + "_" + face + "_" + time.String() + imageFileFormat
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println(err)
			return "error on create file"
		}
		defer f.Close()

		if imageFileFormat == ImageFileFormatJpg {
			var opt jpeg.Options
			opt.Quality = 80
			jpeg.Encode(f, faceImg, &opt)

		} else if imageFileFormat == ImageFileFormatPng {
			png.Encode(f, faceImg)
		} else {
			return "not supported extension"
		}
	}

	images[face] = filename

	defer wg.Done()
	return filename
}

//ReadPixelClamped get the pixel color of equirectangular image to put in cube face
func readPixelClamped(img image.Image, x float64, y float64, w float64, h float64) color.RGBA64 {
	colour := color.RGBA64{255, 255, 255, 255}

	lat := int(math.Max(0, math.Min(x, w-1)))
	long := int(math.Max(0, math.Min(y, h-1)))
	r, g, b, a := img.At(lat, long).RGBA()

	colour.R = uint16(r)
	colour.G = uint16(g)
	colour.B = uint16(b)
	colour.A = uint16(a)
	return colour
}
