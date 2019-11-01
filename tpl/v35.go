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
                    <input type="text" name="v35_channel_old"  value="{{.V35Channel}}"  style="display:none;"> 
                </div>
                <div class="input">
                    <span>V35版本:</span>
                    <input type="radio" name="v35_version" id="" value="0"> <label class="white">干线</label>
                    <input type="radio" name="v35_version" id="" value="1"> <label class="white">动中通</label>
                    <input type="text" name="v35_version_old"  value="{{.V35Version}}"  style="display:none;"> 

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

                <input type="button" onClick="v35Change()" class="btn" value="修改">
            </div>
        </form>
    </div>
</div>
<script >
    // 初始化radio的选择情况
    $.each($('input[type="radio"][name="v35_channel"]'), function (i, val) { 
        if($(val).val() == '{{.V35Channel}}'){
            $(val).attr('checked','checked')
        }
    });
    $.each($('input[type="radio"][name="v35_version"]'), function (i, val) { 
        if($(val).val() == '{{.V35Version}}'){
            $(val).attr('checked','checked')
        }
    });
    $.each($('input[type="radio"][name="v35_send_clock"]'), function (i, val) { 
        if($(val).val() == '{{.V35SendClock}}'){
            $(val).attr('checked','checked')
        }
    });
    $.each($('input[type="radio"][name="v35_rec_clock"]'), function (i, val) { 
        if($(val).val() == '{{.V35ReceiveClock}}'){
            $(val).attr('checked','checked')
        }
    });
// 提交前检查内容是否合理
function v35Change() {  

    emptyAlertMsg()
    let v35Channel = $('input[name="v35_channel"]:checked').val()
    let v35Version = $('input[name="v35_version"]:checked').val()
    if (v35Channel == undefined){
        $('#errorMSG').text('请选择V35通道')
        return false
    }
    if (v35Version == undefined){
        $('#errorMSG').text('请选择V35版本')
        return false
    }

    if(($('input[name="v35_channel_old"]').val() ==v35Channel) && ($('input[name="v35_version_old"]').val() == v35Version)){
        $('#errorMSG').text('未发生变化,请更改后再提交')
        return false
    }
    $('#v35Form').submit()
}
</script>
{{end}}`
