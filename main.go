package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wordAvailable/constant"
	"wordAvailable/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	mw := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "./log/log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	log.SetOutput(mw)
	constant.ReadConfig(".env")
}

func main() {
	gin.SetMode(viper.GetString("RUN_MODE"))
	port := viper.GetString("PORT")
	routesInit := routes.InitRouter()
	server := &http.Server{
		Addr:           port,
		Handler:        routesInit,
		ReadTimeout:    time.Duration(viper.GetInt("READ_TIMEOUT")) * time.Second,
		WriteTimeout:   time.Duration(viper.GetInt("WRITE_TIMEOUT")) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGILL, syscall.SIGFPE)
	go func() {
		s := <-c
		log.Printf("[info] Signal Recived: %s", s)
		os.Exit(0)
	}()

	log.Printf("[info] start http server listening %s", port)
	server.ListenAndServe()
}
