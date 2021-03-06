//package main

package equitocube

//Config save the initial module data
type Config struct {
	ImageFileFormat      string `json:"imageFileFormat"`      //values jpg, png
	ImageCompresion      int    `json:"imageCompresion"`      //values number 1 to 100 only for jpg format
	TileSize             int    `json:"tileSize"`             //exported image size must be power of two (256,512,1024,2408, etc)
	TempDir              string `json:"tempDir"`              //temporal directory
	ImageDataFormat      string `json:"imageDataFormat"`      //returned image data format (saved path, base64)
	InputImageDataFormat string `json:"inputImageDataFormat"` //the input image data format (path, base64)

}
