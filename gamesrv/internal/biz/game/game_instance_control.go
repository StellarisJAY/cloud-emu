package game

import "encoding/json"

type ControllerPlayer struct {
	ControllerId int
	Label        string
	UserId       int64
}

// handlePlayerControl 玩家控制消息处理
func (g *Instance) handlePlayerControl(keyCode string, action byte, userId int64) {
	for k, v := range g.controllerMap {
		if v == userId {
			g.e.SubmitInput(k, keyCode, action == MsgPlayerControlButtonPressed)
		}
	}
}

// handleChat 聊天消息处理, 广播给所有数据通道
func (g *Instance) handleChat(msg *Message) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	resp := &Message{Type: MsgChat, From: msg.From, To: 0, Data: msg.Data}
	for _, conn := range g.connections {
		raw, _ := json.Marshal(resp)
		_ = conn.dataChannel.Send(raw)
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
	g.controllerMap = make(map[int]int64)
	for _, cp := range cp {
		g.controllerMap[cp.ControllerId] = cp.UserId
	}
}

func (g *Instance) handleResetController(controllerId int) {

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
