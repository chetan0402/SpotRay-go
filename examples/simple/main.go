package main

import spotraygo "github.com/chetan0402/SpotRay-go"

func main() {
	a := spotraygo.NewApp()
	a.SetState(&spotraygo.State{
		Options:  []spotraygo.Item{spotraygo.Item{Key: "key", Value: "value"}},
		Keybinds: map[int32]func(*spotraygo.App){},
	})
	a.Start()
}
