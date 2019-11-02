package tpl

const CONFIG_AUDIO = `
{{define "content"}}
 <div id="version">
    <div>
        <h3>设备状态</h3>
        <ul >
            <li>
                <a href="/config/network">网络</a>
            </li>
            <li  class="active">
                <a href="/config/audio">音频</a>
            </li>
            <li>
                <a href="/config/video">视频</a>
            </li>
            <li>
                <a href="/config/check">诊断</a>
            </li>
            <li >
                <a href="/config/version">版本</a>
            </li>
            <li>
                <a href="/config/subtitle">字幕</a>
            </li>
        </ul>
            <div class="well video">
            <form action="/config/audio" method="post">

                <div class="input white" style="height: 30px;">
                    <span>调音增益:</span>
                    <input type="text" name="sound_gain" id="" value="{{.AudioSoundGain}}"  style="width: 25%"> 
                    <input type="text" name="sound_gain_old" id="" value="{{.AudioSoundGain}}"  style="display:none;"> 
                    <input type="button" onClick="audioSoundGain(this)" class="btn" value="修改"> db(0~16)
                </div>

                <div class="input white" style="height: 30px;">
                    <span>麦克增益:</span>
                    <input type="text" name="mac_gain" id="" value="{{.AudioMacGain}}"  style="width: 25%"> 
                    <input type="text" name="mac_gain_old" id="" value="{{.AudioMacGain}}"  style="display:none;"> 
                    <input type="button" onClick="audioMacGain(this)" class="btn" value="修改"> db(0~16)
                </div> 

                <div class="input white" style="height: 30px;">
                    <span>音量控制:</span>
                    <input type="text" name="volume" id="" value="{{.AudioOutVolume}}"  style="width: 25%"> 
                    <input type="text" name="volume_old" id="" value="{{.AudioOutVolume}}"   style="display:none;"> 
                    <input type="button" onClick="audioVolumnUpdate(this)" class="btn" value="修改"> db(0~32)
                </div> 

                <div class="input" style="height: 30px;">
                    <span>静音输出:</span> 
                    <input type="radio" name="mute" id="" value="1"> <label class="white">是</label>
                    <input type="radio" name="mute" id="" value="0"> <label class="white">否</label>
                    <input type="text" name="mute_old" id="" value="{{.AudioOutMute}}"  style="display:none;"> 

                </div> 

                <div class="input" style="height: 30px;">
                    <span>音频编码:</span>
                    <input type="radio" name="enc_type" id="" value="6"> <label class="white">AAC-LC</label>
                    <input type="radio" name="enc_type" id="" value="7"> <label class="white">HE-AAC</label>
                    <input type="radio" name="enc_type" id="" value="5"> <label class="white">G.726</label>
                    <input type="radio" name="enc_type" id="" value="3"> <label class="white">G.729A</label>
                    <input type="radio" name="enc_type" id="" value="0"> <label class="white">G.711A</label>
                    <input type="radio" name="enc_type" id="" value="4"> <label class="white">ADPCM</label>
                    <input type="text" name="enc_type_old" id="" value="{{.AudioEncodeType}}"  style="display:none;"> 

                </div> 

                <div class="input acc hidden" style="height: 30px;" id="AACLC2HEAAC">
                    <span style="font-size: 14px;" class="aac-type-span"></span>
                       <select id="" name="aac_sample_rate" style="width: 8%">
                            <option value="8000">8K</option>
                            <option value="16000">16K</option>
                            <option value="32000">32K</option>
                            <option value="48000">48K</option>
                            <option value="44100">44.1K</option>
                        </select>
                        <input type="text" name="aac_sample_rate_old" id="" value="{{.AudioAacSampleRate}}"  style="display:none;"> 

                    <span class="aac-type-span"></span>
                       <select id="" name="aac_bit_rate" style="width: 8%">
                            <option value="8000">8K</option>
                            <option value="16000">16K</option>
                            <option value="32000">32K</option>
                            <option value="64000">64K</option>
                            <option value="128000">128K</option>
                            <option value="192000">192K</option>
                            <option value="288000">288K</option>
                            <option value="320000">320K</option>
                        </select>
                        <input type="text" name="aac_bit_rate_old" id="" value="{{.AudioAacBitRate}}"  style="display:none;"> 

                    <span class="aac-type-span"></span>
                       <select id="" name="aac_encap_type" style="width: 8%">
                            <option value="0">ADTS</option>
                            <option value="1">LATM</option>
                        </select>
                        <input type="text" name="aac_encap_type_old" id="" value="{{.AudioAacEncapType}}"  style="display:none;"> 

                </div> 

                <input type="submit" onClick="removeDisabled()" class="btn ml35" value="修改">

            </form>

            </div>
    </div>
</div>
<script>
$('select[name="aac_sample_rate"]').val({{.AudioAacSampleRate}})
$('select[name="aac_bit_rate"]').val({{.AudioAacBitRate}})
$('select[name="aac_encap_type"]').val({{.AudioAacEncapType}})



// 初始化 静音输出 音频编码 aac类别 的选择情况
$.each($('input[type="radio"][name="mute"]'), function (i, val) { 
    if($(val).val() == '{{.AudioOutMute}}'){
        $(val).attr('checked','checked')
    }
});
$.each($('input[type="radio"][name="enc_type"]'), function (i, val) { 
    if($(val).val() == '{{.AudioEncodeType}}'){
        $(val).attr('checked','checked')
    }
});
// $.each($('input[type="radio"][name="aac_type"]'), function (i, val) { 
//     if($(val).val() == '{{.AudioAacType}}'){
//         $(val).attr('checked','checked')
//     }
// });

setAacType()


function audioSoundGain(event) {  
    audioMacVolunmAjax(event,'sound_gain','sound_gain_old',16)
}
function audioMacGain(event) {  
    audioMacVolunmAjax(event,'mac_gain','mac_gain_old',16)
}
function audioVolumnUpdate(event) {  
    audioMacVolunmAjax(event,'volume','volume_old',32)
}

function beforeCommit(type,size) {  
 
}

// 调音增益 麦克增益 音量控制 的ajax
function audioMacVolunmAjax(event,lated,old,size) {  
    emptyAlertMsg()
    let value = $('input[name="'+lated+'"]').val()
    if(!/^[0-9]{1,2}$/.test(value)){
        $('#errorMSG').text('请正确输入参数')
        return false
    }
    if(parseInt(value) >size || parseInt(value)<0){
        $('#errorMSG').text('参数不在指定范围内')
        return false
    }


    let newValue = $('input[name="'+lated+'"]').val()
    let oldValue = $('input[name="'+old+'"]').val()
    if (newValue == oldValue){
        $('#errorMSG').text('请先更改内容后再修改')
        return false
    }
    // 添加
    if($(event).hasClass('disabled')){
        return false
    }
    $(event).addClass('disabled')

    let data = {}
    data[old] = oldValue
    data[lated] = newValue
    $.ajax({
        type: "get",
        url: "/config/audio_mac_volunm_update",
        data: data,
        success: function (response) {
            let res = JSON.parse(response)
            $(event).removeClass('disabled')
            if(res['Code'] == '0'){
                $('#successMSG').text(res['MsgSuccess'])
                $('input[name="'+old+'"]').val($('input[name="'+lated+'"]').val())
            }else{
                $('#errorMSG').text(res['MsgError'])
            }
        },error: function () {
            $(event).removeClass('disabled')
            $('#errorMSG').text('服务端错误')
        }
    });
}




// 当音频编码选择变化时候触发
$('input[name="enc_type"]').change(function (event) { 
    setAacType()
});

// 在页面初始化时候 判断下面三个的内容
function setAacType() {  
    let aacType = $('input[name="enc_type"]:checked').val()

 
    var aacLcArray = ['AAC-LC采样率:','AAC-LC码率:','AAC-LC封装:']
    var heAacArray = ['HE-AAC采样率:','HE-AAC码率:','HE-AAC封装:']
    console.log(aacType)

    if(aacType == '6'){
        // AAC-LC
        $('#AACLC2HEAAC').removeClass('hidden')
        $.each($('.aac-type-span'), function (i, value) { 
            $(value).text(aacLcArray[i])
        });
    }else if (aacType == '7'){
        $('#AACLC2HEAAC').removeClass('hidden')
        $.each($('.aac-type-span'), function (i, value) { 
            $(value).text(heAacArray[i])
        });
    }else{
        // 隐藏下面三个
        $('#AACLC2HEAAC').addClass('hidden')
    }
}

</script>
{{end}}`
