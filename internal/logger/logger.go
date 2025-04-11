package logger

import log "github.com/sirupsen/logrus"

func InitLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}
