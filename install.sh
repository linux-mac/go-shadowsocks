#!/bin/bash

version="0.3.0"
platform="linux"
os_uname=`uname`
if [[ $os_uname == 'Linux' ]]; then
    platform='linux'
elif [[ $os_uname == 'Darwin' ]]; then
    platform='darwin'
fi

if [ ! -d "/ect/systemd" ]; then
    echo "😞您的操作系统不支持systemd系统管理服务，本脚本暂不支持您的系统意见安装。请更换您的操作操作系统为Ubuntu16.04LTS, Centos7等系统，再执行该命令"
    exit 1
fi

curl -Lo go-shadowsocks-server.tar.gz https://github.com/sedgwickz/go-shadowsocks/releases/download/$version/ss-server-$platform-amd64.tar.gz
tar xzf go-shadowsocks-server.tar.gz && sudo mv ss-server-$platform-amd64 /usr/local/bin/ssserver && rm go-shadowsocks-server.tar.gz
mkdir -p ~/.shadowsocks && curl -Lo ss-config.json https://raw.githubusercontent.com/sedgwickz/go-shadowsocks/master/sample-config.json && mv ss-config.json ~/.shadowsocks/config.json
curl -Lo https://github.com/sedgwickz/go-shadowsocks/raw/master/script/go-shadowsocks.service && mv go-shadowsocks.service /ect/systemd/system
systemctl daemon-reload && systemctl start go-shadowsocks
echo "🍻安装成功，配置文件位于 ~/.shadowsocks/config.json，建议您及时更改端口和密码。更改成功后使用 systemctl restart go-shadowsocks 重启服务"