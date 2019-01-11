package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"github.com/marzhall/zanth/libraries"

    "fmt"
    "time"
    "os"
)

var winTitle string = "Zanth"
var winWidth, winHeight int = 1920, 1080
var joysticks [16]*sdl.Joystick

type Player struct {
    Xpos, Ypos int
    Height, Width int
    Sprite *sdl.Texture
    Rect *sdl.Rect
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

func loadTiles(tileSources *[]string, renderer *sdl.Renderer) []*sdl.Texture {
	tileSet := make([]*sdl.Texture, len(*tileSources))
	for i,v := range *tileSources {
		image, err := img.Load(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load Image: %s\n", err)
			return nil
		}
		defer image.Free()

		texture, err := renderer.CreateTextureFromSurface(image)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
			return nil
		}

		tileSet[i] = texture
	}

	return tileSet
}

func main() {
    fmt.Println("start");
    xpos, ypos := int32(0), int32(0)

    // main loop
    tileSources := make([]string, 1)
    tileSources[0] = "assets/cliffy.PNG"

	inputState := libraries.GenerateInputState()
	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_FULLSCREEN_DESKTOP | sdl.WINDOW_INPUT_GRABBED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return
	}
	defer renderer.Destroy()

    tiles := loadTiles(&tileSources, renderer)
    if tiles == nil {
        return
    }

	sdl.JoystickEventState(sdl.ENABLE)
    for inputState.Running {
        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch t := event.(type) {
                case *sdl.QuitEvent:
                    inputState.Running = false
                case *sdl.KeyDownEvent:
					if (t.Keysym.Sym == sdl.K_UP && t.Repeat <= 0) {
						inputState.MovingYU = true
					}
					if (t.Keysym.Sym == sdl.K_DOWN && t.Repeat <= 0) {
						inputState.MovingYD = true
					}
					if (t.Keysym.Sym == sdl.K_RIGHT && t.Repeat <= 0) {
						inputState.MovingXR = true
					}
					if (t.Keysym.Sym == sdl.K_LEFT && t.Repeat <= 0) {
						inputState.MovingXL = true
                    }
                    if (t.Keysym.Sym == sdl.K_ESCAPE && t.Repeat <= 0) {
                        inputState.Running = false
                    }
                case *sdl.KeyUpEvent:
					if (t.Keysym.Sym == sdl.K_UP && t.Repeat <= 0) {
						inputState.MovingYU = false
					}
					if (t.Keysym.Sym == sdl.K_DOWN && t.Repeat <= 0) {
						inputState.MovingYD = false
					}
					if (t.Keysym.Sym == sdl.K_RIGHT && t.Repeat <= 0) {
						inputState.MovingXR = false
					}
					if (t.Keysym.Sym == sdl.K_LEFT && t.Repeat <= 0) {
						inputState.MovingXL = false
					}
                case *sdl.JoyAxisEvent:
                    fmt.Printf("[%d ms] JoyAxis\ttype:%d\twhich:%c\taxis:%d\tvalue:%d\n",
                        t.Timestamp, t.Type, t.Which, t.Axis, t.Value)
                case *sdl.JoyBallEvent:
                    fmt.Printf("[%d ms] JoyBall\ttype:%d\twhich:%d\tball:%d\txrel:%d\tyrel:%d\n",
                        t.Timestamp, t.Type, t.Which, t.Ball, t.XRel, t.YRel)
                case *sdl.JoyButtonEvent:
                    fmt.Printf("[%d ms] JoyButton\ttype:%d\twhich:%d\tbutton:%d\tstate:%d\n",
                        t.Timestamp, t.Type, t.Which, t.Button, t.State)
                case *sdl.JoyHatEvent:
                    fmt.Printf("[%d ms] JoyHat\ttype:%d\twhich:%d\that:%d\tvalue:%d\n",
                        t.Timestamp, t.Type, t.Which, t.Hat, t.Value)
                case *sdl.JoyDeviceEvent:
                    if t.Type == sdl.JOYDEVICEADDED {
                        joysticks[int(t.Which)] = sdl.JoystickOpen(t.Which)
                        if joysticks[int(t.Which)] != nil {
                            fmt.Printf("Joystick %d connected\n", t.Which)
                        }
                    } else if t.Type == sdl.JOYDEVICEREMOVED {
                        if joystick := joysticks[int(t.Which)]; joystick != nil {
                            joystick.Close()
                        }
                        fmt.Printf("Joystick %d disconnected\n", t.Which)
                    }
                default:
                    fmt.Printf("Some event\n")
            }
        }

        if inputState.MovingYU {
            if ypos > 0 {
                ypos -= 20*1
            } else {
                ypos = 0
            }
        }
        if inputState.MovingYD {
            ypos += 20*1
        }
        if inputState.MovingXR {
            xpos += 20*1
        }
        if inputState.MovingXL {
            if xpos > 0 {
                xpos -= 20*1
            } else {
                xpos = 0
            }
        }

        renderer.Clear()
        renderer.Copy(tiles[0], nil, &sdl.Rect{xpos,ypos, int32(100), int32(100)})
        renderer.Present()
        time.Sleep(10)
    }
}
