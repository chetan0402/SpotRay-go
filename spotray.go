package spotraygo

import (
	"errors"
	"strings"
	"time"

	"github.com/flopp/go-findfont"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type debounce int

const (
	FIRST debounce = iota
	REPEAT
)

type keyState struct {
	time int64
	d    debounce
}

type Item struct {
	Key   string
	Value string
}

type State struct {
	Options  []Item
	Keybinds map[int32]func(*App)
	Font     string
}

type App struct {
	state          *State
	lastPressed    map[int32]keyState
	searchTerm     string
	prevSearchTerm string
	selectedIndex  int
	filtered       []Item
	font           rl.Font
}

func NewApp() *App {
	a := &App{
		state:          nil,
		lastPressed:    map[int32]keyState{},
		searchTerm:     "",
		prevSearchTerm: " ",
		selectedIndex:  0,
		filtered:       nil,
	}
	rl.SetConfigFlags(rl.FlagWindowUndecorated | rl.FlagWindowTransparent)
	rl.InitWindow(1, 1, "")
	rl.SetWindowSize(rl.GetMonitorWidth(0), rl.GetMonitorHeight(0)-26)
	rl.SetWindowPosition(0, 26)
	rl.SetTargetFPS(60)
	return a
}

func (a *App) Start() error {
	if a.state == nil {
		return errors.New("state is nil")
	}
	for !rl.WindowShouldClose() {
		key := rl.GetCharPressed()
		for key > 0 {
			a.searchTerm += string(byte(key))
			key = rl.GetCharPressed()
		}

		if a.prevSearchTerm != a.searchTerm {
			a.filtered = []Item{}
			for _, v := range a.state.Options {
				if strings.Contains(strings.ToLower(v.Value), strings.ToLower(a.searchTerm)) {
					a.filtered = append(a.filtered, v)
				}
			}
		}

		for i, v := range a.state.Keybinds {
			rn := time.Now().UnixMilli()
			if rl.IsKeyPressed(i) {
				a.lastPressed[i] = keyState{time: rn, d: FIRST}
				v(a)
			} else if rl.IsKeyDown(i) {
				if (a.lastPressed[i].d == FIRST && a.lastPressed[i].time+500 < rn) || (a.lastPressed[i].d == REPEAT && a.lastPressed[i].time+50 < rn) {
					a.lastPressed[i] = keyState{time: rn, d: REPEAT}
					v(a)
				}
			}
		}

		if len(a.filtered) > 0 {
			a.selectedIndex = (a.selectedIndex%len(a.filtered) + len(a.filtered)) % len(a.filtered)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Color{R: 30, G: 45, B: 50, A: 200})
		rl.DrawTextEx(a.font, a.searchTerm, rl.Vector2{X: 5, Y: 5}, 20, 0, rl.White)
		for i, v := range a.filtered {
			rl.DrawTextEx(a.font, v.Value, rl.Vector2{X: 5, Y: float32(5 + (i+1)*20)}, 20, 0, rl.White)
		}

		if len(a.filtered) > 0 {
			rl.DrawRectangle(5, int32(5+20*(a.selectedIndex+1)), int32(rl.GetMonitorWidth(0)), 20, rl.ColorAlpha(rl.White, 0.2))
		}

		rl.EndDrawing()
	}
	return nil
}

func (a *App) SetState(s *State) error {
	a.state = s
	a.prevSearchTerm = a.searchTerm + " "

	s.Keybinds[rl.KeyBackspace] = func(a *App) {
		if len(a.searchTerm) > 0 {
			a.searchTerm = a.searchTerm[:len(a.searchTerm)-1]
		}
	}
	s.Keybinds[rl.KeyDown] = func(a *App) { a.selectedIndex++ }
	s.Keybinds[rl.KeyUp] = func(a *App) { a.selectedIndex-- }

	for i := range s.Keybinds {
		a.lastPressed[i] = keyState{time.Now().UnixMilli(), FIRST}
	}

	if s.Font != "" {
		path, err := findfont.Find(s.Font)
		if err != nil {
			return err
		}
		a.font = rl.LoadFontEx(path, 20, nil)
	}
	return nil
}
