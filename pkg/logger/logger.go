package logger

import (
   "os"
   "sync"

   "github.com/sirupsen/logrus"
   "github.com/joho/godotenv"
)

var (
   // Customize the logger as per your requirements - logrus is the global logger
   Logger *logrus.Logger
   once   sync.Once
)

func Initialize() {
   once.Do(func() {
      err := godotenv.Load()
      if err != nil {
         logrus.Warn("Error loading .env file")
      }

      appEnv := os.Getenv("APP_ENV")

      Logger = logrus.New()
      Logger.SetOutput(os.Stdout)
      Logger.SetFormatter(&logrus.TextFormatter{
         FullTimestamp: true,
      })
      Logger.SetLevel(getLogLevel(appEnv))

      // Add a default field
      Logger = Logger.WithField("service", "my-service").Logger

      Info("Logger initialized", logrus.Fields{"appEnv": appEnv})
   })
}

func getLogLevel(env string) logrus.Level {
   switch env {
   case "production":
      return logrus.WarnLevel
   case "dev":
      return logrus.DebugLevel
   default:
      logrus.Warn("Unrecognized APP_ENV, defaulting to DebugLevel")
      return logrus.DebugLevel
   }
}

func notInitialized() bool {
   if Logger == nil {
      logrus.Error("Logger not initialized")
      return true
   }
   return false
}

// logrus.Fields is map[string]interface{} for structured logging
func Info(msg string, fields logrus.Fields) {
   if notInitialized() {
      return
   }
   Logger.WithFields(fields).Info(msg)
}

func Debug(msg string, fields logrus.Fields) {
   if notInitialized() {
      return
   }
   Logger.WithFields(fields).Debug(msg)
}

func Warn(msg string, fields logrus.Fields) {
   if notInitialized() {
      return
   }
   Logger.WithFields(fields).Warn(msg)
}

func Error(msg string, fields logrus.Fields) {
   if notInitialized() {
      return
   }
   Logger.WithFields(fields).Error(msg)
}

func Fatal(msg string, fields logrus.Fields) {
   if notInitialized() {
      logrus.Fatal(msg)
      return
   }
   Logger.WithFields(fields).Fatal(msg)
}

func Panic(msg string, fields logrus.Fields) {
   if notInitialized() {
      logrus.Panic(msg)
      return
   }
   Logger.WithFields(fields).Panic(msg)
}

