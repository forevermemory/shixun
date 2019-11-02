package tpl

import (
	"html/template"
	"net/http"
)

type Msg struct {
	Code int

	// 编码解码器
	MenuEncodeDecode bool
	MenuV35          bool
	MemuUpgrade      bool
	MemuReboot       bool
	// 消息提示 ok
	MsgError   string
	MsgSuccess string

	// 版本信息 ok
	HW    string
	SW    string
	Panel string
	V35   string
	Mcu   string
	Fpga  string

	// 音频设置 ok
	// 调音增益   静音输出 输出音量控制	  音频编码
	AudioSoundGain  int
	AudioMacGain    int
	AudioOutMute    int
	AudioOutVolume  int
	AudioEncodeType int
	// aac 码率、采样率、封装设置    aac-lc   he-aac
	AudioAacType       int
	AudioAacBitRate    int
	AudioAacSampleRate int
	AudioAacEncapType  int

	// 视频设置 输入源 输出源 输出制式 视频编码 ok
	VideoInParam  int
	VideoOutParam int
	VideoStd      int
	VideoEncFrame int

	// 字幕设置 ok
	SubEncOnoff    int
	SubEncContent  string
	SubEncPx       int
	SubEncPy       int
	SubEncPosition int
	SubEncSize     int
	SubEncColor    int

	// 网络设置 ok
	ConnetMode      int
	VideoEncodeRate int
	IP              string
	Netmask         string
	Gateway         string

	// 系统诊断
	// 两秒后 打印
	RecvPacket int
	SendPacket int
	LossPacket int
	// 环回自测  声音自测
	OnoffLoop  int
	OnoffSound int
	OnoffIp    int

	// v35
	V35Channel      int
	V35Version      int
	V35SendClock    int
	V35ReceiveClock int

	// 串口设置  波特率  设备编号
	UartBuadRate int
	UartDeviceIp int
}

func Render(w http.ResponseWriter, htm string, msg *Msg) {
	t := template.Must(template.New("layout").Parse(LAYOUT))
	t = template.Must(t.Parse(htm))
	w.Header().Set("P3P", "CP='CAO PSA OUR'")
	t.ExecuteTemplate(w, "layout", msg)
}
