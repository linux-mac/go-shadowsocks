version="0.1.0"
platform="linux"
os_uname=`uname`
if [[ $os_uname == 'Linux' ]]; then
    platform='linux'
elif [[ $os_uname == 'Darwin' ]]; then
    platform='Darwin'
fi

curl -Lo go-shadowsocks-server.tar.gz https://github.com/sedgwickz/go-shadowsocks/releases/download/$version/ss-server-$platform-amd64.tar.gz
tar xzf go-shadowsocks-server.tar.gz && sudo mv ss-server /usr/local/bin/ss-server && rm go-shadowsocks-server.tar.gz
mkdir -p ~/.shadowsocks && curl -Lo ss-config.json https://raw.githubusercontent.com/sedgwickz/go-shadowsocks/master/sample-config.json && mv ss-config.json ~/.shadowsocks/config.json