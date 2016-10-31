package main

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

	"github.com/adolfobushi/equirectangular-to-cubic/lib"
)

var (
	inWidth      float64 = 5064
	inHeight     float64 = 2532
	tileSize             = 1024
	outputWidth          = 4096
	outputHeight         = 6144
	wg           sync.WaitGroup
)

func main() {
	inputFilename := "//home/adolfo/Descargas/prueba_cubo/imagen4.jpg"

	GetCubicImage("/temp/", "imagenprueba", inputFilename, 2048)
}

//GetCubicImage transform equirectangular image to cubic and return the 6 image paths
func GetCubicImage(dir, name, imageData string, tilesize int) {
	tileSize = tilesize

	cubemap, err := lib.NewCubemap()
	if err != nil {
		fmt.Printf("cubemap: ", err.Error())
	}

	w, h := cubemap.Resize(tileSize, tileSize)
	cubemap.TileSize.X = w
	cubemap.TileSize.Y = h
	reader, err := os.Open(imageData)
	//var im image
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	im, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	equirectToCubemap(im, *cubemap)
	//equirectToCubemap
}

//EquirectToCubemap convert an image equirectangular to cubic
func equirectToCubemap(equiImage image.Image, cubemap lib.Cubemap) {

	cubeImage := image.NewNRGBA(image.Rect(0, 0, outputWidth, outputHeight))

	for face, faceOffset := range cubemap.FaceMap {
		wg.Add(1)
		go processFace(equiImage, *cubeImage, face, faceOffset, "imagen1", cubemap)

	}

	wg.Wait()

	time := time.Now()
	filename := "../img/imagen1_face" + time.String() + ".jpg"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	var opt jpeg.Options
	opt.Quality = 80
	jpeg.Encode(f, cubeImage, &opt)

}

//process a face an generate an image
func processFace(equiImage image.Image, cubeImage image.NRGBA, face string, faceOffset lib.VectorArray3, name string, cubemap lib.Cubemap) {
	var colour = cubemap.GetFaceColor(face)

	var viewVector lib.Vector3
	var latLong lib.LatLong

	for fy := 0; fy < cubemap.TileSize.Y; fy++ {
		for fx := 0; fx < cubemap.TileSize.X; fx++ {
			var screenPos lib.LatLong
			x := fx + (faceOffset.X * cubemap.TileSize.X)
			y := fy + (faceOffset.Y * cubemap.TileSize.Y)

			vx := float64(fx) / float64(cubemap.TileSize.X)
			vy := float64(fy) / float64(cubemap.TileSize.Y)
			viewVector = cubemap.ScreenToWorld(face, vx, vy)

			if viewVector.X != 0 {

				latLong = lib.ViewToLatLon(viewVector) //	 0.9s

				screenPos = lib.GetScreenFromLatLong(latLong.X, latLong.Y, inWidth, inHeight)

				colour = ReadPixelClamped(equiImage, screenPos.X, screenPos.Y, inWidth, inHeight)
			}

			cubeImage.Set(x, y, colour)
		}
	}

	defer wg.Done()

}

//ReadPixelClamped get the pixel color of equirectangular image to put in cube face
func ReadPixelClamped(img image.Image, x float64, y float64, w float64, h float64) color.RGBA64 {
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
