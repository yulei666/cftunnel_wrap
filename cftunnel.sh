#!/bin/sh
# cftunnel 启动脚本：后台运行 cftunnel，支持 start | stop | restart | status。
#
# 默认在脚本所在目录查找二进制和配置，可用环境变量覆盖：
#   CFTUNNEL_BIN      二进制路径，默认 <脚本目录>/cftunnel
#   CFTUNNEL_CONFIG   配置文件路径，默认 <脚本目录>/config.yaml
#   CFTUNNEL_PIDFILE  本脚本使用的 PID 文件，默认 <脚本目录>/cftunnel.pid
#   CFTUNNEL_LOG      标准输出/错误重定向到的文件，默认 <脚本目录>/cftunnel.log
#   CFTUNNEL_STOP_TIMEOUT  stop 时等待优雅退出的秒数，超时后 SIGKILL，默认 35
#                          （cloudflared 收到 SIGTERM 后最多等待 30 秒处理进行中的请求）

DIR=$(cd "$(dirname "$0")" && pwd)
BIN=${CFTUNNEL_BIN:-$DIR/cftunnel}
CONFIG=${CFTUNNEL_CONFIG:-$DIR/config.yaml}
PIDFILE=${CFTUNNEL_PIDFILE:-$DIR/cftunnel.pid}
LOG=${CFTUNNEL_LOG:-$DIR/cftunnel.log}
STOP_TIMEOUT=${CFTUNNEL_STOP_TIMEOUT:-35}

# 注意：这里的 PID 文件由脚本在启动时写入，与配置里的 pidfile 无关
# （后者要等隧道连上才写，退出时也不删除，不适合用来管理进程）。

# 输出正在运行的 PID；未运行时返回非 0，并清理过期的 PID 文件
running_pid() {
	[ -f "$PIDFILE" ] || return 1
	pid=$(cat "$PIDFILE" 2>/dev/null)
	case $pid in
	'' | *[!0-9]*) rm -f "$PIDFILE"; return 1 ;;
	esac
	# PID 可能已被其他进程复用，核对进程名（comm 最长 15 个字符）
	name=$(basename "$BIN" | cut -c1-15)
	if kill -0 "$pid" 2>/dev/null && [ "$(ps -p "$pid" -o comm= 2>/dev/null | xargs basename 2>/dev/null)" = "$name" ]; then
		echo "$pid"
		return 0
	fi
	rm -f "$PIDFILE"
	return 1
}

start() {
	if pid=$(running_pid); then
		echo "cftunnel is already running (pid $pid)"
		return 0
	fi
	if [ ! -x "$BIN" ]; then
		echo "binary not found or not executable: $BIN" >&2
		return 1
	fi
	if [ ! -f "$CONFIG" ]; then
		echo "config not found: $CONFIG" >&2
		return 1
	fi

	nohup "$BIN" -c "$CONFIG" >>"$LOG" 2>&1 </dev/null &
	pid=$!
	echo "$pid" >"$PIDFILE"

	# 配置错误会立即退出，稍等片刻再确认进程还在
	sleep 1
	if ! kill -0 "$pid" 2>/dev/null; then
		rm -f "$PIDFILE"
		echo "cftunnel failed to start, last lines of $LOG:" >&2
		tail -n 20 "$LOG" >&2
		return 1
	fi
	echo "cftunnel started (pid $pid), log: $LOG"
}

stop() {
	if ! pid=$(running_pid); then
		echo "cftunnel is not running"
		return 0
	fi
	kill -TERM "$pid"
	printf "stopping cftunnel (pid %s)" "$pid"
	i=0
	while kill -0 "$pid" 2>/dev/null; do
		if [ "$i" -ge "$STOP_TIMEOUT" ]; then
			printf " timeout, sending SIGKILL"
			kill -KILL "$pid" 2>/dev/null
			break
		fi
		printf "."
		sleep 1
		i=$((i + 1))
	done
	echo
	rm -f "$PIDFILE"
	echo "cftunnel stopped"
}

status() {
	if pid=$(running_pid); then
		echo "cftunnel is running (pid $pid)"
	else
		echo "cftunnel is not running"
		return 3
	fi
}

case $1 in
start) start ;;
stop) stop ;;
restart) stop && start ;;
status) status ;;
*)
	echo "usage: $0 {start|stop|restart|status}" >&2
	exit 2
	;;
esac
