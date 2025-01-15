package emulator

import (
	"errors"
	"image"
	"image/color"
	"time"
)

const (
	CodeNesgo       = "NESGO"
	CodeChip8       = "CHIP8"
	CodeDummy       = "DUMMY"
	CodeGoboy       = "GOBOY"
	CodeFoglemanNES = "FOGLEMAN/NES"
	CodeMagia       = "MAGIA"
	CodeDawnGB      = "DAWNGB"
)

const FrameInterval = time.Second / 60

var (
	supportedEmulators = make(map[string]Info)
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

type Info struct {
	EmulatorType           string
	Provider               string
	Description            string
	Name                   string
	SupportSave            bool
	SupportGraphicSettings bool
	EmulatorCode           string
}

var (
	ErrorEmulatorNotSupported = errors.New("emulator not supported")
)

func MakeEmulator(emulatorCode string, options IEmulatorOptions) (IEmulator, error) {
	switch emulatorCode {
	case CodeNesgo:
		return makeNESEmulatorAdapter(options)
	case CodeDummy:
		return MakeDummyAdapter(options)
	case CodeGoboy:
		return newGoboyAdapter(options)
	case CodeFoglemanNES:
		return newFoglemanNesAdapter(options)
	case CodeMagia:
		return newMagiaAdapter(options)
	case CodeDawnGB:
		return newDawnGbAdapter(options)
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

func MakeEmptyBaseFrame(rect image.Rectangle) *BaseFrame {
	return &BaseFrame{
		image:  image.NewYCbCr(rect, image.YCbCrSubsampleRatio420),
		width:  rect.Dx(),
		height: rect.Dy(),
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

func (b *BaseFrame) FromRGBA(rgba *image.RGBA) {
	for x := 0; x < rgba.Rect.Dx(); x++ {
		for y := 0; y < rgba.Rect.Dy(); y++ {
			r0, g0, b0, _ := rgba.At(x, y).RGBA()
			Y, Cb, Cr := color.RGBToYCbCr(uint8(r0), uint8(g0), uint8(b0))
			yOff := b.image.YOffset(x, y)
			cOff := b.image.COffset(x, y)
			b.image.Y[yOff] = Y
			b.image.Cb[cOff] = Cb
			b.image.Cr[cOff] = Cr
		}
	}
}

func (b *BaseFrame) FromRGBARaw(data []byte) {
	for i := 0; i < len(data); i += 4 {
		r0, g0, b0 := data[i], data[i+1], data[i+2]
		Y, Cb, Cr := color.RGBToYCbCr(r0, g0, b0)
		x := (i / 4) % b.width
		y := (i / 4) / b.width
		yOff := b.image.YOffset(x, y)
		cOff := b.image.COffset(x, y)
		b.image.Y[yOff] = Y
		b.image.Cb[cOff] = Cb
		b.image.Cr[cOff] = Cr
	}
}

func (b *BaseFrame) FromNRGBAColors(rgbaColors []color.NRGBA) {
	for i, rgba := range rgbaColors {
		Y, Cb, Cr := color.RGBToYCbCr(rgba.R, rgba.G, rgba.B)
		x := i % b.width
		y := i / b.width
		yOff := b.image.YOffset(x, y)
		cOff := b.image.COffset(x, y)
		b.image.Y[yOff] = Y
		b.image.Cb[cOff] = Cb
		b.image.Cr[cOff] = Cr
	}
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

func GetSupportedEmulators() []Info {
	var infos []Info
	for _, v := range supportedEmulators {
		infos = append(infos, v)
	}
	return infos
}

func GetSupportedEmulatorTypes() []string {
	types := make(map[string]struct{})
	for _, v := range supportedEmulators {
		types[v.EmulatorType] = struct{}{}
	}
	var result []string
	for k := range types {
		result = append(result, k)
	}
	return result
}
