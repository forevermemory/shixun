package tpl

const CONFIG_VERSION = `
{{define "content"}}
<div id="version">
    <div>
        <h3>设备状态</h3>
        <ul >
            <li>
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
            <li  class="active">
                <a href="/config/version">版本</a>
            </li>
            <li>
                <a href="/config/subtitle">字幕</a>
            </li>
        </ul>
        <form action="" method="post">
            <div class="well video">
                <div class="input">
                    <span>软件版本:</span>
                    <span style="width: 25%;text-align: left;padding-left:10px;">{{.HW}}</span>
                </div>
                <div class="input">
                    <span>硬件版本:</span>
                    <span style="width: 25%;text-align: left;padding-left:10px;">{{.SW}}</span>
                </div>

                <div class="input">
                    <span>显控版本:</span>
                    <span style="width: 25%;text-align: left;padding-left:10px;">{{.Panel}}</span>
                </div>
                 <div class="input">
                    <span>V35版本:</span>
                    <span style="width: 25%;text-align: left;padding-left:10px;">{{.V35}}</span>
                </div>
                <div class="input">
                    <span>MCU版本:</span>
                    <span style="width: 25%;text-align: left;padding-left:10px;">{{.Mcu}}</span>
                </div>
                <div class="input">
                    <span>FPGA版本:</span>
                    <span style="width: 25%;text-align: left;padding-left:10px;">{{.Fpga}}</span>
                </div>                        

            </div>
        </form>
    </div>
</div>
{{end}}`
