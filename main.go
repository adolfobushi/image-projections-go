package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"time"

	"image/color"
	"image/jpeg"

	"github.com/adolfobushi/equirectangular-to-cubic/lib"
)

func main() {
	inputFilename := "//home/adolfo/Descargas/prueba_cubo/site_uploadfilezilla.jpg"
	/*inputLayout := "equirect"
	sampleWidth := 5064
	sampleHeight := 2532
	sampleTime := 0

	outputFilename := "/home/adolfo/Documentos/panopo/panopo/img/img3_ok.jpg"*/

	outputWidth := 4096
	outputHeight := 6144

	cubemap, err := lib.NewCubemap()
	if err != nil {
		fmt.Printf("cubemap: ", err.Error())

	}

	w, h := cubemap.Resize(outputWidth, outputHeight)
	cubemap.TileSize.X = w
	cubemap.TileSize.Y = h
	reader, err := os.Open(inputFilename)
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
	//fmt.Printf("%+v", im.ColorModel)

	//bounds := im.Bounds()

	// Calculate a 16-bin histogram for m's red, green, blue and alpha components.
	//
	// An image's bounds do not necessarily start at (0, 0), so the two loops start
	// at bounds.Min.Y and bounds.Min.X. Looping over Y first and X second is more
	// likely to result in better memory access patterns than X first and Y second.
	/*cubeImage := image.NewNRGBA(image.Rect(0, 0, 5064, 2532))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := im.At(x, y).RGBA()

			col := color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}

			cubeImage.Set(x, y, col)

		}
	}
	f, err := os.OpenFile("ssrgb.jpg", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	var opt jpeg.Options
	opt.Quality = 80
	jpeg.Encode(f, cubeImage, &opt)
	//png.Encode(f, cubeImage)*/

	//fmt.Fprintln("image: %v", im.Bounds())
	//	$start = microtime(true);
}

//GetCubicImage transform equirectangular image to cubic and return the 6 image paths
func GetCubicImage() {

}

//EquirectToCubemap convert an image equirectangular to cubic
func equirectToCubemap(equiImage image.Image, cubemap lib.Cubemap) {

	var inWidth float64 = 5064
	var inHeight float64 = 2532

	outWidth := cubemap.GetImageWidth()
	outHeight := cubemap.GetImageHeight()

	cubeImage := image.NewNRGBA(image.Rect(0, 0, outWidth, outHeight))

	//	re-using class objects saves cpu time in massive loops
	var viewVector lib.Vector3 // := lib.VectorArray3{0, 0, lib.Vector2{0, 0}}
	var latLong lib.LatLong    // := lib.Vector2{0, 0}
	//var sphereImagePos lib.Vector2  //:= lib.Vector2{0, 0}

	//	go through each tile, convert pixel to lat long, then read
	for face, faceOffset := range cubemap.FaceMap {

		var colour = cubemap.GetFaceColor(face)
		var x int
		var y int
		for fy := 0; fy < cubemap.TileSize.Y; fy++ {
			for fx := 0; fx < cubemap.TileSize.X; fx++ {
				var screenPos lib.LatLong
				x = fx + (faceOffset.X * cubemap.TileSize.X)
				y = fy + (faceOffset.Y * cubemap.TileSize.Y)

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
	}

	time := time.Now()
	filename := "../img/" + time.String() + ".jpg"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	var opt jpeg.Options
	opt.Quality = 80
	jpeg.Encode(f, cubeImage, &opt)
	//png.Encode(f, cubeImage)
	//equiImage = cubeImage
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
