package lib

import (
	"errors"
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

const (
	pi = 3.14159265358979323846264338327950288
)

//23ULFRBD
//Cubemap save cubemap data
type Cubemap struct {
	Ratio          Vector2
	TileSize       Vector2
	TileMap        [2][3]string
	FaceMap        map[string]VectorArray3
	SquareTileSize int
}

//Vector2 save vector2 data
type Vector2 struct {
	X int
	Y int
}

type LatLong struct {
	X float64
	Y float64
}

//Vector3 save vector3 data
type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type VectorArray3 struct {
	X int
	Y int
	Z Vector2
}

//DegreesToRadians convert degrees to radians
func DegreesToRadians(deg float64) float64 {
	return deg * pi / 180
}

//RadiansToDegrees convert radian angles to degrees
func RadiansToDegrees(rad float64) float64 {
	return rad * 180.0 / pi
}

//	LatLonToView
/*
func VectorFromCoordsRad($latlon)
{
    //	http://en.wikipedia.org/wiki/N-vector#Converting_latitude.2Flongitude_to_n-vector
    latitude := $latlon->x;
    longitude = $latlon->y;
    las = sin($latitude);
    lac = cos($latitude);
    los = sin($longitude);
    loc = cos($longitude);

    result = new Vector3( $los * $lac, $las, $loc * $lac );
    //assert(fabsf(result.Length() - 1.0f) < 0.01f);

    return $result;
}*/

//NewCubemap creates a new Cubemap
func NewCubemap() (*Cubemap, error) {
	var layout = [6]string{"U", "L", "F", "R", "B", "D"}

	with := 2
	height := 3

	c := new(Cubemap)
	c.FaceMap = make(map[string]VectorArray3)

	ratio, squareTilesize, err := c.GetRatio(with, height)
	if err != nil {
		return c, err
	}

	c.Ratio = ratio

	c.TileSize = Vector2{squareTilesize, squareTilesize}

	//	see if layout fits
	ratioTileCount := c.Ratio.X * c.Ratio.Y
	layoutTileCount := len(layout)

	if ratioTileCount != layoutTileCount {
		//	scale if posss
		remainder := ratioTileCount % layoutTileCount
		if remainder != 0 {
			return c, errors.New("layout: Layout doesn't fit in ratio")
		}
		//	ratio can only scale upwards...
		if ratioTileCount > layoutTileCount {
			return c, errors.New("layout: Layout doesn't have enough tiles to fit in ratio")
		}

		scale := layoutTileCount / ratioTileCount
		c.Ratio.X *= scale
		c.Ratio.Y *= scale
		c.TileSize.X /= scale
		c.TileSize.Y /= scale
	}

	//	make up 2d map
	for x := 0; x < c.GetTileWidth(); x++ {
		for y := 0; y < c.GetTileHeight(); y++ {
			i := y*c.GetTileWidth() + x
			face := layout[i]
			realFace := c.GetRealFace(face)

			//	gr: turn this whole thing into a matrix!
			matrix := c.GetFaceMatrix(face)

			c.TileMap[x][y] = face
			c.FaceMap[realFace] = VectorArray3{x, y, matrix}

		}
	}

	return c, nil
}

//Resize a Cubemap
func (c Cubemap) Resize(width, height int) (int, int) {
	w := width / c.Ratio.X
	h := height / c.Ratio.Y
	return w, h
}

//GetRatio get the ratio of an image
func (c Cubemap) GetRatio(width int, height int) (Vector2, int, error) {
	var vector = Vector2{}
	if width <= 0 || height <= 0 {
		return Vector2{}, 0, errors.New("ratio: with or heigh equal to zero")
	}
	//	square
	if width == height {

		return Vector2{1, 1}, width, nil
	}

	if width > height {
		remainder := width % height
		if remainder == 0 {

			return Vector2{width / height, 1}, height, nil
		}

		vector.X = width / remainder
		vector.Y = height / remainder
		return vector, remainder, nil
	}

	ratio, squareTilesize, err := c.GetRatio(height, width)
	if err != nil {
		return Vector2{}, 0, err
	}

	return Vector2{ratio.Y, ratio.X}, squareTilesize, nil
}

//GetFaceMatrix get the face Matrix
func (c Cubemap) GetFaceMatrix(face string) Vector2 {
	if face == "Z" {
		return Vector2{-1, -1}
	}

	return Vector2{1, 1}

}

//GetFlipFace get the flipped face
func (c Cubemap) GetFlipFace(face string) string {
	if face == "B" {
		return "Z"
	}
	return ""
}

//GetRealFace get the real face
func (c Cubemap) GetRealFace(face string) string {
	if face == "Z" {
		return "B"
	}
	return face
}

/*
func (c Cubemap) IsValid() bool {
	return if(math.IsNaN(c.Ratio.X))
}*/

//GetTileWidth get the tile width
func (c Cubemap) GetTileWidth() int {
	return c.Ratio.X
}

//GetTileHeight get the tile Height
func (c Cubemap) GetTileHeight() int {
	return c.Ratio.Y
}

//GetImageWidth get the image width
func (c Cubemap) GetImageWidth() int {
	return c.Ratio.X * c.TileSize.X
}

//GetImageHeight get the image height
func (c Cubemap) GetImageHeight() int {
	return c.Ratio.Y * c.TileSize.Y
}

func (c Cubemap) getSquareTileSize() int {
	return c.SquareTileSize
}

//ScreenToWorld get the screen position
func (c Cubemap) ScreenToWorld(face string, screenPosX float64, screenPosY float64) Vector3 {
	//	0..1 -> -1..1
	screenPosX *= 2.0
	screenPosY *= 2.0
	screenPosX -= 1.0
	screenPosY -= 1.0
	vector := Vector3{0, 0, 0}
	switch face {
	case "L":
		vector.X = -1
		vector.Y = -screenPosY
		vector.Z = screenPosX
		break
	case "R":
		vector.X = 1
		vector.Y = -screenPosY
		vector.Z = -screenPosX
		break
	case "U":
		vector.X = -screenPosX
		vector.Y = 1
		vector.Z = -screenPosY
		break
	case "D":
		vector.X = -screenPosX
		vector.Y = -1
		vector.Z = screenPosY
		break
	case "F":
		vector.X = screenPosX
		vector.Y = -screenPosY
		vector.Z = 1
		break
	case "B":
		vector.X = -screenPosX
		vector.Y = -screenPosY
		vector.Z = -1
		break
	}

	return vector
}

func (c Cubemap) GetFaceColor(face string) color.RGBA {
	p := color.RGBA{
		255, 255, 255, 255,
	}
	switch face {
	case "U":
		p.R = 255
		p.G = 0
		p.B = 0
		break
	case "L":
		p.R = 0
		p.G = 255
		p.B = 0
		break
	case "F":
		p.R = 0
		p.G = 0
		p.B = 255
		break
	case "R":
		p.R = 255
		p.G = 255
		p.B = 0
		break
	case "B":
		p.R = 0
		p.G = 255
		p.B = 255
		break
	case "D":
		p.R = 255
		p.G = 0
		p.B = 255
		break
	default:
		p.R = 255
		p.G = 255
		p.B = 255
		break
	}
	return p
}

func GetRgb(r, g, b float64) string {
	c := colorful.Color{r, g, b}
	return c.Hex()
}

func ViewToLatLon(view3 Vector3) LatLong {
	var latLong = LatLong{0, 0}

	x := float64(view3.X)
	y := float64(view3.Y)
	z := float64(view3.Z)

	lat := math.Atan2(x, z)

	//	normalize y
	xz := math.Sqrt((x * x) + (z * z))
	normy := y / math.Sqrt((y*y)+(xz*xz))
	lon := math.Asin(normy)
	//$lon = atan2( $y, $xz );

	//	stretch longitude...
	lon *= 2.0

	latLong.X = lat
	latLong.Y = lon
	return latLong
}

func GetScreenFromLatLong(lat float64, lon float64, width float64, height float64) LatLong {
	var screenPos = LatLong{0, 0}
	//	-pi...pi -> -1...1
	lat /= pi
	lon /= pi

	//	-1..1 -> 0..2
	lat += 1.0
	lon += 1.0

	//	0..2 -> 0..1
	lat /= 2.0
	lon /= 2.0

	lon = 1.0 - lon
	lat *= width
	lon *= height

	screenPos.X = lat
	screenPos.Y = lon
	//fmt.Printf("lat: ", screenPos.X)
	return screenPos
}
