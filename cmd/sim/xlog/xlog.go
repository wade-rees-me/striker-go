package xlog

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/wade-rees-me/striker-go/cmd/sim/constants"
)

const (
	SyslogAddress  = "192.168.0.27"
	SyslogPort     = 10514
	SyslogMsgMax   = 1024
	SyslogFacility = 1 << 3 // USER facility
	SyslogEmerg    = 0
	SyslogAlert    = 1
	SyslogCrit     = 2
	SyslogErr      = 3
	SyslogWarning  = 4
	SyslogNotice   = 5
	SyslogInfo     = 6
	SyslogDebug    = 7
)

var (
	conn     *net.UDPConn
	server   *net.UDPAddr
	initOnce sync.Once
	hostname string
)

// InitSyslog sets up the UDP connection to the remote syslog server.
func InitSyslog(remoteHost string, port int) error {
	var err error
	initOnce.Do(func() {
		server, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", remoteHost, port))
		if err != nil {
			return
		}
		conn, err = net.DialUDP("udp", nil, server)
		if err != nil {
			return
		}
		h, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		} else {
			hostname = h
		}
	})
	return err
}

// CloseSyslog closes the UDP socket.
func CloseSyslog() {
	if conn != nil {
		conn.Close()
		conn = nil
	}
}

// xlogSyslog builds and sends the log packet.
func xlogSyslog(severity int, message string) {
	if conn == nil {
		return
	}
	priority := SyslogFacility + severity
	packet := fmt.Sprintf("<%d>%s: [version=%s] [PID=%d] | %s",
		priority, constants.StrikerWhoAmI, constants.StrikerVersion, os.Getpid(), message)

	if len(packet) > SyslogMsgMax {
		packet = packet[:SyslogMsgMax]
	}

	_, _ = conn.Write([]byte(packet))
}

// LogInfo logs an informational message.
func LogInfo(format string, args ...any) {
	xlogSyslog(SyslogInfo, fmt.Sprintf(format, args...))
}

// LogError logs an error message.
func LogError(format string, args ...any) {
	xlogSyslog(SyslogErr, fmt.Sprintf(format, args...))
}

// LogFatal logs a critical (fatal) message.
func LogFatal(format string, args ...any) {
	xlogSyslog(SyslogCrit, fmt.Sprintf(format, args...))
}
