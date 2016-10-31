package config

const (
	ImageFormatPng string = ".png" //exported png format
	ImageFormatJpg string = ".jpg" //exported jpg format
)

//Configuration save the initial module data
type Configuration struct {
	ImageExtension   string `json:"fileFormat"`       //values jpg, png
	TileSize         int    `json:"tileSize"`         //exported image size must be power of two (256,512,1024,2408, etc)
	BaseFilename     string `json:"baseFilename"`     //base name of the generated images
	TempDir          string `json:"tempDir"`          //temporal directory
	ReturnedDataType string `json:"returnedDataType"` //exported data type (route, base64)

}
