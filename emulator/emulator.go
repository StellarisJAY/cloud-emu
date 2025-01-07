package emulator

import (
	"context"
	"errors"
	"image"
)

const (
	TypeNESGO = "NESGO"
	TypeChip8 = "CHIP8"
	TypeDummy = "DUMMY"
	TypeGoboy = "GOBOY"
)

// IFrame 模拟器输出画面接口
type IFrame interface {
	// Width 画面宽度
	Width() int
	// Height 画面高度
	Height() int
	// YCbCr 画面YUV格式数据
	YCbCr() *image.YCbCr
	// Read 读取画面数据， func() 用于释放资源
	Read() (image.Image, func(), error)
}

type IGameFileRepo interface {
	// GetGameData 获取游戏文件数据
	GetGameData(ctx context.Context, game string) ([]byte, error)
}

type ControllerInfo struct {
	ControllerId int
	Label        string
}

// IEmulator 模拟器接口层，用于适配更多种类的模拟器
type IEmulator interface {
	// Start 启动模拟器
	Start() error
	// Pause 暂停模拟器
	Pause() error
	// Resume 继续运行模拟器
	Resume() error
	// Save 保存游戏
	Save() (IEmulatorSave, error)
	// LoadSave 加载存档，如果需要切换游戏，使用gameFileRepo加载游戏文件
	LoadSave(save IEmulatorSave) error
	// Restart 重启模拟器
	Restart(options IEmulatorOptions) error
	// Stop 结束模拟器
	Stop() error
	// SubmitInput 提交模拟器按键输入
	SubmitInput(controllerId int, keyCode string, pressed bool)
	// SetGraphicOptions 修改模拟器画面设置
	SetGraphicOptions(*GraphicOptions)

	GetCPUBoostRate() float64
	SetCPUBoostRate(float64) float64
	// OutputResolution 获取模拟器输出分辨率
	OutputResolution() (width, height int)

	MultiController() bool
	ControllerInfos() []ControllerInfo
}

// IEmulatorOptions 模拟器选项接口
type IEmulatorOptions interface {
	Game() string
	GameData() []byte
	AudioSampleRate() int
	AudioSampleChan() chan float32
	FrameConsumer() func(frame IFrame)
}

type BaseEmulatorOptions struct {
	game            string
	gameData        []byte
	audioSampleRate int
	audioSampleChan chan float32
	frameConsumer   func(frame IFrame)
}

type IEmulatorSave interface {
	GameName() string
	SaveData() []byte
}

type GraphicOptions struct {
	ReverseColor bool
	Grayscale    bool
}

type BaseEmulatorSave struct {
	Game string
	Data []byte
}

type BaseFrame struct {
	image  *image.YCbCr
	width  int
	height int
}

var (
	ErrorEmulatorNotSupported = errors.New("emulator not supported")
)

func MakeEmulator(emulatorType string, options IEmulatorOptions) (IEmulator, error) {
	switch emulatorType {
	case TypeNESGO:
		return makeNESEmulatorAdapter(options)
	case TypeChip8:
		e := makeChip8EmulatorAdapter(options)
		return e, nil
	case TypeDummy:
		return MakeDummyAdapter(options)
	case TypeGoboy:
		return newGoboyAdapter(options)
	default:
		return nil, ErrorEmulatorNotSupported
	}
}

func makeNESEmulatorAdapter(options IEmulatorOptions) (IEmulator, error) {
	e, err := makeNESEmulator(options)
	if err != nil {
		return nil, err
	}
	return &NesEmulatorAdapter{
		e:       e,
		options: options,
	}, nil
}

func (s *BaseEmulatorSave) GameName() string {
	return s.Game
}

func (s *BaseEmulatorSave) SaveData() []byte {
	return s.Data
}

func MakeBaseFrame(image *image.YCbCr, width, height int) IFrame {
	return &BaseFrame{
		image:  image,
		width:  width,
		height: height,
	}
}

func (b *BaseFrame) Width() int {
	return b.width
}

func (b *BaseFrame) Height() int {
	return b.height
}

func (b *BaseFrame) YCbCr() *image.YCbCr {
	return b.image
}

func (b *BaseFrame) Read() (image.Image, func(), error) {
	return b.image, func() {}, nil
}

func MakeBaseEmulatorOptions(game string, gameData []byte, audioSampleRate int, audioSampleChan chan float32,
	frameConsumer func(frame IFrame)) (IEmulatorOptions, error) {
	return &BaseEmulatorOptions{
		game:            game,
		gameData:        gameData,
		audioSampleRate: audioSampleRate,
		audioSampleChan: audioSampleChan,
		frameConsumer:   frameConsumer,
	}, nil
}

func (b *BaseEmulatorOptions) Game() string {
	return b.game
}

func (b *BaseEmulatorOptions) GameData() []byte {
	return b.gameData
}

func (b *BaseEmulatorOptions) AudioSampleRate() int {
	return b.audioSampleRate
}

func (b *BaseEmulatorOptions) AudioSampleChan() chan float32 {
	return b.audioSampleChan
}

func (b *BaseEmulatorOptions) FrameConsumer() func(frame IFrame) {
	return b.frameConsumer
}
