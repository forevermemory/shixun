package tpl

const CONFIG_SUBTITLE = `
{{define "content"}}
 <div id="version">
    <div>
        <h3>设备状态</h3>
        <ul >
            <li>
                <a href="/config/network">网络</a>
            </li>
            <li >
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
            <li  class="active">
                <a href="/config/subtitle">字幕</a>
            </li>
        </ul>
        <form action="/config/subtitle" method="post" id="subtitleForm">
            <div class="well video">

                <div class="input">
                    <span>字幕内容:</span>
                    <input type="text" name="enc_content" id="" value="{{.SubEncContent}}"> 
                    <input type="text" name="enc_content_old" id=""  value="{{.SubEncContent}}" style="display:none;"> 
                </div>
                <div class="input">
                    <span>字幕位置:</span>
                    <input type="radio" name="enc_position" id="" value="0"> <label class="white">自定义</label>
                    <input type="radio" name="enc_position" id="" value="1"> <label class="white">左上</label>
                    <input type="radio" name="enc_position" id="" value="2"> <label class="white">左下</label>
                    <input type="radio" name="enc_position" id="" value="3"> <label class="white">中上</label>
                    <input type="radio" name="enc_position" id="" value="4"> <label class="white">中下</label>
                    <input type="radio" name="enc_position" id="" value="5"> <label class="white">右上</label>
                    <input type="radio" name="enc_position" id="" value="6"> <label class="white">右下</label>
                    <input type="text" name="enc_position_old" id=""  value="{{.SubEncPosition}}" style="display:none;"> 

                </div>
                
                <div class="input {{if gt .SubEncPosition 0 }}  hidden{{else}}  {{end}}" style="color: white;" id="coordinate">
                    <span></span>
                        x坐标:<input type="text" name="enc_px" id="" value="{{.SubEncPx}}" style="width: 30px;">
                        y坐标:<input type="text" name="enc_py" id="" value="{{.SubEncPy}}" style="width: 30px;"> 
                        <input type="text" name="enc_px_old" id=""  value="{{.SubEncPx}}" style="display:none;"> 
                        <input type="text" name="enc_py_old" id=""  value="{{.SubEncPy}}" style="display:none;"> 
                        
                </div>

                <div class="input">
                    <span>字体大小:</span>
                    <input type="text" name="enc_size" id="" value="{{.SubEncSize}}" > <span style="text-align: left;">(1~7)</span>
                    <input type="text" name="enc_size_old" id=""  value="{{.SubEncSize}}" style="display:none;"> 

                </div>
                <div class="input">
                    <span>字体颜色:</span>
                    <input type="radio" name="enc_color" id="" value="0"> <label class="white">白</label>
                    <input type="radio" name="enc_color" id="" value="1"> <label class="white">黑</label>
                    <input type="radio" name="enc_color" id="" value="2"> <label class="white">红</label>
                    <input type="radio" name="enc_color" id="" value="3"> <label class="white">黄</label>
                    <input type="radio" name="enc_color" id="" value="4"> <label class="white">蓝</label>
                    <input type="radio" name="enc_color" id="" value="5"> <label class="white">绿</label>
                    <input type="text" name="enc_color_old" id=""  value="{{.SubEncColor}}" style="display:none;"> 

                </div>
                <div class="input">
                    <span>显示/隐藏:</span>
                    <input type="radio" name="enc_onoff" id="" value="0"> <label class="white">隐藏</label>
                    <input type="radio" name="enc_onoff" id="" value="1"> <label class="white">显示</label>
                    <input type="text" name="enc_onoff_old" id=""  value="{{.SubEncOnoff}}" style="display:none;"> 

                </div>
                <input type="button" onClick="checkValid()" class="btn" value="修改">
            </div>
        </form>
    </div>
</div>

<script>
// 初始化 字幕位置 颜色 显示和隐藏 的选择情况
$.each($('input[type="radio"][name="enc_position"]'), function (i, val) { 
    if($(val).val() == '{{.SubEncPosition}}'){
        $(val).attr('checked','checked')
    }
});
$.each($('input[type="radio"][name="enc_color"]'), function (i, val) { 
    if($(val).val() == '{{.SubEncColor}}'){
        $(val).attr('checked','checked')
    }
});
$.each($('input[type="radio"][name="enc_onoff"]'), function (i, val) { 
    if($(val).val() == '{{.SubEncOnoff}}'){
        $(val).attr('checked','checked')
    }
});

// 判断左边的显示和隐藏
$('input[name="enc_position"]').change(function (event) { 
    let subtitlePosition =  $('input[name="enc_position"]:checked').val()
    if(subtitlePosition == '0'){
        $('#coordinate').removeClass('hidden')
    }else{
        $('#coordinate').addClass('hidden')
    }
});


// 提交前检查内容是否合理
function checkValid() {  
    emptyAlertMsg()
    let content = $('input[name="enc_content"]').val()
    let size = $('input[name="enc_size"]').val()
    let px = $('input[name="enc_px"]').val()
    let py = $('input[name="enc_py"]').val()

    let position = $('input[name="enc_position"]:checked').val()
    let color = $('input[name="enc_color"]:checked').val()
    let onoff = $('input[name="enc_onoff"]:checked').val()

    if (content.trim()==''){
        $('#errorMSG').text('字幕内容不能为空')
        return false
    }
    if (position == undefined){
        $('#errorMSG').text('请选择字幕位置')
        return false
    }
    if (position == '0'){
        if((px.trim()=='' )|| (py.trim()=='')){
            $('#errorMSG').text('请填写字幕自定义坐标')
            return false
        }
    }
    if(!/^[0-9]+$/.test(px)){
        $('#errorMSG').text('请正确填写字幕坐标')
        return false
    }
    if(!/^[0-9]+$/.test(py)){
        $('#errorMSG').text('请正确填写字幕坐标')
        return false
    }
    if(!/[0-9]/.test(size)){
        $('#errorMSG').text('请正确填写字幕')
        return false
    }

    if(size>7 || size<1){
        $('#errorMSG').text('请正确填写字幕范围')
        return false
    }
    if (color == undefined){
        $('#errorMSG').text('请选择字幕颜色')
        return false
    }
    if (onoff == undefined){
        $('#errorMSG').text('请选择是否隐藏')
        return false
    }
    $('#subtitleForm').submit()
}
</script>
{{end}}`
