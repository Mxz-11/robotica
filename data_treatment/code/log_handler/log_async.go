package log_handler

import (
	"data_treatment/config_handler"
	"data_treatment/shared"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type AsyncLogger struct {
	channel  chan string
	wait     sync.WaitGroup
	buffer   []string
	ticker   *time.Ticker
	quit     chan struct{}
	log_file *os.File
}

var logger *AsyncLogger

func prepareLogFile(log_folder string, name string, ext string) (*os.File, error) {
	err := os.MkdirAll(log_folder, 0755)
	if err != nil {
		return nil, err
	}

	for i := 0; ; i++ {
		var filename string
		if i == 0 {
			filename = fmt.Sprintf("%s%s", name, ext)
		} else {
			filename = fmt.Sprintf("%s_%d%s", name, i, ext)
		}
		full_path := filepath.Join(log_folder, filename)
		if _, err := os.Stat(full_path); os.IsNotExist(err) {
			return os.Create(full_path)
		}
	}
}

func CreateAsyncLogger(log_folder string, name string, ext string) error {
	log_file, err := prepareLogFile(log_folder, name, ext)
	if err != nil {
		return err
	}
	buffer_raw, err := config_handler.GetData(shared.Consts, "LOG_MAX_BUFFER_SIZE", config_handler.TYPE_INT)
	if err != nil {
		log.Fatal("\033[31m[ERROR]\033[0m LOG_MAX_BUFFER_SIZE must be integer")
	}
	buffer_s := buffer_raw.(int)

	interval_raw, err := config_handler.GetData(shared.Consts, "LOG_TIME_INTERVAL", config_handler.TYPE_TIME)
	if err != nil {
		log.Fatalf("\033[31m[ERROR]\033[0m LOG_TIME_INTERVAL must be time.Duration")
	}
	interval := interval_raw.(time.Duration)
	logger = &AsyncLogger{
		channel:  make(chan string, 100),
		buffer:   make([]string, 0, buffer_s),
		ticker:   time.NewTicker(interval),
		quit:     make(chan struct{}),
		log_file: log_file,
	}
	logger.wait.Add(1)
	go runLogger()
	return nil
}

func Log(msg string) {
	if logger != nil {
		logger.channel <- msg
	}
}

func runLogger() {
	defer logger.wait.Done()
	for {
		select {
		case msg := <-logger.channel:
			logger.buffer = append(logger.buffer, msg)
			if len(logger.buffer) >= 10 {
				flushLogger()
			}
		case <-logger.ticker.C:
			if len(logger.buffer) > 0 {
				flushLogger()
			}
		case <-logger.quit:
			flushLogger()
			logger.log_file.Close()
			return
		}
	}
}

func flushLogger() {
	if logger == nil || logger.log_file == nil {
		return
	}
	for _, msg := range logger.buffer {
		fmt.Fprintln(logger.log_file, msg)
	}
	logger.buffer = logger.buffer[:0]
}

func FlushAll() {
	if logger == nil {
		return
	}

flush_loop:
	for {
		select {
		case msg := <-logger.channel:
			logger.buffer = append(logger.buffer, msg)
		default:
			break flush_loop
		}
	}
	flushLogger()
}

func Close() {
	if logger != nil {
		close(logger.quit)
		logger.wait.Wait()
	}
}
