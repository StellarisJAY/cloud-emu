package game

import (
	"encoding/json"
	"github.com/StellrisJAY/cloud-emu/emulator"
	"github.com/StellrisJAY/cloud-emu/gamesrv/internal/codec"
)

type emulatorRestartRequest struct {
	game         string
	gameData     []byte
	emulatorCode string
	emulatorType string
	emulatorId   int64
	gameId       int64
}

func (g *Instance) handleRestartEmulator(request *emulatorRestartRequest) ConsumerResult {
	err := g.restartEmulator(request.game, request.gameData, request.emulatorCode, request.emulatorType, request.emulatorId, request.gameId)
	return ConsumerResult{Success: err == nil, Error: err}
}

// 重启模拟器，如果模拟器类型相同，则直接重启，否则先停止当前模拟器再创建新模拟器
func (g *Instance) restartEmulator(game string, gameData []byte, emulatorCode, emulatorType string, emulatorId, gameId int64) error {
	if g.EmulatorCode == emulatorCode {
		opts, err := g.makeEmulatorOptions(g.EmulatorCode, game, gameData)
		if err != nil {
			return err
		}
		if err := g.e.Restart(opts); err != nil {
			return err
		}
		g.emulatorId = emulatorId
		g.gameId = gameId
		return nil
	}
	options, err := g.makeEmulatorOptions(emulatorCode, game, gameData)
	if err != nil {
		return err
	}
	_ = g.e.Stop()
	e, err := emulator.MakeEmulator(emulatorCode, options)
	if err != nil {
		return err
	}
	g.e = e
	g.EmulatorCode = emulatorCode
	g.emulatorId = emulatorId
	g.gameId = gameId
	width, height := g.e.OutputResolution()
	// 模拟器画面规格可能不同，重新创建视频编码器
	videoEncoder, err := codec.NewVideoEncoder("vp8", width, height)
	if err != nil {
		return err
	}
	g.videoEncoder.Close()
	g.videoEncoder = videoEncoder
	if err := g.e.Start(); err != nil {
		return err
	}
	return nil
}

func (g *Instance) RestartEmulator(game string, gameData []byte, emulatorCode, emulatorType string, emulatorId, gameId int64) error {
	ch := make(chan ConsumerResult)
	request := &emulatorRestartRequest{game, gameData, emulatorCode, emulatorType, emulatorId, gameId}
	g.messageChan <- &Message{Type: MsgRestartEmulator, Data: request, resultChan: ch}
	result := <-ch
	close(ch)
	if result.Success {
		g.onRestartSuccess(emulatorId, gameId, emulatorCode, emulatorType)
		return nil
	} else {
		return result.Error
	}
}

func (g *Instance) onRestartSuccess(emulatorId, gameId int64, emulatorCode, emulatorType string) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	content := &struct {
		EmulatorId   int64
		GameId       int64
		EmulatorCode string
		EmulatorType string
	}{
		EmulatorId:   emulatorId,
		GameId:       gameId,
		EmulatorCode: emulatorCode,
		EmulatorType: emulatorType,
	}
	msg := &Message{Type: MsgEmulatorRestart, Data: content}
	raw, _ := json.Marshal(msg)
	for _, conn := range g.connections {
		_ = conn.dataChannel.Send(raw)
	}
}

func (g *Instance) makeEmulatorOptions(emulatorName string, game string, gameData []byte) (emulator.IEmulatorOptions, error) {
	switch emulatorName {
	case emulator.CodeNesgo:
		return emulator.MakeNesEmulatorOptions(game, gameData, g.audioSampleRate, g.audioSampleChan, func(frame emulator.IFrame) {
			g.RenderCallback(frame, nil)
		}), nil
	case emulator.CodeChip8, emulator.CodeDummy, emulator.CodeGoboy:
		return emulator.MakeBaseEmulatorOptions(game, gameData, g.audioSampleRate, g.audioSampleChan, func(frame emulator.IFrame) {
			g.RenderCallback(frame, nil)
		})
	default:
		return nil, emulator.ErrorEmulatorNotSupported
	}
}
