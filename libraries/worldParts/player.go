package worldParts

import "github.com/veandco/go-sdl2/sdl"

type InputState struct {
	MovingXL, MovingXR, MovingYD, MovingYU, Running bool
}

func GenerateInputState() InputState {
	inputState := InputState{
		MovingXL: false,
		MovingXR: false,
		MovingYD: false,
		MovingYU: false,
		Running: true }

	return inputState;
}

type Player struct {
    Sprite *sdl.Texture
    Rect *sdl.Rect
    InputsState InputState
}

func (player Player) GetX() int {
	return int(player.Rect.X)
}

func (player Player) GetY() int {
	return int(player.Rect.Y)
}

func (player Player) GetHeight() int {
	return int(player.Rect.H)
}

func (player Player) GetWidth() int {
	return int(player.Rect.W)
}

func (player Player) GetTexture() *sdl.Texture {
	return player.Sprite
}

func (player Player) GetRect() *sdl.Rect {
	return player.Rect
}

func (player *Player) SetX(x int) {
	player.Rect.X = int32(x)
}

func (player *Player) SetY(y int) {
	player.Rect.Y = int32(y)
}

func (player *Player) SetHeight(h int) {
	player.Rect.H = int32(h)
}

func (player *Player) SetWidth(w int) {
	player.Rect.W = int32(w)
}

func (player *Player) SetTexture(t *sdl.Texture) {
	player.Sprite = t
}

