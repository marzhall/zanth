package libraries

import "github.com/veandco/go-sdl2/sdl"

type Entity interface {
	GetX() int
	GetY() int
	GetHeight() int
	GetWidth() int
	GetTexture() *sdl.Texture
	GetRect() *sdl.Rect

	SetX(x int)
	SetY(y int)
	SetHeight(h int)
	SetWidth(w int)
	SetTexture(t *sdl.Texture)
}
