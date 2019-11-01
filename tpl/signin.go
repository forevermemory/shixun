package tpl

const SIGNIN = `
<!doctype html>
<html lang="zh">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>登录</title>
    <link rel="stylesheet" href="/static/css/reset.css">
    <link rel="stylesheet" href="/static/css/login.css">
</head>

<body>
    <div class="login" id="login">
        <div style="text-align: center;padding-top: 60px;">
            <span style="color:red;">{{.MsgError}}</span>
            <span style="color:green;">{{.MsgSuccess}}</span>
        </div>
        <div class="cell">
            <form action="/signin" method="post">
                <div class="header">通讯管理系统</div>
                <div class="first">
                    <span>管理员：</span>
                    <input type="text" name="username" autocomplete="off">
                </div>
                <div class="second">
                    <span>密码：</span>
                    <input type="password" name="password" autocomplete="off">
                </div>
                <div class="third">
                    <input type="checkbox" name="remember_me" value="yes" checked>
                    <span>记住我</span>
                </div>
                <div class="submit">
                    <input type="submit" value="登录">
                </div>
            </form>
        </div>
    </div>
</body>
<script>
window.onload=function(){
    var getDOM = function(id) {
        return document.getElementById(id)
    }
    var sheight = document.documentElement.clientHeight || document.body.clientHeight
    window.onresize = function() {
        getDOM('login').style.height = sheight + 'px'
    }
    getDOM('login').style.height = sheight + 'px'
}
</script>
</html>`
