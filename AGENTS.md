# AGENTS.md

给在本仓库工作的 AI 编码代理（以及人类贡献者）的说明。

## 项目概览

cftunnel 是一个 Go 程序，把 cloudflared 作为 Go 模块静态链接进来，只支持用 token 运行远程托管隧道。

| 路径 | 说明 |
|---|---|
| `main.go` | 全部业务逻辑：读取 YAML 配置，拼出 `tunnel run` 参数，调用 cloudflared 的 `tunnel.Init()` / `tunnel.Commands()` |
| `cftunnel.sh` | 后台启动、停止、重启和查看状态的脚本（Linux / macOS） |
| `config.example.yaml` | 配置示例；真正的 `config.yaml` 含 token，已被 `.gitignore` 忽略，**不要提交** |
| `scripts/build-release.sh` | 交叉编译并打包发布产物到 `dist/`，CI 和本地共用 |
| `scripts/genlicenses` | 重新生成 `NOTICE` 和 `THIRD_PARTY_LICENSES.md` |
| `.github/workflows/release.yml` | 编译四个平台；推送 tag 时发布 GitHub Release |

## 常用命令

```bash
go build -trimpath -ldflags "-s -w" -o cftunnel .   # 本地构建
go vet ./...                                         # 静态检查
gofmt -l .                                           # 格式检查（应无输出）
go run ./scripts/genlicenses                         # 依赖变化后重新生成许可证文件
scripts/build-release.sh 0.0.0-dev                   # 本地生成四个平台的发布包
```

## 修改代码的约定

- 文档和代码注释使用中文；commit message 使用英文，标题用祈使句。
- 用户可见的行为（配置项、命令行选项、`cftunnel.sh` 的行为）有变化时，同步更新 `README.md`。
- 依赖有任何变化（`go.mod` / `go.sum`）后，必须运行 `go run ./scripts/genlicenses`，并把 `NOTICE` 和 `THIRD_PARTY_LICENSES.md` 一起提交。
- 发布包必须保持 `CGO_ENABLED=0`，可以在 Linux runner 上交叉编译全部平台。

### 升级 cloudflared

1. cloudflared 的 tag（如 `2026.10.7`）不是 Go 认可的 semver，用 `go get github.com/cloudflare/cloudflared@<tag>` 升级。go.mod 中会记录为 pseudo-version，需要手动把行尾注释改成对应的 tag。
2. 对照新版本 cloudflared 的 `go.mod`，同步 `go` 版本和全部 `replace` 指令（replace 不会传递给下游模块，不能删除）。
3. `go mod tidy`，然后编译并运行 `./cftunnel -h`。`tunnel.Init` 等是 cloudflared 的内部 API，升级后可能编译不过。
4. 运行 `go run ./scripts/genlicenses`，并更新 README 中“当前锁定在 cloudflared …”的版本号。

## 分支与发布流程（必须遵守）

1. **不要直接向 `main` 提交或推送。** 所有改动都在功能分支上完成（如 `<user>/<topic>`），推送后向 `main` 开 PR。
2. PR 的 CI（`release` workflow 的 build 任务）通过后再合并。
3. **只在合并到 `main` 之后打 tag 和发布。** tag 必须打在 `main` 的提交上：

   ```bash
   git switch main && git pull --ff-only
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push origin vX.Y.Z
   ```

   推送 tag 后，workflow 会编译 linux/amd64、linux/arm64、darwin/amd64、darwin/arm64，并创建 GitHub Release。release 任务会检查 tag 是否在 `main` 上，不在则直接失败，不会发布。
4. 版本号遵循 semver，tag 格式为 `vX.Y.Z`；带 `-` 后缀的（如 `v0.3.0-rc.1`）发布为 pre-release。发布前把 `main.go` 中的默认 `version` 改成同一个版本号，随 PR 一起合并。
5. 不要移动或删除已经推送的 tag。发布有问题时，修复后发布新的 patch 版本。
6. 打 tag、推送 tag 和发布都是对外可见的操作，代理执行前需要得到用户明确同意。
