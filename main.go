// cftunnel 内嵌 cloudflared，仅支持通过 token 运行远程托管隧道。
//
// 用法: cftunnel <config.yaml>
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/cloudflare/cloudflared/cmd/cloudflared/cliutil"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/tunnel"
	"github.com/cloudflare/cloudflared/tracing"
)

const version = "0.1.0"

type Config struct {
	Token    string `yaml:"token"`
	Protocol string `yaml:"protocol"`
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <config.yaml>\n", os.Args[0])
		os.Exit(2)
	}
	cfg, err := loadConfig(os.Args[1])
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
	// 等价于: cloudflared tunnel --no-autoupdate --protocol <p> run --token <t>
	// 参数只在进程内传递，不会出现在 ps 里
	args := []string{
		"cloudflared", "tunnel",
		"--no-autoupdate",
		"--protocol", cfg.Protocol,
		"run", "--token", cfg.Token,
	}
	if err := app.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
