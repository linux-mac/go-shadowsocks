# go-shadowsocks

Inspired by [shadowsocks-go](https://github.com/shadowsocks/shadowsocks-go)

go-shadowsocks是shadowsocks协议go语言实现，仅供个人学习参考。

目前支持以下协议：

1. aes-256-cfb
2. aes-256-gcm (ToDo)
3. chacha20-ietf-poly1305（ToDo）

## Getting Started

go-shadowsocks一键服务端安装脚本目前仅适合支持systemd服务管理的 *nux 系统，请选择相应操作系统，如Ubuntu 16.04 LTS, Centos7等

若不使用一键脚本，可使用supervisor等进程管理程序进行go-shadowsocks进程管理。

go-shadowsocks支持服务器多端口不同加密算法，具体使用参考配置部分。

### Installing

服务端执行下面脚本直接安装，密码及端口设置参考下方配置文件：

`curl -sSL https://git.io/shadowsocks | bash`

### Config

安装完成后默认配置文件在 `~/.shadowsocks/.config.json`，字段说明请参考simple-config.json文件

您可根据需要使用以下参数：

```
-d 打印调试信息
-c path 指定配置文件路径，默认路径为 ~/.shadowsocks/.config.json
-v 打印当前使用版本号
```

## License

Apache 2.0