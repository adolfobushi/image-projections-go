package lib

import "math"

//Cubemap save cubemap data
type Cubemap struct {
	Ratio    Vector2
	TileSize Vector2
	TileMap  [][]int
	FaceMap  Vector2
}

//Vector2 save vector2 data
type Vector2 struct {
	X int
	Y int
}

//Vector3 save vector3 data
type Vector3 struct {
	X int
	Y int
	Z int
}

//DegreesToRadians convert degrees to radians
func DegreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

//RadiansToDegrees convert radian angles to degrees
func RadiansToDegrees(rad float64) float64 {
	return rad * 180.0 / math.Pi
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

//Resize a Cubemap
func (c Cubemap) Resize(width, height int) {
	c.TileSize.X = width / c.Ratio.X
	c.TileSize.Y = height / c.Ratio.Y

}

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

func (c Cubemap) GetFaceOffset() {

}
