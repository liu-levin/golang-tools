package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/liu-levin/golang-tools/pkg/cleaner"
	"github.com/liu-levin/golang-tools/pkg/logger"

	"github.com/jinzhu/configor"
	"github.com/robfig/cron/v3"
)

type config struct {
	ARCH         string `yaml:"ARCH"`
	FILE_PATH    string `yaml:"FILE_PATH"`
	LOG_FILE     string `yaml:"LOG_FILE"`
	CRON_TIME    string `yaml:"CORN_TIME""`
	EXPIRED_DAYS int    `yaml:"EXPIRED_DAYS"`
}

func main() {
	conf := getConfig()
	logger.InitLogger(conf.LOG_FILE)
	runner := cron.New()
	clean := cleaner.NewCleaner(conf.FILE_PATH)
	_, err := runner.AddJob(conf.CRON_TIME, clean)
	if err != nil {
		panic(err)
	}
	runner.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

	cronCTX := runner.Stop()
	select {
	case <-cronCTX.Done():
		fmt.Println("exit done")
	case <-time.After(5 * time.Minute):
		fmt.Println("exit timeout")
	}

}

func getConfig() config {
	var conf config
	err := configor.Load(&conf, "./clean.yaml")
	if err != nil {
		panic(fmt.Errorf("load config failed: %w", err))
	}
	fmt.Printf("config is %+v\n", conf)

	if conf.ARCH == "" {
		panic(errors.New("no arch"))
	}

	if conf.LOG_FILE == "" {
		panic(errors.New("no log file"))
	}

	if conf.CRON_TIME == "" {
		panic(errors.New("no cron time"))
	}

	return conf
}
