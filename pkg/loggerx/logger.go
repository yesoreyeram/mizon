// shared logging utilities
package loggerx

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Level int

const (
	LevelDebug Level = iota + 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var currentLevel = LevelInfo
var jsonMode = false

// Setup configures the global logger format and level via LOG_LEVEL.
// LOG_LEVEL can be one of: DEBUG, INFO, WARN, ERROR, FATAL (default: INFO)
func Setup() {
	// Determine level
	if lvlStr := strings.TrimSpace(os.Getenv("LOG_LEVEL")); lvlStr != "" {
		SetLevel(ParseLevel(lvlStr))
	}
	// Determine format
	fmtStr := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	if fmtStr == "json" || EnvBool("LOG_JSON", false) {
		jsonMode = true
		log.SetFlags(0)
	} else {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}
}

func SetLevel(l Level) { currentLevel = l }

func ParseLevel(s string) Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return LevelDebug
	case "INFO", "":
		return LevelInfo
	case "WARN", "WARNING":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	default:
		return LevelInfo
	}
}

func levelName(l Level) string {
	switch l {
	case LevelDebug:
		return "🔍 DEBUG"
	case LevelInfo:
		return "▶︎ INFO"
	case LevelWarn:
		return "⚠️ WARN"
	case LevelError:
		return "🚨 ERROR"
	case LevelFatal:
		return "☠️FATAL"
	default:
		return "ℹ️ INFO"
	}
}

type logEntry struct {
	Ts     string                 `json:"ts"`
	Level  string                 `json:"level"`
	Msg    string                 `json:"msg"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

func logAt(l Level, msg string, fields map[string]interface{}) {
	if l < currentLevel && l != LevelFatal {
		return
	}
	if jsonMode {
		e := logEntry{Ts: time.Now().UTC().Format(time.RFC3339Nano), Level: levelName(l), Msg: msg}
		if len(fields) > 0 {
			e.Fields = fields
		}
		b, _ := json.Marshal(e)
		log.Print(string(b))
	} else {
		if len(fields) > 0 {
			log.Printf("[%s] %s | %+v", levelName(l), msg, fields)
		} else {
			log.Printf("[%s] %s", levelName(l), msg)
		}
	}
	if l == LevelFatal {
		os.Exit(1)
	}
}

// Message helpers
func Debug(msg string) { logAt(LevelDebug, msg, nil) }
func Info(msg string)  { logAt(LevelInfo, msg, nil) }
func Warn(msg string)  { logAt(LevelWarn, msg, nil) }
func Error(msg string) { logAt(LevelError, msg, nil) }
func Fatal(msg string) { logAt(LevelFatal, msg, nil) }

// Formatted helpers
func Debugf(format string, a ...any) { logAt(LevelDebug, fmt.Sprintf(format, a...), nil) }
func Infof(format string, a ...any)  { logAt(LevelInfo, fmt.Sprintf(format, a...), nil) }
func Warnf(format string, a ...any)  { logAt(LevelWarn, fmt.Sprintf(format, a...), nil) }
func Errorf(format string, a ...any) { logAt(LevelError, fmt.Sprintf(format, a...), nil) }
func Fatalf(format string, a ...any) { logAt(LevelFatal, fmt.Sprintf(format, a...), nil) }

// Structured helpers
func Debugw(msg string, fields map[string]interface{}) { logAt(LevelDebug, msg, fields) }
func Infow(msg string, fields map[string]interface{})  { logAt(LevelInfo, msg, fields) }
func Warnw(msg string, fields map[string]interface{})  { logAt(LevelWarn, msg, fields) }
func Errorw(msg string, fields map[string]interface{}) { logAt(LevelError, msg, fields) }
func Fatalw(msg string, fields map[string]interface{}) { logAt(LevelFatal, msg, fields) }
