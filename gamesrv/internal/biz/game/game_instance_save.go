package game

import "github.com/StellrisJAY/cloud-emu/emulator"

type Save struct {
	Data       []byte
	ExitSave   bool
	EmulatorId int64
	GameId     int64
}

type emulatorLoadSaveRequest struct {
	emulatorId   int64
	gameId       int64
	emulatorType string
	gameName     string
	gameData     []byte
	saveData     []byte
}

func (g *Instance) SaveGame() (*Save, error) {
	ch := make(chan ConsumerResult)
	g.messageChan <- &Message{Type: MsgSaveGame, resultChan: ch}
	result := <-ch
	close(ch)
	if result.Success {
		data := result.Data.([]byte)
		return &Save{
			Data:       data,
			EmulatorId: g.emulatorId,
			GameId:     g.gameId,
		}, nil
	} else {
		return nil, result.Error
	}
}

func (g *Instance) LoadSave(emulatorId, gameId int64, emulatorType, gameName string, gameData, saveData []byte) error {
	ch := make(chan ConsumerResult)
	g.messageChan <- &Message{Type: MsgLoadSave, resultChan: ch, Data: &emulatorLoadSaveRequest{
		emulatorId:   emulatorId,
		gameId:       gameId,
		emulatorType: emulatorType,
		gameName:     gameName,
		gameData:     gameData,
		saveData:     saveData,
	}}
	result := <-ch
	close(ch)
	if result.Success {
		return nil
	} else {
		return result.Error
	}
}

func (g *Instance) handleSaveGame() ConsumerResult {
	save, err := g.e.Save()
	if err != nil {
		return ConsumerResult{Error: err}
	}
	return ConsumerResult{Success: true, Data: save.SaveData()}
}

func (g *Instance) handleLoadSave(request *emulatorLoadSaveRequest) ConsumerResult {
	// 存档的模拟器或游戏与当前模拟器游戏不同，需要重启游戏实例
	if g.emulatorId != request.emulatorId || g.gameId != request.gameId {
		err := g.restartEmulator(request.gameName, request.gameData, request.emulatorType, request.emulatorId, request.gameId)
		if err != nil {
			return ConsumerResult{Error: err}
		}
	}
	// 切换存档
	err := g.e.LoadSave(&emulator.BaseEmulatorSave{
		Game: request.gameName,
		Data: request.saveData,
	})
	if err != nil {
		return ConsumerResult{Error: err}
	}
	return ConsumerResult{Success: true}
}
