package config

type Config struct {
	Game               string // Game 游戏文件路径
	EnableTrace        bool   // EnableTrace 是否在控制台打印trace
	Disassemble        bool   // Disassemble 打印程序的反汇编结果
	SnapshotSerializer string
	MuteApu            bool
	Debug              bool
}
