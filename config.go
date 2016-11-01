package equitocube

//Config save the initial module data
type Config struct {
	ImageFileFormat string `json:"imageFileFormat"` //values jpg, png
	TileSize        int    `json:"tileSize"`        //exported image size must be power of two (256,512,1024,2408, etc)
	TempDir         string `json:"tempDir"`         //temporal directory
	ImageDataFormat string `json:"imageDataFormat"` //returned image data format (saved path, base64)

}
