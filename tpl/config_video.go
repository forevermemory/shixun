package tpl

const CONFIG_VIDEO = `
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
            <li class="active">
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
        <form action="/config/video" method="post" id="videoCommit">
            <div class="well video">
                <div class="input">
                    <span>输入源:</span>
                    <input type="radio" name="in_param" id="" value="0"> <label class="white">标清</label>
                    <input type="radio" name="in_param" id="" value="1"> <label class="white">HDMI-720P</label>
                    <input type="radio" name="in_param" id="" value="2"> <label class="white">HDMI-1080P</label>
                    <input type="radio" name="in_param" id="" value="3"> <label class="white">SDI-720P</label>
                    <input type="radio" name="in_param" id="" value="4"> <label class="white">SDI-1080P</label>
                    <input type="text" name="in_param_old" id="" value="{{.VideoInParam}}"  style="display:none;"> 
                </div>
                <div class="input">
                    <span>输出源:</span>
                    <input type="radio" name="out_param" id="" value="0"> <label class="white">标清</label>
                    <input type="radio" name="out_param" id="" value="1"> <label class="white">HDMI-720P@50</label>
                    <input type="radio" name="out_param" id="" value="2"> <label class="white">HDMI-720P@60</label>
                    <input type="radio" name="out_param" id="" value="3"> <label class="white">HDMI-1080P@50</label>
                    <input type="radio" name="out_param" id="" value="4"> <label class="white">HDMI-1080P@60</label>
                    <input type="text" name="out_param_old" id="" value="{{.VideoOutParam}}"  style="display:none;"> 
                </div>
                <div class="input">
                    <span>输出制式:</span>
                    <input type="radio" name="std" id="" value="0"> <label class="white">NTSC</label>
                    <input type="radio" name="std" id="" value="1"> <label class="white">PAL</label>
                    <input type="text" name="std_old" id="" value="{{.VideoStd}}"  style="display:none;"> 

                </div>
                <div class="input">
                    <span>视频编码:</span>
                    <input type="radio" name="enc_frame" id="" value="0"> <label class="white">H265</label>
                    <input type="radio" name="enc_frame" id="" value="1"> <label class="white">H264</label>
                    <input type="text" name="enc_frame_old" id="" value="{{.VideoEncFrame}}"  style="display:none;"> 

                </div>


                <input type="button"  onClick="doSubmitForm()" class="btn" value="修改" style="margin-left: 30px;">
            </div>
        </form>
    </div>
</div>
<script>

// 初始化视频 输入源 输出源 输出制式 视频编码 的选择情况
$.each($('input[type="radio"][name="in_param"]'), function (i, val) { 
    if($(val).val() == '{{.VideoInParam}}'){
        $(val).attr('checked','checked')
    }
});
$.each($('input[type="radio"][name="out_param"]'), function (i, val) { 
    if($(val).val() == '{{.VideoOutParam}}'){
        $(val).attr('checked','checked')
    }
});
$.each($('input[type="radio"][name="std"]'), function (i, val) { 
    if($(val).val() == '{{.VideoStd}}'){
        $(val).attr('checked','checked')
    }
});
$.each($('input[type="radio"][name="enc_frame"]'), function (i, val) { 
    if($(val).val() == '{{.VideoEncFrame}}'){
        $(val).attr('checked','checked')
    }
});

// video表单提交前判断 是否有参数发生变化
function doSubmitForm() {  
    emptyAlertMsg()
    let inParamOld = $('input[name="in_param_old"]').val()
    let outParamOld = $('input[name="out_param_old"]').val()
    let  stdOld= $('input[name="std_old"]').val()
    let encFrameOld = $('input[name="enc_frame_old"]').val()

    
    let inParam = $('input[name="in_param"]:checked').val()
    let outParam = $('input[name="out_param"]:checked').val()
    let std = $('input[name="std"]:checked').val()
    let encFrame = $('input[name="enc_frame"]:checked').val()

    if((inParam == inParamOld)&&(outParam == outParamOld)&&(stdOld == std)&&(encFrameOld == encFrame)){
        $('#errorMSG').text('无参数变化,提交失败')
        return false
    }
    
    $('#videoCommit').submit()
}

</script>
{{end}}`
