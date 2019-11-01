## 准备及编译步骤
执行 go generate 打包静态文件
statik -src=static
```bash
sudo dpkg --add-architecture i386
sudo apt update
sudo apt install zlib1g:i386
env CGO_ENABLED=1 GOARCH=arm64 GOARM=6 CC=aarch64-himix100-linux-gcc go build  -ldflags='-s -w' -vx .
env CGO_ENABLED=1 GOARCH=arm64 CC=aarch64-himix100-linux-gcc go build  -ldflags='-s -w' -x .
```
