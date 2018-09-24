package ansi

import (
	"fmt"
	"time"
)

// cleanup clears the line and redisplays the cursor
func cleanup() {
	fmt.Print(EraseLine(2))
	fmt.Print(ShowCursor())
}

// LoaderControl holds internal data for the loader
type LoaderControl struct {
	anim   []string     // The animation sprite list
	value  string       // The label value
	cursor int          // The current position of the animation
	max    int          // The length of the animation list
	ticker *time.Ticker // Pointer to the ticker instance
	done   chan bool    // Channel used to stop the loader
}

// Start will start the loader
func (loader *LoaderControl) Start() {
	go func() {
		for {
			select {
			case <-loader.done:
				return
			case <-loader.ticker.C:
				fmt.Print(HideCursor())
				fmt.Print(EraseLine(2))
				fmt.Print(CursorStart(1))

				if loader.cursor >= loader.max {
					loader.cursor = 0
				}

				fmt.Print(loader.anim[loader.cursor] + " " + loader.value)

				loader.cursor++
			}
		}
	}()
}

// Stop will stop the loader
func (loader *LoaderControl) Stop() {
	loader.done <- true
	cleanup()
}

// SetValue will set a new value for the label
func (loader *LoaderControl) SetValue(value string) {
	loader.value = value
}

// Loader creates a new loader with the given `anim` that animates at the specified `speed`
func Loader(anim []string, speed int) *LoaderControl {
	return &LoaderControl{
		cursor: 0,
		anim:   anim,
		value:  "",
		max:    len(anim),
		ticker: time.NewTicker(time.Millisecond * time.Duration(speed)),
		done:   make(chan bool),
	}
}
