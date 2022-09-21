#!/bin/bash

# WeChat.exe
# msg-wechatd

while true
do
    wechatnum=`ps -ef|grep "WeChat.exe"|grep -v grep|wc -l`
    procnum=`ps -ef|grep "msg-wechatd"|grep -v grep|wc -l`
    #echo "WeChat.exe=$wechatnum"
    #echo "msg-wechatd=$procnum"
    if [ $wechatnum -ne 0 ]; then
        #echo "wechat exist"
        if [ $procnum -eq 0 ]; then
            echo "restart msg-wechatd"
            /usr/bin/msg-wechatd &
        fi
    else
        #echo "wechat not exist"
        if [ $procnum -ne 0 ]; then
            echo "stop msg-wechatd"
            kill -9 $(pidof msg-wechatd)
        fi
    fi

    sleep 10
done
