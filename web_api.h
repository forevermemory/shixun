#include <stdio.h>


void avc_web_database_init();

//音频调音增益、麦克增益
void avc_web_audioInGainGet(int source,int *gain);
void avc_web_audioInGainSet(int source,int gain);

//静音输出
void avc_web_audioOutMuteGet(int *mute);
void avc_web_audioOutMuteSet(int mute);

//音量控制
void avc_web_audioOutVolumeGet(int *volume);
void avc_web_audioOutVolumeSet(int volume);

//音频编码
void avc_web_audioEncodeTypeGet(int *type);
void avc_web_audioEncodeTypeSet(int type);

//aac 码率、采样率、封装设置
void avc_web_audioAacGet(int *type, int *bitRate , int *sampleRate, int *encapType);
void avc_web_audioAacSet(int type, int bitRate , int sampleRate, int encapType);




//输入源
void avc_web_videoInParamGet(int *def_res);
void avc_web_videoInParamSet(int def_res);

//输出源
void avc_web_videoOutParamGet(int *def_res);
void avc_web_videoOutParamSet(int def_res);

//输出制式
void avc_web_videoStdGet(int *std);
void avc_web_videoStdSet(int std);

//视频编码
void avc_web_video_enc_frameGet(int *frame);
void avc_web_video_enc_frameSet(int frame);





//字幕显示隐藏
void avc_web_captionCtrlGet(char *onoff);
void avc_web_captionCtrlSet(char onoff);

//字幕内容
void avc_web_captionNameGet(char *name);
void avc_web_captionNameSet(char *name);

//字幕坐标
void avc_web_captionCoordGet(unsigned int *px, unsigned int *py);
void avc_web_captionCoordSet(unsigned int px, unsigned int py);

//字幕位置
void avc_web_captionPositGet(int *position);
void avc_web_captionPositSet(int position);

//字幕字体大小
void avc_web_captionSizeSet(int size);
void avc_web_captionSizeGet(int *size);

//字幕字体颜色
void avc_web_captionColorGet(int *color);
void avc_web_captionColorSet(int color);






//V35通道
void avc_web_V35ChnGet(int *v35chn);
void avc_web_V35ChnSet(int v35chn);
// v35版本
void avc_web_V35VerGet(int *v35ver);
void avc_web_V35VerSet(int v35ver);
//V35发送时钟
void avc_web_v35TxClkGet(int *txclk);
//v35接收时钟
void avc_web_v35RxClkGet(int *rxclk);


//版本信息 *
void avc_web_boardInfoGet(char *sw, char *hw, char *panel, char *v35, char *mcu, char *fpga);



//环回测试  0停止 1远端环回 3本地还回
void avc_web_loopGet(char *loop);
void avc_web_loopSet(int loop);

//声音测试 0停止 1测试
void avc_web_audioLoopGet(char *onoff);
void avc_web_audioLoopSet(int onoff);

//ping测试  0停止 1测试
void avc_web_pingTestSet(int onoff,char *ip);
void avc_web_pingTestPrint(int *sendPacket,int *recvPacket,int *lossPacket);




//连接模式
void avc_web_ConnetModeGet(int *mode);
void avc_web_ConnetModeSet(int mode);

//码率
void avc_web_videoEncodeRateGet(int *rate);
void avc_web_videoEncodeRateSet(int rate);

//本机IP地址，子网掩码，网关
void avc_web_boardNetGet(char *ipaddr,char *mask,char *gateway);
void avc_web_boardNetSet(char *ipaddr,char *mask,char *gateway);

//发送包、接收包、丢包率
void avc_web_audio_sendReceiveInfoGet(float *s_rate, float *s_loss, float *r_rate, float *r_loss);
void avc_web_video_sendReceiveInfoGet(float *s_rate, float *s_loss, float *r_rate, float *r_loss);




//固件升级  需要md5
void avc_web_firmwareUpdate();

//v35 需要md5
void avc_web_V35Update();

//mcu  不需要md5比对
void avc_web_McuUpdate();
//md5文件v35Buf
int avc_checksum_md5(char *md5_linux, char *md5_win);




//重启
void avc_web_board_reboot(void);









