package game

import (
	"time"
)

type ControllerPlayer struct {
	ControllerId int
	Label        string
	UserId       int64
}

type applyMacroRequest struct {
	UserId int64
	Keys   []string
}

// handlePlayerControl 玩家控制消息处理
func (g *Instance) handlePlayerControl(keyCode string, action byte, userId int64) {
	for k, v := range g.controllerMap {
		if v == userId {
			g.e.SubmitInput(k, keyCode, action == MsgPlayerControlButtonPressed)
		}
	}
}

func (g *Instance) SetController(cp []ControllerPlayer) []ControllerPlayer {
	resultChan := make(chan ConsumerResult)
	g.messageChan <- &Message{
		Type:       MsgSetController,
		Data:       cp,
		resultChan: resultChan,
	}
	result := <-resultChan
	return result.Data.([]ControllerPlayer)
}

func (g *Instance) handleSetController(cp []ControllerPlayer) {
	clear(g.controllerMap)
	for _, cp := range cp {
		g.controllerMap[cp.ControllerId] = cp.UserId
	}
}

func (g *Instance) GetControllerPlayers() []ControllerPlayer {
	resultChan := make(chan ConsumerResult)
	g.messageChan <- &Message{
		Type:       MsgGetController,
		resultChan: resultChan,
	}
	result := <-resultChan
	if result.Success {
		return result.Data.([]ControllerPlayer)
	}
	return nil
}

func (g *Instance) getControllerPlayer() []ControllerPlayer {
	controllers := g.e.ControllerInfos()
	players := make([]ControllerPlayer, len(controllers))
	for i, c := range controllers {
		players[i] = ControllerPlayer{
			ControllerId: c.ControllerId,
			Label:        c.Label,
			UserId:       g.controllerMap[c.ControllerId],
		}
	}
	return players
}

func (g *Instance) ApplyMacro(userId int64, keys []string) {
	resultChan := make(chan ConsumerResult)
	g.messageChan <- &Message{
		Type:       MsgApplyMacro,
		Data:       applyMacroRequest{UserId: userId, Keys: keys},
		resultChan: resultChan,
	}
	<-resultChan
}

func (g *Instance) handleApplyMacro(request applyMacroRequest) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	controller := -1
	for k, v := range g.controllerMap {
		if v == request.UserId {
			controller = k
			break
		}
	}
	if controller == -1 {
		return
	}
	for _, key := range request.Keys {
		g.e.SubmitInput(controller, key, true)
	}
	// TODO 模拟按键按下后，等待一段时间，再释放按键
	timer := time.NewTimer(time.Millisecond * 300)
	go func() {
		defer timer.Stop()
		<-timer.C
		for _, key := range request.Keys {
			g.e.SubmitInput(controller, key, false)
		}
	}()
}

func (g *Instance) resetController() {
	clear(g.controllerMap)
}
