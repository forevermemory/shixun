package tpl

const REBOOT = `
{{define "content"}}
<div id="version">
	<h3>265 设备重启</h3>

    <div class="progress" style="width: 60%;margin-top: 150px;margin-left: 1%;position: relative; display: none;">
        <div style="height: 16px;background-color: #dedede;"></div>
        <div class="bar" style="height: 16px;width:0%;top: 0px;background: #9f7d7d;position: absolute;z-index:10;"></div>
        <span class="time" style="position: absolute; color: #fff; left: 0;top: -25px;">正在重启中,请稍后</span>
    </div>
    <div class="uploading" style="display:none;left: 45%; top: 30%; position: absolute; z-index: 100000;">
        <img src="/static/img/loading.svg" style="width: 70%"/>
    </div>
    <div class="upgrate" style="left: 40%; position: absolute; z-index: 100000;"></div>


</div>

<script>
 $('.progress').css('display', 'block')
$.ajax({
    type: "get",
    url: "/reboot_start",
    success: function (response) {
        let res = JSON.parse(response)
        if(res['Code'] == '0'){
            var w = 0;
			var timer = setInterval(function(){
			        w++
			        $(".bar").css("width", (w + "%"))
			        $(".time").text("正在重启中,请稍后"+getDotNum(w))
			        if (w>99){
			        	w = 0
			        }
			        // 继续调用测试接口
			        getRebootTest((res)=>{
			        	if(res){
				            clearInterval(timer)
				            window.location.href = '/signin'
			        	}
			        })
			},1000)
        }else{
            $('#errorMSG').text('重启失败')
        }
    },error: function () {
        $('#errorMSG').text('服务端错误')
    }
});



function getDotNum(w){
	let len = w % 5 +1
	let dots = ""
	for (let i = 0; i < len; i++) {
		dots +="。"
	}
	return dots
}

function getRebootTest(callback) {  
    $.ajax({
        type: "get",
        url: "/reboot_end",
        success: function (response) {
            let res = JSON.parse(response)
            if(res['Code'] == '0'){
                callback(true)
            }
        },error: function () {
            callback(false)
        }
    });
}
</script>
{{end}}`
