#!/bin/bash

echo "安装GO"
sudo pacman -S go

echo "配置GO环境"
git config --global http.proxy 'socks5://127.0.0.1:7890'
git config --global https.proxy 'socks5://127.0.0.1:7890'
export https_proxy=http://127.0.0.1:7890;export http_proxy=http://127.0.0.1:7890;export all_proxy=socks5://127.0.0.1:7890
export GOPROXY=https://proxy.golang.com.cn,direct
go env -w GO111MODULE=auto

go get github.com/shirou/gopsutil/process
go get github.com/esiqveland/notify
go get github.com/godbus/dbus

echo "编译msg-wechat"
cd ./src
#go run ./msg-wechat.go
go build msg-wechat.go
echo "编译完成"
