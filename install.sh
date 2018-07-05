#!/bin/bash

version="0.3.0"
platform="linux"
os_uname=`uname`
if [[ $os_uname == 'Linux' ]]; then
    platform='linux'
elif [[ $os_uname == 'Darwin' ]]; then
    platform='darwin'
fi

if [ ! -d "/etc/systemd" ]; then
    echo -e "安装失败！\n您的操作系统不支持systemd系统管理服务，本脚本暂不支持您的系统意见安装。\n请更换您的操作操作系统为Ubuntu16.04 LTS, Centos 7等系统，再执行该命令"
    exit 1
fi

goshadowsocks=`systemctl status go-shadowsocks | grep Active | awk '{print $3}' | cut -d "(" -f2 | cut -d ")" -f1`

curl -Lo go-shadowsocks-server.tar.gz https://github.com/sedgwickz/go-shadowsocks/releases/download/$version/ss-server-$platform-amd64.tar.gz \
&& tar xzf go-shadowsocks-server.tar.gz && sudo mv ss-server-$platform-amd64 /usr/local/bin/ssserver && rm go-shadowsocks-server.tar.gz \
&& curl -Lo ss-config.json https://raw.githubusercontent.com/sedgwickz/go-shadowsocks/master/sample-config.json \ 
&& mkdir -p ~/.shadowsocks && mv ss-config.json ~/.shadowsocks/config.json \
&& curl -Lo https://github.com/sedgwickz/go-shadowsocks/raw/master/script/go-shadowsocks.service \
&& mv go-shadowsocks.service /ect/systemd/system \
&& systemctl daemon-reload \

echo "🍻已安装成功！"

if [ $goshadowsocks ==  "running" ]; then
    systemctl restart go-shadowsocks
    echo "已重启go-shadowsocks服务"
else
    systemctl start go-shadowsocks
    echo "go-shadowsocks服务已启动"
fi

echo -e "配置文件位于 ~/.shadowsocks/config.json \n建议您及时更改端口和密码 \n更改成功后使用 systemctl restart go-shadowsocks 重启服务"