package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"github.com/marzhall/zanth/libraries"
	"github.com/marzhall/zanth/libraries/worldParts"

    "fmt"
    "time"
    "os"
)

var winTitle string = "Zanth"
var winWidth, winHeight int = 1920, 1080
var joysticks [16]*sdl.Joystick

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

const EventTypes (
    EntityCreated = iota
)

const EntityTypes (
    Player = iota
    Bullet
)

type EntityCreatedEvent struct {
    EventType EventTypes
    EntityType EntityTypes
    Xpos int
    Ypos int
}

type Event interface {
    GetEventType() EventTypes
}

func HandleEntityCreated(event EntityCreatedEvent) {
    eCEvent, ok := event.(EntityCreatedEvent)
    if (ok != nil) {
        panic("Got an EntityCreated event that didn't cast to EntityCreatedEvent!")
    }

    switch eCEvent.EntityType {
    case Player:
        inputState := worldParts.GenerateInputState()
        player := worldParts.Player{tiles[0], &sdl.Rect{ecEvent.Xpos, eCEvent.YPos, 100, 100}, inputState}
        libraries.AllGameEntities.Add(&player)
    }
}


func GameEventProcessor(eventBus chan Event, done chan int) {
    select {
    case event := <-eventBus:
        switch event.GetEventType() {
        case EntityCreated {
            HandleEntityCreated(event)
        }
    }
}

func main() {
    fmt.Println("start");

    // main loop
    tileSources := make([]string, 1)
    tileSources[0] = "assets/cliffy.PNG"

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

    //player := worldParts.Player{tiles[0], &sdl.Rect{0, 0, 100, 100}}
    //libraries.AllGameEntities.Add(&player)

    ebus := make(chan Event)
	sdl.JoystickEventState(sdl.ENABLE)
    running := true
    for running {
        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch t := event.(type) {
                case *sdl.QuitEvent:
                    running = false
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
                        running = false
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
            if player.GetY() > 0 {
                player.SetY(player.GetY() - 20)
            } else {
                player.SetY(0)
            }
        }
        if inputState.MovingYD {
            player.SetY(player.GetY() + 20)
        }
        if inputState.MovingXR {
            player.SetX(player.GetX() + 20)
        }
        if inputState.MovingXL {
            if player.GetX() > 0 {
                player.SetX(player.GetX() - 20)
            } else {
                player.SetX(0)
            }
        }

        renderer.Clear()
        renderer.Copy(player.GetTexture(), nil, player.GetRect())
        renderer.Present()
        time.Sleep(10)
    }
}
