//go:generate statik -src=./static
//go:generate go fmt statik/statik.go
package main

/*
#include <stdio.h>
#include <stdlib.h>
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -L${SRCDIR} -lwebInterface
#include "web_api.h"
*/
import "C"
import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/rakyll/statik/fs"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	_ "shixun600/statik"
	"shixun600/tpl"
	"strconv"
	"strings"
	"unsafe"
)

type Upload struct {
	Result int `json:"result"`
}

type PingStat struct {
	RecvPacket int `json:"recv_packet"`
	SendPacket int `json:"send_packet"`
	LossPacket int `json:"loss_packet"`
}

type NetworkConnInfo struct {
	SendVideoRate float32
	SendVideoLoss float32
	SendAudioRate float32
	SendAudioLoss float32
	RecVideoRate  float32
	RecVideoLoss  float32
	RecAudioRate  float32
	RecAudioLoss  float32
}

var (
	Username, Password = "admin", "passwd"
)

var (
	// 版本信息
	hwBuf    = C.calloc(1, (C.size_t)(20))
	swBuf    = C.calloc(1, (C.size_t)(20))
	panelBuf = C.calloc(1, (C.size_t)(20))
	v35Buf   = C.calloc(1, (C.size_t)(20))
	mcuBuf   = C.calloc(1, (C.size_t)(20))
	fpgaBuf  = C.calloc(1, (C.size_t)(20))

	ipaddrBuf  = C.calloc(1, (C.size_t)(20))
	netmaskBuf = C.calloc(1, (C.size_t)(20))
	gatewayBuf = C.calloc(1, (C.size_t)(20))

	ipRemoteBuf = C.calloc(1, (C.size_t)(20))
	ipPingBuf   = C.calloc(1, (C.size_t)(20))

	sendIpaddr    = C.calloc(1, (C.size_t)(20))
	recvIpaddr    = C.calloc(1, (C.size_t)(20))
	encsubcontent = C.calloc(1, (C.size_t)(20))

	onoffLoop  = C.calloc(1, (C.size_t)(20))
	onoffSound = C.calloc(1, (C.size_t)(20))
)

func checkSignin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("signed")
	if err != nil || cookie.Value != "true" {
		http.Redirect(w, r, "/signin", 302)
		return
	}
	log.Printf("cookie signed is %s.\n", cookie.Value)
}

func upgradeFirmwareError(w http.ResponseWriter) {
	log.Println("升级失败")
	tpl.Render(w, tpl.UPGRADE, &tpl.Msg{
		MemuUpgrade: true,
		MsgError:    "升级失败",
	})
	return
}

func jsonOk(w http.ResponseWriter) {
	js, _ := json.Marshal(tpl.Msg{
		Code:       0,
		MsgSuccess: "设置成功",
	})
	w.Write(js)
	return
}

func jsonError(w http.ResponseWriter) {
	js, _ := json.Marshal(tpl.Msg{
		Code:     1,
		MsgError: "设置失败",
	})
	w.Write(js)
	return
}

