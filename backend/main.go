package main

import (
	"context"
	"db-querier/controller"
	"db-querier/middleware"
	"db-querier/service"
	"db-querier/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	settingLogger()
	watchConfig("config/system.toml")
	watchConfig("config/query.toml")

	binding.Validator = new(utils.ValidatorV9)

	engine := gin.New()
	engine.Use(cors.Default())
	engine.Use(middlewares.Logger)
	engine.Use(gin.RecoveryWithWriter(log.Logger))

	controller.Register(engine)
	service.Initialize()

	server := &http.Server{
		Addr:    viper.GetString("server.host") + ":" + viper.GetString("server.port"),
		Handler: engine,
	}

	go func() {
		log.Info().Msgf("Listening and serving HTTP on %s\n", server.Addr)
		_ = server.ListenAndServe()
	}()
	gracefulShutdown(server)
}

func settingLogger() {
	// initialize logger
	zerolog.DurationFieldUnit = time.Second
	if gin.Mode() == gin.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false}).With().Stack().Caller().Timestamp().Logger()
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Logger = zerolog.New(os.Stdout).With().Stack().Caller().Timestamp().Logger()
	}
}

func watchConfig(name string) {
	v := viper.New()
	v.SetConfigFile(name)
	if err := v.ReadInConfig(); err != nil {
		log.Panic().Err(err).Msg("failed reading config files")
	}
	mergeInViper(v)
	v.OnConfigChange(func(e fsnotify.Event) {
		mergeInViper(v)
		service.Initialize()
		log.Info().Str("config_name", e.Name).Msg("config modified")
	})
	v.WatchConfig()
}

func mergeInViper(v *viper.Viper) {
	if err := viper.MergeConfigMap(v.AllSettings()); err != nil {
		log.Panic().Err(err).Msg("failed reading config files")
	}
}

func gracefulShutdown(server *http.Server) {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)

	<-sig
	log.Info().Msg("Shutdown Server ...")

	service.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed shutdown")
	}
	log.Info().Msg("Server closed")
}
