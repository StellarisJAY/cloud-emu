package dummy

import (
	"context"
	"time"
)

type Emulator struct {
	cancelFunc     context.CancelFunc
	renderCallback func(frame *Frame)
	pauseChan      chan struct{}
	resumeChan     chan struct{}
	frames         []*Frame
	i              int
	Width          int
	Height         int
	ticker         *time.Ticker
}

func MakeDummyEmulator(game []byte, frameConsumer func(frame *Frame)) (*Emulator, error) {
	e := &Emulator{
		renderCallback: frameConsumer,
		pauseChan:      make(chan struct{}),
		resumeChan:     make(chan struct{}),
		i:              0,
		ticker:         time.NewTicker(time.Second),
	}
	frames, w, h, err := parseGif(game)
	if err != nil {
		return nil, err
	}
	e.frames = frames
	e.Width = w
	e.Height = h
	return e, nil
}

const dummyEmulatorFrameTick = 5 * time.Millisecond

func (d *Emulator) Start(ctx context.Context) {
	d.ticker.Reset(dummyEmulatorFrameTick)
	for {
		select {
		case <-ctx.Done():
			return
		case <-d.ticker.C:
			frame := d.nextFrame()
			d.renderCallback(frame)
			d.ticker.Reset(frame.duration)
		}
	}
}

func (d *Emulator) nextFrame() *Frame {
	if d.i == len(d.frames) {
		d.i = 0
	}
	frame := d.frames[d.i]
	d.i++
	return frame
}

func (d *Emulator) Pause() {
	d.ticker.Stop()
}

func (d *Emulator) Resume() {
	d.ticker.Reset(dummyEmulatorFrameTick)
}

func (d *Emulator) Stop() {
	d.cancelFunc()
}