func main() {
	log.Println(" start    avc_web_database_init -----------------------------")
	// 板子初始化
	C.avc_web_database_init()
	log.Println(" end    avc_web_database_init -------------------------------")
	statikFS, err := fs.New()
	if err != nil {
		log.Println("Error:", err.Error())
		os.Exit(1)
	}
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	// 基础路由
	log.Println("start router --------------------------------------------")
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static/", http.FileServer(statikFS))
		fs.ServeHTTP(w, r)
	})
	r.Get("/signin", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("username") == Username && r.URL.Query().Get("password") == Password {
			c := http.Cookie{
				Name:  "signed",
				Value: "true",
			}
			http.SetCookie(w, &c)
			http.Redirect(w, r, "/", 302)
			return
		}
		tpl.Render(w, tpl.SIGNIN, &tpl.Msg{})
	})
	r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("username") == Username && r.FormValue("password") == Password {
			c := http.Cookie{
				Name:  "signed",
				Value: "true",
			}
			http.SetCookie(w, &c)
			tpl.Render(w, tpl.INDEX, &tpl.Msg{MsgSuccess: "登陆成功"})
			return
		}
		tpl.Render(w, tpl.SIGNIN, &tpl.Msg{MsgError: "账号或密码错误"})
	})
	r.Get("/signout", func(w http.ResponseWriter, r *http.Request) {
		c := http.Cookie{
			Name:  "signed",
			Value: "false",
		}
		http.SetCookie(w, &c)
		http.Redirect(w, r, "/signin", 302)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		checkSignin(w, r)
		tpl.Render(w, tpl.INDEX, &tpl.Msg{})
	})
	r.Get("/reboot", func(w http.ResponseWriter, r *http.Request) {
		checkSignin(w, r)
		tpl.Render(w, tpl.REBOOT, &tpl.Msg{MemuReboot: true})
	})
	r.Get("/reboot_start", func(w http.ResponseWriter, r *http.Request) {
		log.Println("调用重启接口")
		C.avc_web_board_reboot()
		w.Write([]byte(`{"Code": "0"}`))
	})
	r.Get("/reboot_end", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"Code": "0"}`))
	})

	r.Get("/upgrade", func(w http.ResponseWriter, r *http.Request) {
		checkSignin(w, r)
		tpl.Render(w, tpl.UPGRADE, &tpl.Msg{MemuUpgrade: true})
	})
	r.Post("/upgrade", func(w http.ResponseWriter, r *http.Request) {
		checkSignin(w, r)
		w.Header().Set("Content-Type", "text/html")
		r.ParseMultipartForm(32 << 20)

		fileFirmware, fh, fhErr := r.FormFile("file_firmware")
		v35Firmware, v35f, v35fErr := r.FormFile("v35_firmware")
		// muc 可以传两个文件            mcu  不需要md5比对
		mcuFirmware1, mcuf1, mcufErr1 := r.FormFile("mcu_firmware1")
		mcuFirmware2, mcuf2, mcufErr2 := r.FormFile("mcu_firmware2")
		fileMD5, md5file, md5Err := r.FormFile("file_md5")

		if fhErr != nil && v35fErr != nil && mcufErr1 != nil && mcufErr2 != nil && md5Err != nil {
			log.Println("读取文件失败,未上传任何文件", fhErr, md5Err, mcufErr2, mcufErr1, v35fErr)
			tpl.Render(w, tpl.UPGRADE, &tpl.Msg{
				MemuUpgrade: true,
				MsgError:    "读取文件失败,未上传任何文件",
			})
			return
		}
		var mcus bool = false
		// 先存两个不需要md校验的 mcu
		if mcufErr1 == nil && mcufErr2 == nil {
			// mcu文件名要完全匹配
			log.Println("mcu1 name-------", mcuf1.Filename)
			log.Println("muc2 name-------", mcuf2.Filename)

			var mcuNames = map[string]interface{}{"DZT_H2.57A.bin": "1", "DZT_H2.57B.bin": "1"}
			if mcuNames[mcuf1.Filename] != "1" || mcuNames[mcuf2.Filename] != "1" {
				tpl.Render(w, tpl.UPGRADE, &tpl.Msg{
					MemuUpgrade: true,
					MsgError:    "mcu固件文件名不匹配,请正确上传",
				})
				return
			}
			log.Println("save mcus----")
			mcus = true
			f, _ := os.OpenFile("/tmp/"+mcuf1.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			io.Copy(f, mcuFirmware1)
			f.Close()
			f2, _ := os.OpenFile("/tmp/"+mcuf2.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			io.Copy(f2, mcuFirmware2)
			f2.Close()
			// 调用mcu升级的api
			log.Println("开始升级mcu")
			C.avc_web_McuUpdate()
			log.Println("升级mcu成功")
		}
		// 有没有md5 没有直接返回错误
		if !mcus && md5Err != nil {
			log.Println("未升级mcu固件,必须上传md5")
			tpl.Render(w, tpl.UPGRADE, &tpl.Msg{
				MemuUpgrade: true,
				MsgError:    "未升级mcu固件,必须上传md5",
			})
			return
		}
		if v35fErr != nil && fhErr != nil && !mcus {
			log.Println("上传了md5，没有上传v35或者系统固件")
			tpl.Render(w, tpl.UPGRADE, &tpl.Msg{
				MemuUpgrade: true,
				MsgError:    "升级失败,请上传v35或者系统固件",
			})
			return
		}
		if md5Err == nil {
			log.Println("md5 name-------", md5file.Filename)
			var md5Name = map[string]interface{}{"md5checksum.txt": "1"}
			if md5Name[md5file.Filename] != "1" {
				tpl.Render(w, tpl.UPGRADE, &tpl.Msg{
					MemuUpgrade: true,
					MsgError:    "md5文件名不匹配,请正确上传",
				})
				return
			}
			log.Println("md5保存中。。。")
			f3, _ := os.OpenFile("/tmp/"+md5file.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			io.Copy(f3, fileMD5)
			f3.Close()
			log.Println("md5保存成功")

		}
		// fileMD5, md5file, md5Err := r.FormFile("file_md5")
		// v35和系统固件
		log.Println("v35fErr----", v35fErr)
		if v35fErr == nil {
			log.Println("v35 name-------", v35f.Filename)
			var v35Name = map[string]interface{}{"v35file": "1"}
			if v35Name[v35f.Filename] != "1" {
				tpl.Render(w, tpl.UPGRADE, &tpl.Msg{
					MemuUpgrade: true,
					MsgError:    "v35文件名不匹配,请正确上传",
				})
				return
			}

			log.Println("v35保存中。。。")
			f4, _ := os.OpenFile("/tmp/"+v35f.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			io.Copy(f4, v35Firmware)
			f4.Close()
			log.Println("v35保存成功")

			// 读取文件 和md5 进行校验
			// TODO
			log.Println("read ---- v35")
			v35fFirData, err := os.Open("/tmp/" + v35f.Filename)
			if err != nil {
				upgradeFirmwareError(w)
			}
			defer v35fFirData.Close()
			h := md5.New()
			log.Println("io.Copy(h, v35fFirData)-----------------")
			if _, err := io.Copy(h, v35fFirData); err != nil {
				upgradeFirmwareError(w)
			}
			log.Println("read md5---")
			md5Data, err := ioutil.ReadFile("/tmp/" + md5file.Filename)
			if err != nil {
				upgradeFirmwareError(w)
			}

			firMd5 := strings.TrimSpace(hex.EncodeToString(h.Sum([]byte(""))))
			md5Str := strings.Split(string(md5Data), " ")[0]
			log.Println("v35-------", firMd5)
			log.Println("md5-------", md5Str)
			log.Println("是否md5匹配-------", firMd5 == md5Str)
			if firMd5 == md5Str {
				log.Println("升级 v35 .........")
				C.avc_web_V35Update()
				log.Println("升级 v35 成功")

			}

		}
		// fileFirmware, fh, fhErr := r.FormFile("file_firmware")

		log.Println("系统固件err", fhErr)
		if fhErr == nil {
			log.Println("save 系统固件...")
			f5, _ := os.OpenFile("/tmp/"+fh.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			io.Copy(f5, fileFirmware)
			f5.Close()
			// 读取系统固件 和md5 进行校验
			log.Println("save 系统固件ok")
			log.Println("open 系统固件")
			// TODO
			fhFirData, err := os.Open("/tmp/" + fh.Filename)
			if err != nil {
				upgradeFirmwareError(w)
			}
			defer fhFirData.Close()
			h := md5.New()
			log.Println("io.Copy(h, fhFirData)---")
			if _, err := io.Copy(h, fhFirData); err != nil {
				upgradeFirmwareError(w)
			}
			md5Data, err := ioutil.ReadFile("/tmp/" + md5file.Filename)
			if err != nil {
				upgradeFirmwareError(w)
			}

			fir2Md5 := strings.TrimSpace(hex.EncodeToString(h.Sum([]byte(""))))
			md5Str := strings.Split(string(md5Data), " ")[0]

			log.Println("firMd5-------", fir2Md5)
			log.Println("md5-------", md5Str)
			log.Println("是否md5匹配-------", fir2Md5 == md5Str)

			if fir2Md5 == md5Str {
				log.Println("固件升级......")
				C.avc_web_firmwareUpdate()
				log.Println("固件升级成功")
			}
		}

		log.Println("升级成功")
		tpl.Render(w, tpl.UPGRADE, &tpl.Msg{
			MemuUpgrade: true,
			MsgSuccess:  "升级成功",
		})
		return
	})

	// v35
	r.Get("/v35", func(w http.ResponseWriter, r *http.Request) {
		checkSignin(w, r)
		log.Println("进入v35get----不查询数据")
		tpl.Render(w, tpl.V35, &tpl.Msg{
			V35Channel:      66,
			V35Version:      66,
			V35SendClock:    66,
			V35ReceiveClock: 66,
			MenuV35:         true,
		})
	})
	// v35查询比较慢  这里用ajax拿对应数据
	r.Get("/v35_info", func(w http.ResponseWriter, r *http.Request) {
		var v35Channel, v35Version, v35SendClock, v35ReceiveClock C.int
		log.Println("进入v35get----ajax查询")
		log.Println("进入v35get----ajax查询")
		log.Println("进入v35get----ajax查询")
		if _, err := C.avc_web_V35ChnGet(&v35Channel); err != nil {
			log.Println("get v35 channel  error:", err.Error())
		}
		if _, err := C.avc_web_V35VerGet(&v35Version); err != nil {
			log.Println("get v35 version  error:", err.Error())
		}
		if _, err := C.avc_web_v35TxClkGet(&v35SendClock); err != nil {
			log.Println("get v35 send clock  error:", err.Error())
		}
		if _, err := C.avc_web_v35RxClkGet(&v35ReceiveClock); err != nil {
			log.Println("get v35 rec clock  error:", err.Error())
		}
		log.Println("end v35----ajax查询")
		log.Println(v35Channel, "--", v35Version, "-", v35SendClock, "-", v35ReceiveClock)

		v35Data, _ := json.Marshal(tpl.Msg{
			V35Channel:      int(v35Channel),
			V35Version:      int(v35Version),
			V35SendClock:    int(v35SendClock),
			V35ReceiveClock: int(v35ReceiveClock),
		})
		w.Write(v35Data)
	})
	r.Post("/v35_channel_version", func(w http.ResponseWriter, r *http.Request) {
		checkSignin(w, r)
		// 这里数据类型可能存在问题
		v35ChannelInt, v35ChannelErr := strconv.Atoi(r.FormValue("v35_channel"))
		if v35ChannelErr != nil {
			log.Println("get post v35 channel error:", v35ChannelErr.Error())
		}
		v35VersionInt, v35VersionErr := strconv.Atoi(r.FormValue("v35_version"))
		if v35VersionErr != nil {
			log.Println("get post v35 version error:", v35VersionErr.Error())
		}

		if v35ChannelErr == nil {
			log.Println("start set v35 channel -- ")
			C.avc_web_V35ChnSet(C.int(v35ChannelInt))
			log.Println("end set v35 channel -- ")

		}
		if v35VersionErr == nil {
			log.Println("start set v35 ver -- ")
			C.avc_web_V35VerSet(C.int(v35VersionInt))
			log.Println("end set v35 channel -- ")
		}

		// 重新查询新的 send 和 rec clock
		var v35SendClock, v35ReceiveClock C.int
		if _, err := C.avc_web_v35TxClkGet(&v35SendClock); err != nil {
			log.Println("after update  get v35 send clock  error:", err.Error())
		}
		if _, err := C.avc_web_v35RxClkGet(&v35ReceiveClock); err != nil {
			log.Println("after update get v35 rec clock  error:", err.Error())
		}

		var newV35Clock = map[string]interface{}{"Code": "0", "sendClock": int(v35SendClock), "receiveClock": int(v35ReceiveClock)}
		res, _ := json.Marshal(newV35Clock)
		w.Write(res)
	})
	// config start
	r.Route("/config", func(r chi.Router) {
		r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
			checkSignin(w, r)
			C.avc_web_boardInfoGet(
				(*C.char)(unsafe.Pointer(swBuf)),
				(*C.char)(unsafe.Pointer(hwBuf)),
				(*C.char)(unsafe.Pointer(panelBuf)),
				(*C.char)(unsafe.Pointer(v35Buf)),
				(*C.char)(unsafe.Pointer(mcuBuf)),
				(*C.char)(unsafe.Pointer(fpgaBuf)))

			hwGo := C.GoString((*C.char)(unsafe.Pointer(hwBuf)))
			swGo := C.GoString((*C.char)(unsafe.Pointer(swBuf)))
			panelGo := C.GoString((*C.char)(unsafe.Pointer(panelBuf)))
			v35Go := C.GoString((*C.char)(unsafe.Pointer(v35Buf)))
			mcuGo := C.GoString((*C.char)(unsafe.Pointer(mcuBuf)))
			fpgaGo := C.GoString((*C.char)(unsafe.Pointer(fpgaBuf)))
			tpl.Render(w, tpl.CONFIG_VERSION, &tpl.Msg{
				SW:               swGo,
				HW:               hwGo,
				Panel:            panelGo,
				V35:              v35Go,
				Mcu:              mcuGo,
				Fpga:             fpgaGo,
				MenuEncodeDecode: true,
			})
		})
		r.Get("/video", func(w http.ResponseWriter, r *http.Request) {
			checkSignin(w, r)
			var inParam, outParam, std, encFrame C.int
			if _, err := C.avc_web_videoInParamGet(&inParam); err != nil {
				log.Println("get video in param error:", err.Error())
			}
			if _, err := C.avc_web_videoOutParamGet(&outParam); err != nil {
				log.Println("get video out param error:", err.Error())
			}
			if _, err := C.avc_web_videoStdGet(&std); err != nil {
				log.Println("get video std error:", err.Error())
			}
			if _, err := C.avc_web_video_enc_frameGet(&encFrame); err != nil {
				log.Println("get video enc frame error:", err.Error())
			}
			tpl.Render(w, tpl.CONFIG_VIDEO, &tpl.Msg{
				VideoInParam:  int(inParam),
				VideoOutParam: int(outParam),
				VideoStd:      int(std),
				VideoEncFrame: int(encFrame),
			})
		})
		r.Post("/video", func(w http.ResponseWriter, r *http.Request) {
			checkSignin(w, r)

			inParamInt, inParamErr := strconv.Atoi(r.FormValue("in_param"))
			if inParamErr != nil {
				log.Println("inParamErr:", inParamErr.Error())
			}
			outParamInt, outParamErr := strconv.Atoi(r.FormValue("out_param"))
			if outParamErr != nil {
				log.Println("outParamErr:", outParamErr.Error())
			}
			stdInt, stdErr := strconv.Atoi(r.FormValue("std"))
			if stdErr != nil {
				log.Println("stdErr:", stdErr.Error())
			}
			encFrameInt, encFrameErr := strconv.Atoi(r.FormValue("enc_frame"))
			if encFrameErr != nil {
				log.Println("encFrameErr:", encFrameErr.Error())
			}

			if inParamErr == nil && outParamErr == nil && stdErr == nil && encFrameErr == nil {
				if r.FormValue("in_param") != r.FormValue("in_param_old") {
					if _, err := C.avc_web_videoInParamSet(C.int(inParamInt)); err != nil {
						log.Println("set video in param  error:", err.Error())
					}
				}
				if r.FormValue("out_param") != r.FormValue("out_param_old") {
					if _, err := C.avc_web_videoOutParamSet(C.int(outParamInt)); err != nil {
						log.Println("set video out param  error:", err.Error())
					}
				}
				if r.FormValue("std") != r.FormValue("std_old") {
					if _, err := C.avc_web_videoStdSet(C.int(stdInt)); err != nil {
						log.Println("set video std  error:", err.Error())
					}
				}
				if r.FormValue("enc_frame") != r.FormValue("enc_frame_old") {
					if _, err := C.avc_web_video_enc_frameSet(C.int(encFrameInt)); err != nil {
						log.Println("set video enc frame  error:", err.Error())
					}
				}

				tpl.Render(w, tpl.CONFIG_VIDEO, &tpl.Msg{
					VideoInParam:     inParamInt,
					VideoOutParam:    outParamInt,
					VideoStd:         stdInt,
					VideoEncFrame:    encFrameInt,
					MenuEncodeDecode: true,
					MsgSuccess:       "设置成功"})
				return
			}
			tpl.Render(w, tpl.CONFIG_VIDEO, &tpl.Msg{
				VideoInParam:     inParamInt,
				VideoOutParam:    outParamInt,
				VideoStd:         stdInt,
				VideoEncFrame:    encFrameInt,
				MenuEncodeDecode: true,
				MsgError:         "设置失败"})
		})

		r.Get("/subtitle", func(w http.ResponseWriter, r *http.Request) {
			var encOnoff C.char
			var encPx, encPy C.uint
			var encSize, encColor, encPosition C.int

			if _, err := C.avc_web_captionCtrlGet(&encOnoff); err != nil {
				log.Println("get suntitle onoff  error:", err.Error())
			}
			if _, err := C.avc_web_captionNameGet((*C.char)(unsafe.Pointer(encsubcontent))); err != nil {
				log.Println("get subtitle content  error:", err.Error())
			}
			enccontent := C.GoString((*C.char)(unsafe.Pointer(encsubcontent)))
			if _, err := C.avc_web_captionCoordGet(&encPx, &encPy); err != nil {
				log.Println("get subtitle px py  error:", err.Error())
			}
			if _, err := C.avc_web_captionSizeGet(&encSize); err != nil {
				log.Println("get subtitle size  error:", err.Error())
			}
			if _, err := C.avc_web_captionColorGet(&encColor); err != nil {
				log.Println("get subtitle color  error:", err.Error())
			}
			if _, err := C.avc_web_captionPositGet(&encPosition); err != nil {
				log.Println("get subtitle position  error:", err.Error())
			}
			log.Println("on-off", encOnoff)
			log.Println("color--", encColor)
			log.Println("size--", encSize)
			log.Println("position--", encPosition)
			log.Println("px--", encPx)
			log.Println("py--", encPy)
			tpl.Render(w, tpl.CONFIG_SUBTITLE, &tpl.Msg{
				SubEncOnoff:      int(encOnoff),
				SubEncContent:    enccontent,
				SubEncPx:         int(encPx),
				SubEncPy:         int(encPy),
				SubEncPosition:   int(encPosition),
				SubEncSize:       int(encSize),
				SubEncColor:      int(encColor),
				MenuEncodeDecode: true,
			})
		})

		r.Post("/subtitle", func(w http.ResponseWriter, r *http.Request) {
			encOnoff, encOnoffErr := strconv.Atoi(r.FormValue("enc_onoff"))
			encPx, encPxErr := strconv.Atoi(r.FormValue("enc_px"))
			encPy, encPyErr := strconv.Atoi(r.FormValue("enc_py"))
			encSize, encSizeErr := strconv.Atoi(r.FormValue("enc_size"))
			encPosition, encPositionErr := strconv.Atoi(r.FormValue("enc_position"))
			encColor, encColorErr := strconv.Atoi(r.FormValue("enc_color"))

			encContent := C.CString(r.FormValue("enc_content"))
			defer C.free(unsafe.Pointer(encContent))

			log.Println(encOnoff, "-显示、隐藏-", r.FormValue("enc_onoff_old"), "--", encOnoffErr)
			log.Println(encPx, "-px-", r.FormValue("enc_px_old"), "--", encPxErr)
			log.Println(encPy, "-py-", r.FormValue("enc_py_old"), "--", encPyErr)
			log.Println(encSize, "-size-", r.FormValue("enc_size_old"), "--", encSizeErr)
			log.Println(encPosition, "-*****position******-", r.FormValue("enc_position_old"), "--", encPositionErr)
			log.Println(encColor, "-color-", r.FormValue("enc_color_old"), "--", encColorErr)
			log.Println(encContent, "-content-", r.FormValue("enc_content_old"))
			log.Println("-----------------------------------------------")
			log.Println("-----------------------------------------------")
			log.Println("-----------------------------------------------")
			if encOnoffErr == nil && encPxErr == nil && encPyErr == nil && encSizeErr == nil && encPositionErr == nil && encColorErr == nil {
				log.Println("enc_onoff--")
				if r.FormValue("enc_onoff") != r.FormValue("enc_onoff_old") {
					if _, err := C.avc_web_captionCtrlSet(C.char(encOnoff)); err != nil {
						log.Println("set subtitle onoff  error:", err.Error())
					}
				}
				log.Println("enc_content--")

				if r.FormValue("enc_content") != r.FormValue("enc_content_old") {
					if _, err := C.avc_web_captionNameSet(encContent); err != nil {
						log.Println("set subtitle content  error:", err.Error())
					}
				}
				log.Println("enc_px--")

				if (r.FormValue("enc_px") != r.FormValue("enc_px_old")) || r.FormValue("enc_py") != r.FormValue("enc_py_old") {
					if _, err := C.avc_web_captionCoordSet(C.uint(encPx), C.uint(encPy)); err != nil {
						log.Println("set subtitle px py  error:", err.Error())
					}
				}
				log.Println("--enc_color")

				if r.FormValue("enc_size") != r.FormValue("enc_size_old") {
					if _, err := C.avc_web_captionSizeSet(C.int(encSize)); err != nil {
						log.Println("set subtitle size  error:", err.Error())
					}
				}
				log.Println("enc_color--")

				if r.FormValue("enc_color") != r.FormValue("enc_color_old") {
					if _, err := C.avc_web_captionColorSet(C.int(encColor)); err != nil {
						log.Println("set subtitle color  error:", err.Error())
					}
				}
				log.Println("enc_position----------------------")

				if r.FormValue("enc_position") != r.FormValue("enc_position_old") {
					log.Println("enter enc_position")
					if _, err := C.avc_web_captionPositSet(C.int(encPosition)); err != nil {
						log.Println("error---------------")
						log.Println("set subtitle position  error:", err.Error())
					}
					log.Println("leave enc_position")

				}
				log.Println("update subtitle success ----------------------------------")
				log.Println("update subtitle success ----------------------------------")
				tpl.Render(w, tpl.CONFIG_SUBTITLE, &tpl.Msg{
					SubEncOnoff:      int(encOnoff),
					SubEncContent:    r.FormValue("enc_content"),
					SubEncPx:         int(encPx),
					SubEncPy:         int(encPy),
					SubEncSize:       int(encSize),
					SubEncColor:      int(encColor),
					SubEncPosition:   int(encPosition),
					MenuEncodeDecode: true,
					MsgSuccess:       "设置成功",
				})
				log.Println("before return subtitle  ----------------------------------")
				log.Println("before return subtitle  ----------------------------------")
				return
			}
			tpl.Render(w, tpl.CONFIG_SUBTITLE, &tpl.Msg{
				SubEncOnoff:      int(encOnoff),
				SubEncContent:    r.FormValue("enc_content"),
				SubEncPx:         int(encPx),
				SubEncPy:         int(encPy),
				SubEncSize:       int(encSize),
				SubEncColor:      int(encColor),
				SubEncPosition:   int(encPosition),
				MenuEncodeDecode: true,
				MsgError:         "设置失败",
			})
		})
		r.Get("/audio", func(w http.ResponseWriter, r *http.Request) {
			log.Println("enter audio----------------------------------------")
			log.Println("enter audio----------------------------------------")
			var soundGain, macGain, isMute, volume, encType, aacType, aacBitRate, aacSampleRate, aacEncapType C.int
			if _, err := C.avc_web_audioInGainGet(C.int(0), &soundGain); err != nil {
				log.Println("get audio sound gain  error:", err.Error())
			}
			if _, err := C.avc_web_audioInGainGet(C.int(1), &macGain); err != nil {
				log.Println("get audio mac gain  error:", err.Error())
			}
			if _, err := C.avc_web_audioOutMuteGet(&isMute); err != nil {
				log.Println("get audio out mute  error:", err.Error())
			}
			if _, err := C.avc_web_audioOutVolumeGet(&volume); err != nil {
				log.Println("get audio out volume  error:", err.Error())
			}
			if _, err := C.avc_web_audioEncodeTypeGet(&encType); err != nil {
				log.Println("get audio encode type  error:", err.Error())
			}
			if _, err := C.avc_web_audioAacGet(&aacType, &aacBitRate, &aacSampleRate, &aacEncapType); err != nil {
				log.Println("get audio aacs  error:", err.Error())
			}

			log.Println("after get audio info --------------------------")
			log.Println("get aacType -----------------------", aacType)
			log.Println("get aacType -----------------------", aacType)
			log.Println("get mute -----------------------", isMute)
			log.Println("get mute -----------------------", isMute)
			log.Println("render page------------------------------")

			tpl.Render(w, tpl.CONFIG_AUDIO, &tpl.Msg{
				AudioSoundGain:     int(soundGain),
				AudioMacGain:       int(macGain),
				AudioOutMute:       int(isMute),
				AudioOutVolume:     int(volume),
				AudioEncodeType:    int(encType),
				AudioAacType:       int(aacType),
				AudioAacBitRate:    int(aacBitRate),
				AudioAacSampleRate: int(aacSampleRate),
				AudioAacEncapType:  int(aacEncapType),
				MenuEncodeDecode:   true,
			})
		})

		// set audio/mac gain or volumn
		r.Get("/audio_mac_volunm_update", func(w http.ResponseWriter, r *http.Request) {
			soundGainInt, soundGainErr := strconv.Atoi(r.FormValue("sound_gain"))
			if soundGainErr != nil {
				log.Println("sound gain error:", soundGainErr.Error())
			} else {
				if r.FormValue("sound_gain") != r.FormValue("sound_gain_old") {
					if _, err := C.avc_web_audioInGainSet(C.int(0), C.int(soundGainInt)); err != nil {
						log.Println("set audio sound gain  error:", err.Error())
					}
					jsonOk(w)
					return
				}
			}
			macGainInt, macGainErr := strconv.Atoi(r.FormValue("mac_gain"))
			if macGainErr != nil {
				log.Println("mac gain error:", macGainErr.Error())
			} else {
				if r.FormValue("mac_gain") != r.FormValue("mac_gain_old") {
					if _, err := C.avc_web_audioInGainSet(C.int(1), C.int(macGainInt)); err != nil {
						log.Println("set audio mac gain  error:", err.Error())
					}
					jsonOk(w)
					return
				}
			}
			volumeInt, volumeErr := strconv.Atoi(r.FormValue("volume"))
			if volumeErr != nil {
				log.Println("volume error:", volumeErr.Error())
			} else {
				if r.FormValue("volume") != r.FormValue("volume_old") {
					if _, err := C.avc_web_audioOutVolumeSet(C.int(volumeInt)); err != nil {
						log.Println("set audio out volume  error:", err.Error())
					}
					jsonOk(w)
					return
				}
			}
			jsonError(w)
		})

		r.Post("/audio", func(w http.ResponseWriter, r *http.Request) {

			checkSignin(w, r)
			log.Println("enter post audio-----------------------")
			log.Println("enter post audio-----------------------")
			soundGainInt, soundGainErr := strconv.Atoi(r.FormValue("sound_gain"))
			if soundGainErr != nil {
				log.Println("sound gain error:", soundGainErr.Error())
			}
			macGainInt, macGainErr := strconv.Atoi(r.FormValue("mac_gain"))
			if macGainErr != nil {
				log.Println("mac gain error:", macGainErr.Error())
			}
			volumeInt, volumeErr := strconv.Atoi(r.FormValue("volume"))
			if volumeErr != nil {
				log.Println("volume error:", volumeErr.Error())
			}
			muteInt, muteErr := strconv.Atoi(r.FormValue("mute"))
			if muteErr != nil {
				log.Println("mute error:", muteErr.Error())
			}
			encType, encErr := strconv.Atoi(r.FormValue("enc_type"))
			if encErr != nil {
				log.Println("enc type error:", encErr.Error())
			}
			// aacTypeInt, aacTypeErr := strconv.Atoi(r.FormValue("aac_type"))
			// if aacTypeErr != nil {
			// 	log.Println("aac type error:", aacTypeErr.Error())
			// }
			aacBitRateInt, aacBitRateErr := strconv.Atoi(r.FormValue("aac_bit_rate"))
			if aacBitRateErr != nil {
				log.Println("aac bit rate error:", aacBitRateErr.Error())
			}
			aacSampleRateInt, aacSampleRateErr := strconv.Atoi(r.FormValue("aac_sample_rate"))
			if aacSampleRateErr != nil {
				log.Println("aac sample rate error:", aacSampleRateErr.Error())
			}
			aacEncapTypeInt, aacEncapTypeErr := strconv.Atoi(r.FormValue("aac_encap_type"))
			if aacEncapTypeErr != nil {
				log.Println("aac encap type error:", aacEncapTypeErr.Error())
			}

			log.Println("音频编码为-------------------", encType)
			log.Println("音频编码为-------------------", r.FormValue("enc_type"))
			if soundGainErr == nil && macGainErr == nil && muteErr == nil && volumeErr == nil && aacBitRateErr == nil && aacSampleRateErr == nil && aacEncapTypeErr == nil {

				if r.FormValue("mute") != r.FormValue("mute_old") {
					if _, err := C.avc_web_audioOutMuteSet(C.int(muteInt)); err != nil {
						log.Println("set audio out mute  error:", err.Error())
					}

				}
				if r.FormValue("enc_type") != r.FormValue("enc_type_old") {
					if _, err := C.avc_web_audioEncodeTypeSet(C.int(encType)); err != nil {
						log.Println("set audio encode type  error:", err.Error())
					}
				}

				if encType == 6 || encType == 7 {
					if r.FormValue("aac_sample_rate") != r.FormValue("aac_sample_rate_old") || r.FormValue("aac_bit_rate") != r.FormValue("aac_bit_rate_old") || r.FormValue("aac_encap_type") != r.FormValue("aac_encap_type_old") {
						if _, err := C.avc_web_audioAacSet(C.int(encType), C.int(aacBitRateInt), C.int(aacSampleRateInt), C.int(aacEncapTypeInt)); err != nil {
							log.Println("aac_type not change ,set aac params error:", err.Error())
						}
					}
				}
				// if r.FormValue("aac_type") != r.FormValue("aac_type_old") {
				// 	if _, err := C.avc_web_audioAacSet(C.int(aacTypeInt), C.int(aacBitRateInt), C.int(aacSampleRateInt), C.int(aacEncapTypeInt)); err != nil {
				// 		log.Println("set aac params error:", err.Error())
				// 	}
				// } else {

				// }
				log.Println("psot audio end----------------------")
				log.Println("render pages----------------------")
				tpl.Render(w, tpl.CONFIG_AUDIO, &tpl.Msg{
					AudioSoundGain:  soundGainInt,
					AudioMacGain:    macGainInt,
					AudioOutMute:    muteInt,
					AudioOutVolume:  volumeInt,
					AudioEncodeType: encType,
					// AudioAacType:       aacTypeInt,
					AudioAacBitRate:    aacBitRateInt,
					AudioAacSampleRate: aacSampleRateInt,
					AudioAacEncapType:  aacEncapTypeInt,
					MenuEncodeDecode:   true,
					MsgSuccess:         "设置成功"})
				return
			}
			tpl.Render(w, tpl.CONFIG_AUDIO, &tpl.Msg{
				AudioSoundGain:  soundGainInt,
				AudioMacGain:    macGainInt,
				AudioOutMute:    muteInt,
				AudioOutVolume:  volumeInt,
				AudioEncodeType: encType,
				// AudioAacType:       aacTypeInt,
				AudioAacBitRate:    aacBitRateInt,
				AudioAacSampleRate: aacSampleRateInt,
				AudioAacEncapType:  aacEncapTypeInt,
				MenuEncodeDecode:   true,
				MsgError:           "设置失败"})
		})
		r.Get("/network", func(w http.ResponseWriter, r *http.Request) {
			checkSignin(w, r)
			// var  ipaddr,netmask,gateway C.uchar
			var connetMode, videoEncodeRate C.int
			C.avc_web_boardNetGet(
				(*C.char)(unsafe.Pointer(ipaddrBuf)),
				(*C.char)(unsafe.Pointer(netmaskBuf)),
				(*C.char)(unsafe.Pointer(gatewayBuf)),
			)
			ipaddr := C.GoString((*C.char)(unsafe.Pointer(ipaddrBuf)))
			netmask := C.GoString((*C.char)(unsafe.Pointer(netmaskBuf)))
			gateway := C.GoString((*C.char)(unsafe.Pointer(gatewayBuf)))

			if _, err := C.avc_web_ConnetModeGet(&connetMode); err != nil {
				log.Println("get connect mode  error:", err.Error())
			}
			if _, err := C.avc_web_videoEncodeRateGet(&videoEncodeRate); err != nil {
				log.Println("get encode rate  error:", err.Error())
			}

			tpl.Render(w, tpl.CONFIG_NETWORK, &tpl.Msg{
				IP:               ipaddr,
				Netmask:          netmask,
				Gateway:          gateway,
				VideoEncodeRate:  int(videoEncodeRate),
				ConnetMode:       int(connetMode),
				MenuEncodeDecode: true,
			})
		})
		// network 打印
		r.Get("/network_data_info", func(w http.ResponseWriter, r *http.Request) {
			var sendVR, sendVl, sendAr, sendAl, recVR, recVl, recAr, recAl C.float
			C.avc_web_SendReceiveInfoGet(&sendVR, &sendVl, &recVR, &recVl, &sendAr, &sendAl, &recAr, &recAl)
			// C.avc_web_audio_sendReceiveInfoGet(&sendAr, &sendAl, &revAr, &revAl)
			// C.avc_web_video_sendReceiveInfoGet(&sendVR, &sendVl, &recVR, &recVl)
			data, _ := json.Marshal(NetworkConnInfo{
				float32(sendVR),
				float32(sendVl),
				float32(sendAr),
				float32(sendAl),
				float32(recVR),
				float32(recVl),
				float32(recAr),
				float32(recAl),
			})
			w.Write(data)

		})

		// 		SendVideoRate float32
		// SendVideoLoss float32
		// SendAudioRate float32
		// SendAudioLoss float32
		// RecVideoRate  float32
		// RecVideoLoss  float32
		// RecAudioRate  float32
		// RecAudioLoss  float32

		r.Post("/network_mode_rate", func(w http.ResponseWriter, r *http.Request) {
			log.Println("update network_mode or rate --------------------------------")
			log.Println("mode--------------------", r.FormValue("connect_mode"))
			log.Println("mode-old-------------------", r.FormValue("connect_mode_old"))
			log.Println("rate--------------------", r.FormValue("encode_rate"))
			log.Println("rate-old-------------------", r.FormValue("encode_rate_old"))

			connetModeInt, connetModeErr := strconv.Atoi(r.FormValue("connect_mode"))
			videoEncodeRateInt, videoEncodeRateErr := strconv.Atoi(r.FormValue("encode_rate"))
			if connetModeErr == nil && videoEncodeRateErr == nil {
				// 这里先不处理 connetModeInt
				// if videoEncodeRateErr == nil {
				if r.FormValue("connect_mode") != r.FormValue("connect_mode_old") {
					log.Println("mode---not equal---avc_web_ConnetModeSet")
					if _, err := C.avc_web_ConnetModeSet(C.int(connetModeInt)); err != nil {
						log.Println("set connect mode error:", err.Error())
					}
					log.Println("mode set ok")
				}
				if r.FormValue("encode_rate") != r.FormValue("encode_rate_old") {
					log.Println("rate---not equal---avc_web_videoEncodeRateSet")
					if _, err := C.avc_web_videoEncodeRateSet(C.int(videoEncodeRateInt)); err != nil {
						log.Println("set encode rate error:", err.Error())
					}
					log.Println("rate set ok")

				}
				jsonOk(w)
			}
			jsonError(w)

		})

		r.Post("/network_ip_netmask_gateway", func(w http.ResponseWriter, r *http.Request) {
			if len(r.FormValue("ipaddr")) != 0 && len(r.FormValue("netmask")) != 0 && len(r.FormValue("gateway")) != 0 {
				ipaddr := C.CString(r.FormValue("ipaddr"))
				netmask := C.CString(r.FormValue("netmask"))
				gateway := C.CString(r.FormValue("gateway"))
				defer C.free(unsafe.Pointer(ipaddr))
				defer C.free(unsafe.Pointer(netmask))
				defer C.free(unsafe.Pointer(gateway))
				if r.FormValue("ipaddr") != r.FormValue("ipaddr_old") || r.FormValue("netmask") != r.FormValue("netmask_old") || r.FormValue("gateway") != r.FormValue("gateway_old") {
					if _, err := C.avc_web_boardNetSet(ipaddr, netmask, gateway); err != nil {
						log.Println("set ipaddr or netmask or gateway error:", err.Error())
					}
				}
				jsonOk(w)
			}
			jsonError(w)
		})
		r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
			checkSignin(w, r)
			// 初始化 环路测试 声音测试
			if _, err := C.avc_web_loopGet((*C.char)(unsafe.Pointer(onoffLoop))); err != nil {
				log.Println("get test loop  error:", err.Error())
			}
			if _, err := C.avc_web_audioLoopGet((*C.char)(unsafe.Pointer(onoffSound))); err != nil {
				log.Println("get test audio  error:", err.Error())
			}
			onoffLoopStr2 := C.GoString((*C.char)(unsafe.Pointer(onoffLoop)))
			onoffSoundStr2 := C.GoString((*C.char)(unsafe.Pointer(onoffSound)))

			var onoffLoopStr, onoffSoundStr C.char
			if _, err := C.avc_web_loopGet(&onoffLoopStr); err != nil {
				log.Println("get test loop  error:", err.Error())
			}
			if _, err := C.avc_web_audioLoopGet(&onoffSoundStr); err != nil {
				log.Println("get test audio  error:", err.Error())
			}
			log.Println("loop-----------------===", string(onoffLoopStr))
			log.Println("loop-----------------===", onoffLoopStr)
			log.Println("loop-----------------===", int(onoffLoopStr))
			log.Println("sound-----------------===", string(onoffSoundStr))
			log.Println("sound-----------------===", onoffSoundStr)
			log.Println("sound-----------------===", int(onoffSoundStr))

			log.Println("v2----")
			log.Println("loop--------", onoffLoopStr2)
			log.Println("sound--------", onoffSoundStr2)
			tpl.Render(w, tpl.CONFIG_CHECK, &tpl.Msg{
				OnoffLoop:        int(onoffLoopStr),
				OnoffSound:       int(onoffSoundStr),
				OnoffIp:          0,
				MenuEncodeDecode: true,
			})

		})
		// 开启或停止 loop 测试 0停止 1 3
		r.Get("/check_loop", func(w http.ResponseWriter, r *http.Request) {
			var onoffLoop C.int
			onoffLoopInt, onoffErr := strconv.Atoi(r.FormValue("onoff_loop"))
			if onoffErr == nil {
				onoffLoop = C.int(onoffLoopInt)
				if _, err := C.avc_web_loopSet(onoffLoop); err != nil {
					log.Println("test loop  error:", err.Error())
				}
				w.Write([]byte(`{"Code": "0", "Msg": "设置成功"}`))
				return
			}
			w.Write([]byte(`{"Code": "1", "Msg": "设置失败"}`))
		})
		// 开启或停止audio测试 0停止 1
		r.Get("/check_audio", func(w http.ResponseWriter, r *http.Request) {
			var onoffSoundAudio C.int
			onoffSoundInt, onoffErr := strconv.Atoi(r.FormValue("onoff_sound"))
			if onoffErr == nil {
				onoffSoundAudio = C.int(onoffSoundInt)
				if _, err := C.avc_web_audioLoopSet(onoffSoundAudio); err != nil {
					log.Println("test audio loop  error:", err.Error())
				}
				w.Write([]byte(`{"Code": "0", "Msg": "设置成功"}`))
				return
			}
			w.Write([]byte(`{"Code": "1", "Msg": "设置失败"}`))
		})
		// 开启或停止ping测试
		r.Get("/check_ping", func(w http.ResponseWriter, r *http.Request) {
			var onoff C.int
			onoffInt, onoffErr := strconv.Atoi(r.FormValue("ip_off"))
			ipRemote := C.CString(r.FormValue("ip"))
			defer C.free(unsafe.Pointer(ipRemote))
			if onoffErr == nil {
				onoff = C.int(onoffInt)
				C.avc_web_pingTestSet(onoff, ipRemote)
				w.Write([]byte(`{"Code": "0", "Msg": "设置成功"}`))
				return
			}
			w.Write([]byte(`{"Code": "1", "Msg": "设置失败"}`))
		})
		// 打印ping测试结果
		r.Get("/check_ping_stat", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			var RecvPacket, SendPacket, LossPacket C.int
			C.avc_web_pingTestPrint(&SendPacket, &RecvPacket, &LossPacket)
			v := PingStat{int(RecvPacket), int(SendPacket), int(LossPacket)}
			js, _ := json.Marshal(v)
			w.Write(js)
		})

	})

	log.Println("start http service  on 80 ------------------------------------------------")
	if err := http.ListenAndServe("0.0.0.0:80", r); err != nil {
		log.Println("start http service error:", err.Error())
	}
}
