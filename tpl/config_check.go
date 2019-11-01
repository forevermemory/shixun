package tpl

const CONFIG_CHECK = `
{{define "content"}}
<div id="version">
    <div>
        <h3>设备状态</h3>
        <ul >
            <li>
                <a href="/config/network">网络</a>
            </li>
            <li  >
                <a href="/config/audio">音频</a>
            </li>
            <li >
                <a href="/config/video">视频</a>
            </li>
            <li class="active">
                <a href="/config/check">诊断</a>
            </li>
            <li >
                <a href="/config/version">版本</a>
            </li>
            <li>
                <a href="/config/subtitle">字幕</a>
            </li>
        </ul>
        <form action="/config/test" method="post">
            <div class="well video">
                <div class="input">
                    <span>还回测试:</span>
                    <select id="onoff_loop" name="onoff_loop" style="width: 25%">
                        <option value="1">远端还回</option>
                        <option value="3">本地还回</option>
                        <!-- <option value="0">停止</option> -->
                    </select>
                    <input type="text" name="onoff_loop_old" id="onoff_loop_old" value="{{.OnoffLoop}}"  style="display:none;"> 
                    <input type="button" onClick="checkLoop(this)" class="btn" value="测试" id="onoff_loop_old_btn">

                </div>
                <div class="input">
                    <span>声音测试:</span>
                    <span style="width: 25%;text-align: left;padding-left:10px;">应能听到连续声音的输出</span>
                    <input type="text" name="onoff_sound" id="onoff_sound" value="{{.OnoffSound}}"  style="display:none;"> 
                    <input type="button" onClick="checkAudio(this)" class="btn" value="测试" id="onoff_sound_btn">

                </div>
                <div class="input">
                    <span>ping测试:</span>
                    <input type="text" name="ip" id="" value=""  style="width: 25%"> 
                    <input type="text" name="ip_off"  value="{{.OnoffIp}}"  style="display:none;"> 
                    <input type="button" onClick="checkPing(this)" class="btn" value="测试">
                </div>
                <div class="input">
                    <div class="ping"><span>发送包:&nbsp;&nbsp;</span><span style="text-align: left;" id="sendIp">0</span></div>
                    <div class="ping"><span>接受包:&nbsp;&nbsp;</span> <span style="text-align: left;" id="recIp">0</span> </div>
                    <div class="ping"><span>丢包率:&nbsp;&nbsp; </span><span style="text-align: left;" id="lossIp">0%</span> </div>
                </div>    


            </div>
        </form>
    </div>
</div>
<script>
// 初始化定时器
var timeInter = {}
// 检查测试声音和回路测试的按钮的值
checkAudioOrLoopBtnValue()


// 测试还回  0停止 1远端环回 3本地还回
function checkLoop(event) {  
    emptyAlertMsg()
    if($(event).hasClass('disabled')){
        return false
    }
    $(event).addClass('disabled')
    let data = {}
    let loopOff = $('#onoff_loop').val()
    // 0停止 1远端环回 3本地还回
    data['onoff_loop'] = loopOff
    if ($(event).val() == '停止'){
        data['onoff_loop'] = 0
    }

    $.ajax({
        type: "get",
        url: '/config/check_loop',
        data: data,
        success: function (response) {
            let res = JSON.parse(response)
            $(event).removeClass('disabled')
            if(res['Code'] == '0'){
                // 将old更新为当前选中的option
                $('#onoff_loop_old').val(data['onoff_loop'])
                 // $('#onoff_loop_old').val() == '1'? $('#onoff_loop_old').val('0'):$('#onoff_loop_old').val('1')
                if ($(event).val() == '停止'){
                    $(event).val('测试')
                }else{
                    $(event).val('停止')
                }
                $('#successMSG').text(res['Msg'])
            }else{
                $('#errorMSG').text(res['Msg'])
            }
        },error: function () {
            $(event).removeClass('disabled')
            $('#errorMSG').text('服务端错误')
        }
    });
}
// 测试声音
function checkAudio(event) {  
    emptyAlertMsg()
    if($(event).hasClass('disabled')){
        return false
    }
    $(event).addClass('disabled')
    let data = {}
    let soundOff = $('#onoff_sound').val()
    // 0 停止 1 开始    点击停止传0 点击测试传1
    // 是1 的话你的页面应该是停止按钮   
    // 点击测试传1   再次avc_web_audioLoopGet查询到的还是1
    // soundOff == '0'? soundOff = '1' : soundOff = '0'
    if($('#onoff_sound_btn').val() == '停止'){
        data['onoff_sound'] = 0
    }else if($('#onoff_sound_btn').val() == '测试'){
        data['onoff_sound'] = 1
    }
    
    $.ajax({
        type: "get",
        url: '/config/check_audio',
        data: data,
        success: function (response) {
            let res = JSON.parse(response)
            $(event).removeClass('disabled')
            if(res['Code'] == '0'){
                 $('#onoff_sound').val() == '1'? $('#onoff_sound').val('0'):$('#onoff_sound').val('1')
                if ($(event).val() == '停止'){
                    $(event).val('测试')
                }else{
                    $(event).val('停止')
                }
                $('#successMSG').text(res['Msg'])
            }else{
                $('#errorMSG').text(res['Msg'])
            }
        },error: function () {
            $(event).removeClass('disabled')
            $('#errorMSG').text('服务端错误')
        }
    });
}





// 测试ping
function checkPing(event) {  
    let ipReg = /^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$/
    let ipValue = $('input[name="ip"]').val()
    if(!ipReg.test(ipValue)){
        $('#errorMSG').text('ip地址格式错误')
        return false
    }
    emptyAlertMsg()
    if($(event).hasClass('disabled')){
        return false
    }
    $(event).addClass('disabled')
    let data = {}
    data['ip'] = ipValue
    // 0 停止 1 开始    初始时候值为0
    let ipOffValue = $('input[name="ip_off"]').val()
    ipOffValue == '0'? ipOffValue = '1' : ipOffValue = '0'
    data['ip_off'] = ipOffValue
    $.ajax({
        type: "get",
        url: "/config/check_ping",
        data: data,
        success: function (response) {
            let res = JSON.parse(response)
            $(event).removeClass('disabled')
            console.log(res)
            console.log(res)
            console.log(res)
            if(res['Code'] == '0'){
                // 改变ip_off的值
                $('input[name="ip_off"]').val() == '0'? $('input[name="ip_off"]').val('1'):$('input[name="ip_off"]').val('0')
                if ($(event).val() == '停止'){
                    $(event).val('测试')
                    // 结束定时器
                    clearInterval(timeInter)
                }else{
                    $(event).val('停止')
                    // 开始定时器
                    timeInter = setInterval(() => {
                        getPingPrintData()
                    }, 1000);
                }
                $('#successMSG').text(res['Msg'])
            }else{
                $('#errorMSG').text(res['Msg'])
            }
        },error: function () {
            $(event).removeClass('disabled')
            $('#errorMSG').text('服务端错误')
        }
    });
}


function getPingPrintData() {  
    $.ajax({
        type: "get",
        url: "/config/check_ping_stat",
        success: function (response) {
            // {"recv_packet":0,"send_packet":20,"loss_packet":100}
            let res = response
            $('#sendIp').text(res['send_packet'])
            $('#recIp').text(res['recv_packet'])
            $('#lossIp').text(res['loss_packet'] +'%')
        }
    });
}
function checkAudioOrLoopBtnValue() {  
    if($('#onoff_sound').val() == '0'){
        $('#onoff_sound_btn').val('测试')
    }else if($('#onoff_sound').val() == '1'){
        $('#onoff_sound_btn').val('停止')
    }
    // 11.01修改
    // if($('#onoff_sound').val() == '0'){
    //     $('#onoff_sound_btn').val('停止')
    // }else if($('#onoff_sound').val() == '1'){
    //     $('#onoff_sound_btn').val('测试')
    // }

    if($('#onoff_loop_old').val() == '0'){
        $('#onoff_loop_old_btn').val('测试')
    }else{
        // 设置option选中
        $('#onoff_loop').val('{{.OnoffLoop}}')
        $('#onoff_loop_old_btn').val('停止')
    }
}


</script>
{{end}}`
