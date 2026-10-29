// cftunnel 内嵌 cloudflared，仅支持通过 token 运行远程托管隧道。
//
// 用法: cftunnel [-c config.yaml]
package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/cloudflare/cloudflared/cmd/cloudflared/cliutil"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/tunnel"
	"github.com/cloudflare/cloudflared/tracing"
)

const version = "0.1.0"

type Config struct {
	Token        string `yaml:"token"`
	Protocol     string `yaml:"protocol"`
	LogLevel     string `yaml:"loglevel"`
	LogFile      string `yaml:"logfile"`
	LogDirectory string `yaml:"log-directory"`
	PidFile      string `yaml:"pidfile"`
	Retries      *int   `yaml:"retries"` // nil 表示沿用 cloudflared 默认值 5
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("%s: token is required", path)
	}
	if err := checkToken(cfg.Token); err != nil {
		return nil, fmt.Errorf("%s: invalid token: %w", path, err)
	}
	switch cfg.Protocol {
	case "":
		cfg.Protocol = "auto"
	case "auto", "quic", "http2":
	default:
		return nil, fmt.Errorf("%s: protocol must be auto, quic or http2, got %q", path, cfg.Protocol)
	}
	switch cfg.LogLevel {
	case "":
		cfg.LogLevel = "info"
	case "debug", "info", "warn", "error", "fatal":
	default:
		return nil, fmt.Errorf("%s: loglevel must be debug, info, warn, error or fatal, got %q", path, cfg.LogLevel)
	}
	if cfg.Retries != nil && *cfg.Retries < 0 {
		return nil, fmt.Errorf("%s: retries must be >= 0, got %d", path, *cfg.Retries)
	}
	return &cfg, nil
}

// checkToken 提前校验 token 格式（base64 编码的 {"a","t","s"}），让错误更早、更清楚。
func checkToken(token string) error {
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return fmt.Errorf("not base64: %w", err)
	}
	var t struct {
		AccountTag string `json:"a"`
		TunnelID   string `json:"t"`
		Secret     string `json:"s"`
	}
	if err := json.Unmarshal(raw, &t); err != nil {
		return fmt.Errorf("not JSON: %w", err)
	}
	if t.AccountTag == "" || t.TunnelID == "" || t.Secret == "" {
		return fmt.Errorf("missing account tag, tunnel id or secret")
	}
	return nil
}

const usage = `cftunnel - 用 token 运行 Cloudflare Tunnel（内嵌 cloudflared）

用法:
  cftunnel [-c <path>]

选项:
  -c, --config <path>  配置文件路径（默认 ./config.yaml）
  -h, --help           显示本帮助
  -v, --version        显示版本

配置文件 (YAML):
  token: eyJhIjoi...   # 必填，Zero Trust 面板中隧道的 token
  protocol: auto       # 可选，auto | quic | http2，默认 auto
  loglevel: info       # 可选，debug | info | warn | error | fatal，默认 info
  logfile: <path>      # 可选，日志额外写入该文件，默认只输出到终端
  log-directory: <dir> # 可选，日志额外写入该目录下的 cloudflared.log（按 1MB 滚动，保留 5 份）；
                       #       与 logfile 同时设置时以 logfile 为准
  pidfile: <path>      # 可选，隧道连接成功后把 PID 写入该文件（退出时不删除）
  retries: 5           # 可选，连接/协议出错时的最大重试次数，默认 5
`

// parseFlags 解析命令行，返回配置文件路径。
func parseFlags(args []string) (string, error) {
	// 错误和用法统一由 main 打印
	fs := flag.NewFlagSet("cftunnel", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var configPath string
	var showVersion bool
	fs.StringVar(&configPath, "c", "config.yaml", "")
	fs.StringVar(&configPath, "config", "config.yaml", "")
	fs.BoolVar(&showVersion, "v", false, "")
	fs.BoolVar(&showVersion, "version", false, "")
	if err := fs.Parse(args); err != nil {
		return "", err
	}
	if fs.NArg() > 0 {
		return "", fmt.Errorf("unexpected argument: %s", fs.Arg(0))
	}
	if showVersion {
		fmt.Println("cftunnel", version)
		os.Exit(0)
	}
	return configPath, nil
}

func main() {
	configPath, err := parseFlags(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		fmt.Print(usage)
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n%s", err, usage)
		os.Exit(2)
	}
	cfg, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 与官方 main 保持一致
	os.Setenv("QUIC_GO_DISABLE_ECN", "1")

	// 关闭该 channel 可触发优雅退出；SIGINT/SIGTERM 由 cloudflared 内部处理
	graceShutdownC := make(chan struct{})
	tunnel.Init(cliutil.GetBuildInfo("embedded", version), graceShutdownC)
	tracing.Init(version)

	app := &cli.App{
		Name:     "cloudflared",
		Version:  version,
		Commands: tunnel.Commands(),
	}
	// 等价于: cloudflared tunnel --no-autoupdate --protocol <p> --loglevel <l> [--logfile <f>] [--log-directory <d>] [--pidfile <p>] [--retries <n>] run --token <t>
	// 参数只在进程内传递，不会出现在 ps 里
	args := []string{
		"cloudflared", "tunnel",
		"--no-autoupdate",
		"--protocol", cfg.Protocol,
		"--loglevel", cfg.LogLevel,
	}
	if cfg.LogFile != "" {
		args = append(args, "--logfile", cfg.LogFile)
	}
	// 两者都指定时 cloudflared 只用 --logfile
	if cfg.LogDirectory != "" {
		args = append(args, "--log-directory", cfg.LogDirectory)
	}
	if cfg.PidFile != "" {
		args = append(args, "--pidfile", cfg.PidFile)
	}
	if cfg.Retries != nil {
		args = append(args, "--retries", strconv.Itoa(*cfg.Retries))
	}
	args = append(args, "run", "--token", cfg.Token)
	if err := app.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
