package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type logLevel int

const (
	debug logLevel = iota
	info
	warn
	fatal
)

func logrusOutputf(l logLevel, format string, args ...interface{}) {
	switch {
	case l == debug:
		logrus.Debugf(format, args)
	case l == info:
		logrus.Infof(format, args)
	case l == warn:
		logrus.Warnf(format, args)
	case l == fatal:
		logrus.Fatalf(format, args)
	}
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func statusCodeColor(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func statusLogLevel(code int) logLevel {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return info
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return debug
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return warn
	default:
		return fatal
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		methodName := c.Request.Method
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		statusColor := statusCodeColor(statusCode)
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		errorMessage := c.Errors.ByType(1).String()
		if raw != "" {
			path = path + "?" + raw
		}

		if Cfg.Debug == true {
			logrus.SetOutput(ioutil.Discard)
			fmtLogStringWithColor := fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
				endTime.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, reset,
				latency,
				clientIP,
				methodColor(methodName), methodName, reset,
				path,
				errorMessage,
			)
			fmt.Print(fmtLogStringWithColor)
		}
		fmtLogStringWithoutColor := fmt.Sprintf("[GIN] | Result: %3d | Latency: %13v | From: %15s | Method: %-7s  URL: %#v %s",
			statusCode,
			latency,
			clientIP,
			methodName,
			path,
			errorMessage)

		logrusOutputf(statusLogLevel(statusCode), "%s", fmtLogStringWithoutColor)
		if Cfg.Debug == true {
			logrus.SetOutput(os.Stdout)
		}
	}
}

func InitLogger() {
	logrus.Infoln("init logger")
	var logLevel = logrus.InfoLevel
	if Cfg.Debug {
		logLevel = logrus.DebugLevel
	}
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/console.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,

		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	logrus.SetLevel(logLevel)
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		TimestampFormat:        time.RFC3339,
	})
	logrus.AddHook(rotateFileHook)
}
