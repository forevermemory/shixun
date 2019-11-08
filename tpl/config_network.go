package tpl

const CONFIG_NETWORK = `
{{define "content"}}
<div id="version">
    <div>
        <h3 id="networkH3">设备状态</h3>
        <ul id="networkUl">
            <li class="active">
                <a href="/config/network">网络</a>
            </li>
            <li>
                <a href="/config/audio">音频</a>
            </li>
            <li>
                <a href="/config/video">视频</a>
            </li>
            <li>
                <a href="/config/check">诊断</a>
            </li>
            <li>
                <a href="/config/version">版本</a>
            </li>
            <li>
                <a href="/config/subtitle">字幕</a>
            </li>
        </ul>
        <form action="/config/network" method="post">
            <div class="well video">
                <div class="input">
                    <span>连接模式：</span>
                    <input type="radio" name="connect_mode"  value="1" id=""> <label class="white">直通模式</label>
                    <input type="radio" name="connect_mode"  value="2" id=""> <label class="white">TS模式</label>
                    <input type="radio" name="connect_mode"  value="3" id=""> <label class="white">RTP模式</label>
                    <input type="text" name="connect_mode_old"  value="{{.ConnetMode}}" id="" style="display:none;"> 
                </div>
                <div class="input">
                    <span>速率：</span>
                    <input type="text" name="encode_rate" value="{{.VideoEncodeRate}}"> <span>kbps(64~1920)</span>
                    <input type="text" name="encode_rate_old" value="{{.VideoEncodeRate}}" style="display:none;"> 
                </div>
                <div class="input">
                    <span></span>
                    <input type="button" onClick="networkModeRateUpdate(this)" class="btn " value="修改" style="width:80px;">
                    <span class="white" style="cursor:pointer" id="showConnectInfo" data-type="1">显示连通信息</span>
                </div>

                <div class="hidden" id="networkShowSendInfo">
                        <div class="input" style="padding: 0">
                            <div class="ping" ><span>发送视频速率:&nbsp;</span><span style="text-align: left;" id="showNetInfo1">0</span></div>
                            <div class="ping" ><span>发送视频丢包率:&nbsp;</span> <span style="text-align: left;" id="showNetInfo2">0</span> </div>
                        </div>    
                        <div class="input" style="padding: 0">
                            <div class="ping"><span>发送音频速率:&nbsp;</span><span style="text-align: left;" id="showNetInfo3">0</span></div>
                            <div class="ping"><span>发送音频丢包率:&nbsp;</span> <span style="text-align: left;" id="showNetInfo4">0</span> </div>
                        </div>    
                        <div class="input" style="padding: 0">
                            <div class="ping"><span>接收视频速率:&nbsp;</span><span style="text-align: left;" id="showNetInfo5">0</span></div>
                            <div class="ping"><span>接收视频丢包率:&nbsp;</span> <span style="text-align: left;" id="showNetInfo6">0</span> </div>
                        </div>    
                        <div class="input" style="padding: 0">
                            <div class="ping"><span>接收音频速率:&nbsp;</span><span style="text-align: left;" id="showNetInfo7">0</span></div>
                            <div class="ping"><span>接收音频丢包率:&nbsp;</span> <span style="text-align: left;" id="showNetInfo8">0</span> </div>
                        </div>    
                </div>

                <div class="" id="networkShowChangeIp">
                    <div class="input">
                        <span>本地IP地址：</span>
                        <input type="text" name="ipaddr" value="{{.IP}}">
                        <input type="text" name="ipaddr_old" value="{{.IP}}"  style="display:none;">
                    </div>
                    <div class="input">
                        <span>子网掩码：</span>
                        <input type="text" name="netmask" value="{{.Netmask}}">
                        <input type="text" name="netmask_old" value="{{.Netmask}}"  style="display:none;">
                    </div>
                    <div class="input">
                        <span>网关地址：</span>
                        <input type="text" name="gateway" value="{{.Gateway}}">
                        <input type="text" name="gateway_old" value="{{.Gateway}}"  style="display:none;">
                    </div>
                    <div class="input">
                        <span></span>
                        <input type="button" onClick="networkIPGatewayNetmaskUpdate(this)" class="btn w80" value="修改" style="width:80px;">
                    </div>
                </div>

            </div>
        </form>
    </div>
</div>

<script>

var timeInter = {}
$('#showConnectInfo').click(function (e) { 
    if ( $(this).data('type') == '1'){
        $('#networkUl').addClass('hidden')
        $('#networkH3').addClass('hidden')

        $(this).text('显示呼叫数据')
        $(this).data('type','2')
    } else if( $(this).data('type') == '2'){
        $(this).text('隐藏呼叫数据')
        $(this).data('type','3')
        // 显示呼叫数据
        $('#networkShowSendInfo').removeClass('hidden')
        $('#networkShowChangeIp').addClass('hidden')

        // 开始连续发送请求拿数据
        timeInter = setInterval(() => {
            $.ajax({
                type: "get",
                url: "/config/network_data_info",
                success: function (response) {
                    let res = JSON.parse(response)
                    $('#showNetInfo1').text(res['SendVideoRate']+'Kbps')
                    $('#showNetInfo2').text(res['SendVideoLoss']+'%')
                    $('#showNetInfo3').text(res['SendAudioRate']+'Kbps')
                    $('#showNetInfo4').text(res['SendAudioLoss']+'%')
                    $('#showNetInfo5').text(res['RecVideoRate']+'Kbps')
                    $('#showNetInfo6').text(res['RecVideoLoss']+'%')
                    $('#showNetInfo7').text(res['RecAudioRate']+'Kbps')
                    $('#showNetInfo8').text(res['RecAudioLoss']+'%')
                }
            });
        }, 1000);
    
    }else if( $(this).data('type') == '3'){
        clearInterval(timeInter)
        $(this).text('显示呼叫数据')
        $(this).data('type','2')
        // 停止上面的计时器

        $('#networkShowSendInfo').addClass('hidden')
        $('#networkShowChangeIp').removeClass('hidden')
    }

    
});


// 初始化 mode 的选择情况
$.each($('input[type="radio"][name="connect_mode"]'), function (i, val) { 
    if($(val).val() == '{{.ConnetMode}}'){
        $(val).attr('checked','checked')
    }
});



// 更新模式 速率
function networkModeRateUpdate(event) {  
    emptyAlertMsg()

    let encodeRateValue = $('input[name="encode_rate"]').val()
    if (!/^[0-9]{2,4}$/.test(encodeRateValue)){
        $('#errorMSG').text('请正确输入速率')
        $('input[name="encode_rate"]').focus()
        return false
    }
    if (encodeRateValue > 1920 || encodeRateValue <64){
        $('#errorMSG').text('请正确输入速率区间')
        $('input[name="encode_rate"]').focus()
        return false
    }

    
    if($(event).hasClass('disabled')){
        return false
    }
    $(event).addClass('disabled')
    let data = {
        'connect_mode': $('input[name="connect_mode"]:checked').val(),
        'connect_mode_old' : $('input[name="connect_mode_old"]').val(),
        'encode_rate': $('input[name="encode_rate"]').val(),
        'encode_rate_old': $('input[name="encode_rate_old"]').val(),
    }
    $.ajax({
        type: "post",
        url: "/config/network_mode_rate",
        data: data,
        success: function (response) {
            $('#successMSG').text('')
            $('#errorMSG').text('')
            $(event).removeClass('disabled')
            let res = JSON.parse(response)
            if(res['Code'] == '0'){
                // success
                $('#successMSG').text(res['MsgSuccess'])
                $('input[name="connect_mode_old"]').val($('input[name="connect_mode"]:checked').val())
                $('input[name="encode_rate_old"]').val($('input[name="encode_rate"]').val())
            }else{
                $('#errorMSG').text(res['MsgError'])
            }
        },error: function () {
            $(event).removeClass('disabled')
            $('#errorMSG').text('服务端错误')
        }
    });
}

// 更新 ip 网关 子网掩码
function networkIPGatewayNetmaskUpdate(event) {  
    emptyAlertMsg()
    // 校验ip正确性
    let ipaddrValue = $('input[name="ipaddr"]').val()
    let ipaddrValueOld = $('input[name="ipaddr_old"]').val()
    let netmaskValue = $('input[name="netmask"]').val()
    let netmaskValueOld = $('input[name="netmask_old"]').val()
    let gatewayValue = $('input[name="gateway"]').val()
    let gatewayValueOld = $('input[name="gateway_old"]').val()
    let ipReg = /^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$/
    if (!ipReg.test(ipaddrValue) || !ipReg.test(netmaskValue) || !ipReg.test(gatewayValue)){
        $('#errorMSG').text('ip地址、子网掩码、或者网关地址格式错误')    
        return false
    }

    if ((ipaddrValue == ipaddrValueOld) &&(netmaskValue == netmaskValueOld)&&(gatewayValue == gatewayValueOld)){
         $('#errorMSG').text('ip地址、子网掩码、网关地址未发生变化,请更改后进行提交 ')    
        return false
    }
    
    if($(event).hasClass('disabled')){
        return false
    }
    $(event).addClass('disabled')
    let data = {
        'ipaddr': ipaddrValue,
        'ipaddr_old': $('input[name="ipaddr_old"]').val(),
        'netmask':  netmaskValue,
        'netmask_old': $('input[name="netmask_old"]').val(),
        'gateway': gatewayValue,
        'gateway_old': $('input[name="gateway_old"]').val(),
    }
    $.ajax({
        type: "post",
        url: "/config/network_ip_netmask_gateway",
        data: data,
        success: function (response) {
            let res = JSON.parse(response)
            $(event).removeClass('disabled')
            if(res['Code'] == '0'){
                // success
                $('#successMSG').text(res['MsgSuccess'])
                $('input[name="ipaddr_old"]').val($('input[name="ipaddr"]').val())
                $('input[name="netmask_old"]').val($('input[name="netmask"]').val())
                $('input[name="gateway_old"]').val($('input[name="gateway"]').val())

            }else{
                $('#errorMSG').text(res['MsgError'])
            }
        },error: function () {
            $(event).removeClass('disabled')
            $('#errorMSG').text('服务端错误')
        }
    });
}



// 连接模式变化自动更新
$('input[name="connect_mode"]').change(function (e) { 
    let connect_mode = $('input[name="connect_mode"]:checked').val()
    console.log(connect_mode)
    let data = {
        'connect_mode': connect_mode,
        'connect_mode_old' : $('input[name="connect_mode_old"]').val(),
    }
    $.ajax({
        type: "post",
        url: "/config/network_mode_rate",
        data: data,
        success: function (response) {
            $('#successMSG').text('')
            $('#errorMSG').text('')
            $(event).removeClass('disabled')
            let res = JSON.parse(response)
            if(res['Code'] == '0'){
                // success
                $('#successMSG').text(res['MsgSuccess'])
                $('input[name="connect_mode_old"]').val($('input[name="connect_mode"]:checked').val())
            }else{
                $('#errorMSG').text(res['MsgError'])
            }
        },error: function () {
            $(event).removeClass('disabled')
            $('#errorMSG').text('服务端错误')
        }
    });
});

</script>
{{end}}`
