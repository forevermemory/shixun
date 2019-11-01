package tpl

const LAYOUT = `
{{define "layout"}}
<!doctype html>
<html lang="zh">

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>265 设备</title>
  <link rel="stylesheet" href="/static/css/reset.css">
  <link rel="stylesheet" href="/static/css/content.css">
  <script type="text/javascript" src="/static/js/jquery.min.js"></script>
  <script type="text/javascript" src="/static/js/slider.js"></script>
  <script type="text/javascript" src="/static/js/jquery.form.js"></script>
  <script>
  // 全局工具函数
// 清空提示信息
function emptyAlertMsg() {  
    $('#successMSG').text('')
    $('#errorMSG').text('')
}
</script>
</head>

<body>

 <div class="content" id="content">
    <div class="header" style="position: relative;height: 95px">
        <span style="color:red;position: absolute;left: 30.5%;top: 45px;" id="errorMSG">{{.MsgError}}</span>
        <span style="color:green; position: absolute;left: 30.5%;top: 45px;" id="successMSG">{{.MsgSuccess}}</span>
        <a href="/signout" class="signout">退出</a>
        <form action="/reboot" method="post">
            <input type="button" value="重启"  onclick="reboot()">
        </form>
    </div>
    <div class="main">
        <div class="menu">
            <ul>
                <li {{if .MenuEncodeDecode}}class="active"  {{end}}>
                    <a href="/config/network">
                        编解码器<p>配置参数</p>
                    </a>
                </li>
                <li {{if .MenuV35}}class="active"  {{end}}>
                    <a href="/v35">
                        v35协转<p>配置参数</p>
                    </a>
                </li>
                <li {{if .MemuUpgrade}}class="active"  {{end}}>
                    <a href="/upgrade">
                        升<span style="visibility: hidden;">升级</span>级<p>固件升级</p>
                    </a>
                </li>
                <li {{if .MemuReboot}}class="active"  {{end}}>
                    <a href="javascript:viod(0)" onclick="reboot()">
                        系<span style="visibility: hidden;">升级</span>统<p>系统重启</p>
                    </a>
                </li>
            </ul>
        </div>
        <!--  -->
        {{template "content" .}}
        <!--  -->
    </div>
</div>

  <script>

function reboot() {  
    var r = confirm("您确认进行系统重启吗!");
    if (r==true){
        // 确认重启
      $.ajax({
            type: "get",
            url: "/reboot",
            success: function (response) {
                let res = JSON.parse(response)
                if(res['Code'] == '0'){
                    window.location.href = '/signout'
                }else{
                    $('#errorMSG').text('重启失败')
                }
            },error: function () {
                $('#errorMSG').text('服务端错误')
            }
        });
    }else{
        return false;
    }
}

  </script>
</body>

</html>
{{end}}`
