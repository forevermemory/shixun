package tpl

const UPGRADE = `
{{define "content"}}



 <div id="version">
    <div>
        <h3>软件升级页面</h3>
 
        <form action="/upgrade" method="post" enctype="multipart/form-data">
            <div class="well video">

                <div class="input">
                    <span>MCU固件:</span>
                    <input type="file" name="mcu_firmware1" id="" >
                </div>
                <div class="input">
                    <span></span>
                    <input type="file" name="mcu_firmware2" id="">
                </div>
                <div class="input">
                    <span>V35固件:</span>
                    <input type="file" name="v35_firmware" id="" >
                </div>

                <div class="input">
                    <span>系统固件:</span>
                    <input type="file" name="file_firmware" id="" >
                </div>

                <div class="input">
                    <span>md5文件:</span>
                    <input type="file" name="file_md5" id="" >
                </div>

               


                <input type="submit"  class="btn ml35" value="修改">
                <input type="button"  class="btn ml35" value="重置" onclick="resetUploadInput()">
            </div>
        </form>
    </div>
</div>




<script>
function resetUploadInput(){
    window.location.href = '/upgrade'
}
</script>
{{end}}`
