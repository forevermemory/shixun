package tpl

const V35 = `
{{define "content"}}
 <div id="version">
    <div>
        <h3>V35软件页面</h3>
 
        <form action="/v35" method="post" id="v35Form">
            <div class="well video">

                <div class="input">
                    <span>V35通道:</span>
                    <input type="radio" name="v35_channel" id="" value="0"> <label class="white">上通道</label>
                    <input type="radio" name="v35_channel" id="" value="1"> <label class="white">下通道</label>
                </div>
                <div class="input" id="normalMode">
                    <span>V35版本:</span>
                    <input type="radio" name="v35_version" id="" value="0"> <label class="white">干&nbsp;&nbsp;&nbsp;&nbsp;线</label>
                    <input type="radio" name="v35_version" id="" value="1"> <label class="white">动中通</label>

                </div>
                <div class="input hidden" id="TSMode">
                    <span>V35版本:</span>
                    <input type="radio"  id="" checked="checked"> <label class="white">TS模式</label>

                </div>
                <!-- 下面两个不接收  只需要修改上面两个 下面两个在修改完之后重新查询一次 -->
                <div class="input">
                    <span>V35发送时钟:</span>
                    <input type="radio" name="v35_send_clock" id="" value="0" > <label class="white">上升沿</label>
                    <input type="radio" name="v35_send_clock" id="" value="1"> <label class="white">下降沿</label>
                </div>
                <div class="input">
                    <span>V35接收时钟:</span>
                    <input type="radio" name="v35_rec_clock" id="" value="0"> <label class="white">上升沿</label>
                    <input type="radio" name="v35_rec_clock" id="" value="1"> <label class="white">下降沿</label>
                </div>

            </div>
        </form>
    </div>
</div>
<script >
v35GetInfo()
// ajax  get   v35_info
function v35GetInfo(){
    $.ajax({    
        type: "get",
        url: "/v35_info",
        success: function (response) {
            emptyAlertMsg()
            var res = JSON.parse(response)
            var v35Channel = res['V35Channel']
            var v35Version = res['V35Version']
            var v35SendClock = res['V35SendClock']
            var v35ReceiveClock = res['V35ReceiveClock']
            var connetMode = res['ConnetMode']

            if (connetMode == '1'){
                $('#TSMode').removeClass('hidden')
                $('#normalMode').addClass('hidden')
            }
            // 更新选中的状态
            $.each($('input[type="radio"][name="v35_channel"]'), function (i, val) { 
                if($(val).val() == v35Channel){
                    $(val).attr('checked','checked')
                }
            });
            $.each($('input[type="radio"][name="v35_version"]'), function (i, val) { 
                if($(val).val() == v35Version){
                    $(val).attr('checked','checked')
                }
            });
            $.each($('input[type="radio"][name="v35_send_clock"]'), function (i, val) { 
                if($(val).val() == v35SendClock){
                    $(val).attr('checked','checked')
                }
            });
            $.each($('input[type="radio"][name="v35_rec_clock"]'), function (i, val) { 
                if($(val).val() == v35ReceiveClock){
                    $(val).attr('checked','checked')
                }
            });
        },error: function () {
            $('#errorMSG').text('服务端错误,获取v35信息失败')
        }
    });
}

// v35通道 修改
$('input[name="v35_channel"]').change(function (e) { 
    setV35VersionAndChannel('v35_channel')
});
// v35版本修改
$('input[name="v35_version"]').change(function (e) { 
    setV35VersionAndChannel('v35_version')
});

function setV35VersionAndChannel(v35){
    emptyAlertMsg()
    var v35_value = $('input[name="'+v35+'"]:checked').val()
    var data = {}
    data[v35] = v35_value
    $.ajax({
        type: "post",
        url: "/v35_channel_version",
        data: data,
        success: function (response) {
            var res = JSON.parse(response)
            if(res['Code'] == '0'){
                $('#successMSG').text("设置成功")
                // 更新时钟的选中情况
                $.each($('input[type="radio"][name="v35_send_clock"]'), function (i, val) { 
                    if($(val).val() == res['sendClock']){
                        $(val).attr('checked','checked')
                    }
                });
                $.each($('input[type="radio"][name="v35_rec_clock"]'), function (i, val) { 
                    if($(val).val() == res['receiveClock']){
                        $(val).attr('checked','checked')
                    }
                });

            }
        },error: function () {
            $('#errorMSG').text('服务端错误')
        }
    });
}


</script>
{{end}}`
