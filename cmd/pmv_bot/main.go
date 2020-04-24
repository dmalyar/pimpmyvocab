package main

import (
	"context"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/bot"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/dmalyar/pimpmyvocab/repo"
	"github.com/dmalyar/pimpmyvocab/service"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
)

const (
	logLevelKey = "log.level"
	logFileKey  = "log.file"
	tokenKey    = "bot.token"
	useProxyKey = "bot.use-proxy"
	proxyURLKey = "bot.proxy-url"
	dbUrlKey    = "db.url"
)

func main() {
	initViper()

	logger := initLogger()
	defer logger.Close()

	botAPI := initBotAPI(logger)

	vocabRepo := initVocabRepo(logger)
	defer vocabRepo.ClosePool()
	vocabService := initVocabService(logger, vocabRepo)

	b := bot.New(logger, botAPI, vocabService)
	b.Run()
}

func initViper() {
	viper.SetDefault(logLevelKey, "debug")

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.pimpmyvocab")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %s", err))
	}
}

func initLogger() *log.LoggerLogrus {
	logger := logrus.StandardLogger()

	level, err := logrus.ParseLevel(viper.GetString(logLevelKey))
	if err != nil {
		logger.Panicf("Error parsing log level: %s", err)
	}
	logger.SetLevel(level)

	logFilePath := viper.GetString(logFileKey)
	var file *os.File
	if logFilePath != "" {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		} else {
			logger.Errorf("Failed to log to file, using default stderr: %s\n", err)
		}
	}
	logger.Info("Logger initialized")
	return log.New(logger, file)
}

func initBotAPI(logger log.Logger) *tgbotapi.BotAPI {
	logger.Info("Initializing bot")
	token := viper.GetString(tokenKey)
	if token == "" {
		logger.Panic("Token is not found in config file")
	}

	client := new(http.Client)
	if viper.GetBool(useProxyKey) {
		proxyRawURL := viper.GetString(proxyURLKey)
		if proxyRawURL == "" {
			logger.Panic("Proxy URL is not found in config file")
		}
		proxyURL, err := url.Parse(proxyRawURL)
		if err != nil {
			logger.Panicf("Error parsing proxy URL: %s", err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	err := tgbotapi.SetLogger(logger)
	if err != nil {
		logger.Errorf("Error setting logger for bot: %s\n", err)
	}

	botAPI, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		logger.Panic(err)
	}
	logger.Info("Bot initialized")
	return botAPI
}

func initVocabRepo(logger log.Logger) *repo.Postgres {
	logger.Info("Initializing repo")
	dbUrl := viper.GetString(dbUrlKey)
	if dbUrl == "" {
		logger.Panic("DB URL is not found in config file")
	}
	connConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		logger.Panicf("Error parsing DB URL: %s\n", err)
	}
	connConfig.ConnConfig.Logger = log.NewPgxAdapter(logger)
	dbPool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		logger.Panicf("Error connecting to DB: %s\n", err)
	}
	logger.Info("Repository initialized")
	return repo.NewPostgresRepo(dbPool)
}

func initVocabService(logger log.Logger, vocabRepo repo.Vocab) service.Vocab {
	localRepoService := service.NewVocabWithLocalRepo(logger, vocabRepo)
	return service.NewConcurrentVocab(localRepoService)
}
