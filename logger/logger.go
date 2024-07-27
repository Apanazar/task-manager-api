package logger

import "log"

func Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func Debug(msg string, args ...interface{}) {
	log.Printf("[DEBUG] "+msg, args...)
}
