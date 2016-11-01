package config

const (
	//ImageFormatPng exported png format
	ImageFileFormatPng string = ".png"

	//ImageFormatJpg exported jpg format
	ImageFileFormatJpg string = ".jpg"

	//ImageDataFormatPath return the image path
	ImageDataFormatPath string = "path"

	//ImageDataFormatBase64 return the image as base64 string
	ImageDataFormatBase64 string = "base64"
)

//Configuration save the initial module data
type Configuration struct {
	ImageFileFormat string `json:"imageFileFormat"` //values jpg, png
	TileSize        int    `json:"tileSize"`        //exported image size must be power of two (256,512,1024,2408, etc)
	TempDir         string `json:"tempDir"`         //temporal directory
	ImageDataFormat string `json:"imageDataFormat"` //returned image data format (saved path, base64)

}
