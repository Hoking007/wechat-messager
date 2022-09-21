#!/bin/bash

#source build.sh

echo "安装msg-wechat"
chmod +x ./src/msg-wechat
sudo /bin/cp -rf ./src/msg-wechat /usr/bin/msg-wechatd
sudo /bin/cp -rf ./src/wechat-messager.sh /usr/bin/
/bin/cp -rf ./src/wechat-messager.sh.desktop ~/.config/autostart/

echo "完成"
