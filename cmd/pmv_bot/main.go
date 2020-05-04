package main

import (
	"context"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/bot"
	"github.com/dmalyar/pimpmyvocab/dictionary"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/dmalyar/pimpmyvocab/repo"
	"github.com/dmalyar/pimpmyvocab/service"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	logLevelKey        = "log.level"
	logFileKey         = "log.file"
	tokenKey           = "bot.token"
	useProxyKey        = "bot.use-proxy"
	proxyURLKey        = "bot.proxy-url"
	dbUrlKey           = "db.url"
	dictionaryTokenKey = "dictionary.token"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	initViper()

	logger := initLogger()
	defer logger.Close()

	botAPI := initBotAPI(logger)

	vocabRepo := initVocabRepo(logger)
	defer vocabRepo.ClosePool()
	vocabEntryService := initVocabEntryService(logger)
	vocabService := initVocabService(logger, vocabRepo, vocabEntryService)

	b := bot.New(logger, botAPI, vocabService)
	b.Run()
}

func initViper() {
	viper.SetDefault(logLevelKey, "debug")

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.pimpmyvocab")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error reading the config file: %s", err))
	}
}

func initLogger() *log.LoggerLogrus {
	logger := logrus.StandardLogger()

	logFilePath := viper.GetString(logFileKey)
	var file *os.File
	if logFilePath != "" {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		} else {
			logger.Errorf("Error setting logger to write to file: %s", err)
		}
	}

	level, err := logrus.ParseLevel(viper.GetString(logLevelKey))
	if err != nil {
		logger.Errorf("Error parsing log level: %s", err)
		level = logrus.DebugLevel
	}
	logger.SetLevel(level)

	logger.Info("Logger initialized")
	return log.New(logger, file)
}

func initBotAPI(logger log.Logger) *tgbotapi.BotAPI {
	logger.Info("Initializing bot")
	token := viper.GetString(tokenKey)
	if token == "" {
		logger.Panic("Token not found in the config file")
	}

	client := new(http.Client)
	if viper.GetBool(useProxyKey) {
		logger.Info("Using proxy for connecting to bot API")
		proxyRawURL := viper.GetString(proxyURLKey)
		if proxyRawURL == "" {
			logger.Panic("Proxy URL not found in the config file")
		}
		proxyURL, err := url.Parse(proxyRawURL)
		if err != nil {
			logger.Panicf("Error parsing proxy URL: %s", err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	err := tgbotapi.SetLogger(logger)
	if err != nil {
		logger.Errorf("Error setting logger for bot: %s", err)
	}

	botAPI, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		logger.Panicf("Error getting bot API: %s", err)
	}
	logger.Info("Bot initialized")
	return botAPI
}

func initVocabRepo(logger log.Logger) *repo.Postgres {
	logger.Info("Initializing repo")
	dbUrl := viper.GetString(dbUrlKey)
	if dbUrl == "" {
		logger.Panic("DB URL not found in the config file")
	}
	connConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		logger.Panicf("Error parsing DB URL: %s", err)
	}
	connConfig.ConnConfig.Logger = log.NewPgxAdapter(logger)
	dbPool, err := pgxpool.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		logger.Panicf("Error connecting to DB: %s", err)
	}
	logger.Info("Repository initialized")
	return repo.NewPostgresRepo(logger, dbPool)
}

func initVocabEntryService(logger log.Logger) service.VocabEntry {
	logger.Info("Initializing vocab entry service")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	dictionaryToken := viper.GetString(dictionaryTokenKey)
	if dictionaryToken == "" {
		logger.Panic("Dictionary token not found in the config file")
	}
	dictionaryURL := fmt.Sprintf(dictionary.URL, dictionaryToken)
	logger.Info("Vocab entry service initialized")
	return dictionary.NewYandexDict(logger, client, dictionaryURL)
}

func initVocabService(logger log.Logger, vocabRepo repo.Vocab, vocabEntryService service.VocabEntry) service.Vocab {
	localRepoService := service.NewVocabWithLocalRepo(logger, vocabRepo, vocabEntryService)
	return service.NewConcurrentVocab(localRepoService)
}
