package game

type Save struct {
	Data       []byte
	ExitSave   bool
	EmulatorId int64
	GameId     int64
}

type emulatorRestartRequest struct {
	game         string
	gameData     []byte
	emulatorType string
	emulatorId   int64
	gameId       int64
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

func (g *Instance) handleSaveGame() ConsumerResult {
	save, err := g.e.Save()
	if err != nil {
		return ConsumerResult{Error: err}
	}
	return ConsumerResult{Success: true, Data: save.SaveData()}
}
