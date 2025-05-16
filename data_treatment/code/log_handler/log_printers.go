package log_handler

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func Success(format string, args ...any) {
	msg := strings.TrimSpace(fmt.Sprintf(format, args...))
	fmt.Println("\033[32m[SUCCESS]\033[0m ", msg)
	Log("[SUCCESS] " + strings.TrimSpace(msg))
}

func Error(format string, args ...any) {
	msg := strings.TrimSpace(fmt.Sprintf(format, args...))
	fmt.Println("\033[31m[ERROR]\033[0m ", msg)
	Log("[ERROR] " + msg)
}

func Fatal(format string, args ...any) {
	msg := strings.TrimSpace(fmt.Sprintf(format, args...))
	fmt.Println("\033[31m[ERROR]\033[0m ", msg)
	Log("[ERROR] " + msg)
	time.Sleep(20 * time.Millisecond)
	FlushAll()
	os.Exit(1)
}
