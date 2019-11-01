var getDOM = function (id) {
    return document.getElementById(id)
}
var sheight = document.documentElement.clientHeight || document.body.clientHeight
window.onresize = function () {
    getDOM('content').style.height = sheight + 'px'
}
getDOM('content').style.height = sheight + 'px'


$(function () {
    // if($(".network").length > 0) {
    //     autoChange()
    //     $("#contmode").change(function(){
    //         autoChange()
    //     })

    //     multiChange()
    //     $("#recv_multi").change(function(){
    //         multiChange()
    //     })

    //     if($('#recv_multi').val() != '0') {
    //         flag()
    //     }

    //     $("#recv_addr").keyup(function(){
    //         if($('#recv_addr').val().split('.').length == 4 && $('#recv_addr').val().split('.')[3]) {
    //             flag()
    //         }
    //     })
    // }

    // if($(".subtitle").length > 0) {
    //     init()
    //     $("#cont").keyup(function(){
    //         init()
    //     })
    // }

    // if(!Array.prototype.indexOf){
    //     Array.prototype.indexOf = function(val){
    //         var value = this;
    //         for(var i =0; i < value.length; i++){
    //             if(value[i] == val) return i;
    //         }
    //         return -1;
    //     };
    // }

    // if($('.video').length > 0) {
    //     isNone()
    //     $("#video_rate").change(function(){
    //         isNone()
    //     })

    //     outChnageOption()
    //     $("#out_defin").change(function(){
    //         outChnageOption()
    //     })

    //     inChnageOption()
    //     $("#in_defin").change(function(){
    //         inChnageOption()
    //     })
    //     $("#rcmode").change(function(){
    //         inChnageOption()
    //     })
    // }

    // if($(".upgrade").length > 0) {
    //     $("#upload").click(function(){
    //         var formData = new FormData($('#upgrade')[0]);
    //         $.ajax({
    //             url:"/upgrade",
    //             type:"post",
    //             enctype: 'multipart/form-data',
    //             data: formData,
    //             cache: false,
    //             processData: false,
    //             contentType: false,
    //             beforeSend: function() {
    //                 $('.uploading').css('display', 'block');
    //                 $('.progress_form').css('display','none')
    //             },
    //             success: function(data) {
    //                 $('#main').html('')
    //                 $('.uploading').css('display', 'none');
    //                 if (data.result == 0) {
    //                     $('.progress_form').css('display','block')
    //                     $('#main').prepend('<span style="color:red;">升级失败</span>')
    //                 } else {
    //                     $.progress(20000);
    //                     setTimeout(function(){
    //                         window.location.href="/upgrade"
    //                     },20000)
    //                     $('#main').prepend('<span style="color:green;">升级成功</span>')
    //                 }
    //             }
    //         })
    //     })
    // }

    $('.pure-menu').click(function () {
        $('.k-config-submenu').css({ 'display': 'none' })
        $('.pure-menu').removeClass('active')
        $('.pure-menu').removeClass('in');
        $(this).addClass('active')
        $(this).siblings('ul').css('display', 'block')
    })
    $('.pure-menu').each(function () {
        if ($(this).attr('href') == window.location.pathname) {
            $(this).addClass('active')
        }
    });
    $('.k-config-submenu li a').each(function () {
        if ($(this).attr('href') == window.location.pathname) {
            $('.k-config-submenu li a').removeClass('in');
            $(this).addClass('in');
            $(this).parents('.k-config-submenu').css('display', 'block');
            $(this).parents('.pure-menu-list').siblings('.pure-menu').addClass('active')
        }
    });
    // if ( window.history.replaceState ) {
    //   window.history.replaceState( null, null, window.location.href );
    // }
})

function autoChange() {
    if ($('#contmode').val() == '0') {
        $('#recv_multi').attr('disabled', 'disabled')
    } else {
        $('#recv_multi').removeAttr('disabled', '')
    }
}

function multiChange() {
    if ($('#recv_multi').val() == '0') {
        $('#recv_addr').attr('disabled', 'disabled')
    } else {
        $('#recv_addr').removeAttr('disabled', '')
    }
}

function flag() {
    if ($('#recv_addr').val().split('.')[0] < 224 || $('#recv_addr').val().split('.')[0] > 239) {
        $('.span').css('display', 'block')
    } else {
        $('.span').css('display', 'none')
    }
}

function init() {
    var len = 30 - $("#cont").val().replace(/[^\x00-\xff]/g, "***").length
    $("#name").text("你还可以输入" + len + "个字节或" + Math.floor(len / 3) + "个汉字")
    if (len <= 0) {
        len = 0
        $("#name").css('color', 'red')
    } else {
        $("#name").css('color', 'white')
    }
}

function isNone() {
    var arr = [512, 640, 768, 832, 896, 960, 1024, 1152, 1536, 1792, 1920, 2048, 3072, 4096]
    if (arr.indexOf(Number($('#video_rate').val())) >= 0) {
        $('.isNone').css('display', 'none')
    } else {
        $('.isNone').css('display', 'inline')
        $('#video_rate').val(0)
    }
}

function outChnageOption() {
    if ($("#out_defin").val() == '0') {
        $('#out_resolution').attr('disabled', 'disabled')
    } else {
        $('#out_resolution').removeAttr('disabled', '')
    }
}

function inChnageOption() {
    if (($("#in_defin").val() == '0') || ($("#rcmode").val() == '1')) {
        $('#in_resolution').attr('disabled', 'disabled')
    } else {
        $('#in_resolution').removeAttr('disabled', '')
    }
}