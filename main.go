package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"image/color"
	_ "image/jpeg"
	"image/png"

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
	fmt.Println("cube: %v", cubemap)

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
	EquirectToCubemap(im, *cubemap)
	//fmt.Fprintln("image: %v", im.Bounds())
	//	$start = microtime(true);
}

func EquirectToCubemap(equiImage image.Image, cubemap lib.Cubemap) {

	var inWidth float64 = 5064
	var inHeight float64 = 2532

	outWidth := cubemap.GetImageWidth()
	outHeight := cubemap.GetImageHeight()

	cubeImage := image.NewNRGBA(image.Rect(0, 0, outWidth, outHeight))

	//	re-using class objects saves cpu time in massive loops
	//var viewVector lib.VectorArray3 // := lib.VectorArray3{0, 0, lib.Vector2{0, 0}}
	//var latLon lib.Vector2          // := lib.Vector2{0, 0}
	//var sphereImagePos lib.Vector2  //:= lib.Vector2{0, 0}

	//	go through each tile, convert pixel to lat long, then read
	for face, faceOffset := range cubemap.FaceMap {

		var colour color.RGBA = cubemap.GetFaceColor(face)
		var x int
		var y int
		for fy := 0; fy < cubemap.TileSize.Y; fy++ {
			for fx := 0; fx < cubemap.TileSize.X; fx++ {
				var screenPos lib.LatLong
				x = fx + (faceOffset.X * cubemap.TileSize.X)
				y = fy + (faceOffset.Y * cubemap.TileSize.Y)

				vx := float64(fx) / float64(cubemap.TileSize.X)
				vy := float64(fy) / float64(cubemap.TileSize.Y)
				viewVector := cubemap.ScreenToWorld(face, vx, vy)

				if viewVector.X != 0 {

					latLong := lib.ViewToLatLon(viewVector) //	 0.9s

					screenPos = lib.GetScreenFromLatLong(latLong.X, latLong.Y, inWidth, inHeight) //	0.41

					colour = ReadPixelClamped(equiImage, screenPos.X, screenPos.Y, inWidth, inHeight) //	 1.77
					//fmt.Println("color: ", colour)
				}

				cubeImage.Set(x, y, colour)
			}
		}
	}

	f, err := os.OpenFile("rgb.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	png.Encode(f, cubeImage)
	//equiImage = cubeImage
}

//	on 512x512 providing w/h saved 0.2sec
func ReadPixelClamped(img image.Image, x float64, y float64, w float64, h float64) color.RGBA {
	colour := color.RGBA{255, 255, 255, 255}

	lat := int(math.Max(0, math.Min(x, w-1)))
	long := int(math.Max(0, math.Min(y, h-1)))
	//fmt.Printf("lat: ", lat, long)
	r, g, b, a := img.At(lat, long).RGBA()

	colour.R = uint8(r)
	colour.G = uint8(g)
	colour.B = uint8(b)
	colour.A = uint8(a)
	return colour
}
