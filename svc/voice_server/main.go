package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	pb "github.com/TRConley/clubhouse-backend-clone/gen/go"
	"github.com/TRConley/clubhouse-backend-clone/pkg/grpc"
	"github.com/TRConley/clubhouse-backend-clone/pkg/signals"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func main() {
	// init the configuration
	initConfig()

	// configure the logger
	logger, err := initZap(viper.GetString("level"))
	if err != nil {
		panic(errors.Wrap(err, "unable to init zap"))
	}
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	// load gRPC configuration
	var grpcConfig grpc.Config
	if err := viper.Unmarshal(&grpcConfig); err != nil {
		logger.Panic("unable to unmarshal grpc config", zap.Error(err))
	}

	logger.Info("starting service",
		zap.String("name", viper.GetString("svc")),
		zap.Int("port", grpcConfig.Port),
	)

	// starting the gRPC server
	grpcServer := grpc.NewServer(&grpcConfig, logger)
	stopCh := signals.SetupSignalHandler()
	pb.RegisterVoiceServiceServer(grpcServer, &VoiceService{})
	grpcServer.ListenAndServe(stopCh)
}

func initConfig() {
	viper.SetConfigFile("config.yaml")
	viper.SetDefault("svc", "placeholder")
	viper.SetDefault("level", "info")
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "unable to init config"))
	}
}

func initZap(logLevel string) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return zapConfig.Build()
}
