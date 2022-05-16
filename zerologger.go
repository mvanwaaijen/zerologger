package zerologger

import (
	"fmt"
	"io"
	"os"

	"github.com/mvanwaaijen/execpath"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	ISO8601TimeString string = "2006-01-02T15:04:05.000-07:00"
	LocalTimeString   string = "15:04:05.000"
	DefaultSizeMB     int    = 5
	DefaultAgeDays    int    = 365
	DefaultBackups    int    = 0
)

type Config struct {
	FileName                string
	LogToFile               bool
	LogToConsole            bool
	MaxSizeMB               int
	MaxAgeDays              int
	MaxBackups              int
	ShowCaller              bool
	FileTimeFormatString    string
	ConsoleTimeFormatString string
}

func NewDefaultConfig() *Config {
	ep, err := execpath.Get()
	if err != nil {
		panic(err)
	}
	return &Config{
		FileName:                fmt.Sprintf("%s.log", ep),
		LogToFile:               true,
		LogToConsole:            true,
		MaxSizeMB:               DefaultSizeMB,
		MaxAgeDays:              DefaultAgeDays,
		MaxBackups:              DefaultBackups,
		ShowCaller:              true,
		FileTimeFormatString:    ISO8601TimeString,
		ConsoleTimeFormatString: LocalTimeString,
	}
}

func NewDefault() zerolog.Logger {
	return New(NewDefaultConfig())
}

func New(cfg *Config) zerolog.Logger {
	writers := make([]io.Writer, 0)
	if cfg.LogToFile {
		writers = append(writers, &lumberjack.Logger{
			Filename:   cfg.FileName,
			MaxSize:    cfg.MaxSizeMB,
			MaxAge:     cfg.MaxAgeDays,
			MaxBackups: cfg.MaxBackups,
		})
	}
	if cfg.LogToConsole {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: cfg.ConsoleTimeFormatString})
	}
	zerolog.TimeFieldFormat = cfg.FileTimeFormatString
	if cfg.ShowCaller {
		log.Logger = log.With().Caller().Logger().Output(io.MultiWriter(writers...))
	} else {
		log.Logger = log.Output(io.MultiWriter(writers...))
	}
	return log.Logger
}
