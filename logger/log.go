package logger

import (
	"github.com/BAN1ce/skyTree/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var Logger *SkyLogger

type SkyLogger struct {
	*zerolog.Logger
}

func Load() {
	var (
		level, err = zerolog.ParseLevel(config.GetConfig().GetLog().Level)
	)
	if err != nil {
		log.Fatal().Err(err)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(level)

	Logger = &SkyLogger{
		Logger: &logger,
	}
	return

	//level, err := zap.ParseAtomicLevel(config.GetConfig().Log.Level)
	//if err != nil {
	//	logger.Fatal("init logger error")
	//}
	//
	//var (
	//	zapConfig = zap.Config{
	//		Level:       level,
	//		Development: true,
	//		Sampling: &zap.SamplingConfig{
	//			Initial:    100,
	//			Thereafter: 100,
	//		},
	//		Encoding: "console",
	//		EncoderConfig: zapcore.EncoderConfig{
	//			TimeKey:        "time",
	//			LevelKey:       "level",
	//			NameKey:        "logger",
	//			CallerKey:      "caller",
	//			FunctionKey:    zapcore.OmitKey,
	//			MessageKey:     "msg",
	//			StacktraceKey:  "stacktrace",
	//			LineEnding:     zapcore.DefaultLineEnding,
	//			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
	//			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
	//			EncodeDuration: zapcore.SecondsDurationEncoder,
	//			EncodeCaller:   zapcore.ShortCallerEncoder,
	//		},
	//		OutputPaths:      []string{"stderr"},
	//		ErrorOutputPaths: []string{"stderr"},
	//	}
	//	l *zap.Logger
	//)
	//l, err = zapConfig.Build()
	//if err != nil {
	//	logger.Fatal("init logger error")
	//}
	//Logger = &SkyLogger{
	//	Logger: l,
	//}
}
