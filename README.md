# cftunnel_wrap

把 [cloudflared](https://github.com/cloudflare/cloudflared) 作为 Go 库内嵌，编译成单个二进制。只支持一种用法：通过 **token** 运行远程托管（remotely-managed）的 Cloudflare Tunnel。

转发规则（public hostname → 本地服务）都在 Cloudflare Zero Trust 面板里配置，本程序只负责用 token 建立隧道。

## 配置

参考 [`config.example.yaml`](config.example.yaml)：

```yaml
# Zero Trust 面板 -> Networks -> Tunnels -> 你的隧道 -> 安装命令里 --token 后面那串
token: eyJhIjoi...
# auto | quic | http2，省略时为 auto；UDP 7844 不通时用 http2
protocol: auto
```

| 字段 | 必填 | 说明 |
|---|---|---|
| `token` | 是 | 隧道 token，即 base64 编码的 `{"a","t","s"}` JSON，启动前会先校验格式 |
| `protocol` | 否 | 连接边缘节点用的协议：`auto`（默认）、`quic` 或 `http2` |

> token 等同于隧道凭据，请不要提交到仓库。`.gitignore` 已经忽略了 `config.yaml`。

## 构建与运行

需要 Go 1.26 或更高版本（go.mod 要求的版本；本地版本较低时，`GOTOOLCHAIN=auto` 会自动下载）。

```bash
go build -trimpath -ldflags "-s -w" -o cftunnel .
./cftunnel                       # 读取当前目录下的 config.yaml
./cftunnel -c /path/to/config.yaml
```

| 选项 | 说明 |
|---|---|
| `-c`, `--config <path>` | 配置文件路径，默认 `./config.yaml` |
| `-h`, `--help` | 显示用法 |
| `-v`, `--version` | 显示版本 |

运行时等价于：

```
cloudflared tunnel --no-autoupdate --protocol <protocol> run --token <token>
```

token 只在进程内部传递，不会出现在 `ps` 的命令行里。按 Ctrl-C 或发送 SIGTERM 会触发优雅退出。

交叉编译示例：

```bash
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o cftunnel.exe .
GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o cftunnel-mac .
GOOS=android GOARCH=arm64 go build -trimpath -ldflags "-s -w" -o libcftunnel.so .
```

## 网络要求

- 能出站访问 `region1.v2.argotunnel.com` / `region2.v2.argotunnel.com` 的 **7844** 端口（QUIC 走 UDP，HTTP/2 走 TCP）。
- 不需要入站端口，也不需要公网 IP。

## 实现说明

- 调用 cloudflared 的 `tunnel.Init()` 和 `tunnel.Commands()`，在进程内执行 `tunnel run` 子命令。连接、重连和远程配置下发都沿用官方逻辑。
- `go.mod` 中的 4 条 `replace` 是从 cloudflared 的 go.mod 原样复制的（它依赖 urfave/cli 和 quic-go 的 fork 版本）。Go 的 replace 不会传递给下游模块，所以**不能删除**。升级 cloudflared 时，需要对照新版本的 go.mod 同步更新。
- `tunnel.Init` 等是 cloudflared 的内部 API，不保证跨版本稳定。当前锁定在 cloudflared `2026.9.3`，升级后需要重新编译验证。

## 许可证

本项目以 [Apache License 2.0](LICENSE) 发布。

本项目静态链接了 cloudflared（Copyright Cloudflare, Inc.，Apache License 2.0），未做修改。编译出的二进制还包含其他第三方 Go 模块，均为 Apache-2.0、MIT、BSD-2-Clause 或 BSD-3-Clause 等宽松许可证：

- [`NOTICE`](NOTICE)：署名声明，以及 Apache-2.0 依赖要求保留的 NOTICE 内容；
- [`THIRD_PARTY_LICENSES.md`](THIRD_PARTY_LICENSES.md)：所有链接模块的版本和许可证清单。

分发二进制时，请一并附上 `LICENSE`、`NOTICE` 和 `THIRD_PARTY_LICENSES.md`。

升级 cloudflared 或调整依赖后，在仓库根目录重新生成这两个文件：

```bash
go run ./scripts/genlicenses
```

它会编译 cftunnel，从二进制的 build info 里读出实际链接的模块（已应用 `replace`），再到 `GOMODCACHE` 中查找各模块的 LICENSE 和 NOTICE。如果有模块的许可证无法识别或缺失，命令会以非零状态退出，这时需要人工确认。如果出现 MPL、GPL 等非宽松许可证，也要人工评估。

本项目与 Cloudflare, Inc. 无关联，"Cloudflare" 是 Cloudflare, Inc. 的商标。
