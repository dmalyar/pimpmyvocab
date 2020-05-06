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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const (
	logLevelKey        = "log.level"
	logFilePathKey     = "log.file"
	tokenKey           = "bot.token"
	useProxyKey        = "bot.use-proxy"
	proxyURLKey        = "bot.proxy-url"
	dbUrlKey           = "db.url"
	dbMigrationPathKey = "db.migration-path"
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
	viper.AddConfigPath("$HOME/.pimpmyvocab") // local
	viper.AddConfigPath("/pimpmyvocab")       // docker
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Error reading the config file: %s", err))
	}
}

func initLogger() *log.LoggerLogrus {
	logger := logrus.StandardLogger()

	logFilePath := viper.GetString(logFilePathKey)
	var lumberjackLogger *lumberjack.Logger
	if logFilePath != "" {
		lumberjackLogger = &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    50,
			MaxBackups: 10,
			MaxAge:     0,
			LocalTime:  true,
		}
		logrus.SetOutput(lumberjackLogger)
	}

	level, err := logrus.ParseLevel(viper.GetString(logLevelKey))
	if err != nil {
		logger.Errorf("Error parsing log level: %s", err)
		level = logrus.DebugLevel
	}
	logger.SetLevel(level)

	logger.Info("Logger initialized")
	return log.New(logger, lumberjackLogger)
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
	migrateSchema(logger, dbUrl)
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

func migrateSchema(logger log.Logger, dbUrl string) {
	logger.Info("Migrating DB schema")
	migrationPath := viper.GetString(dbMigrationPathKey)
	if migrationPath == "" {
		logger.Panic("DB migration path not found in the config file")
	}
	m, err := migrate.New(
		"file://"+migrationPath,
		dbUrl,
	)
	if err != nil {
		logger.Panicf("Error migrating schema: %s", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Panicf("Error migrating schema: %s", err)
	}
	logger.Info("DB schema migrated")
}

func initVocabEntryService(logger log.Logger) *dictionary.Yandex {
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

func initVocabService(logger log.Logger, vocabRepo repo.Vocab, vocabEntryService service.VocabEntry) *service.ConcurrentVocab {
	localRepoService := service.NewVocabWithLocalRepo(logger, vocabRepo, vocabEntryService)
	return service.NewConcurrentVocab(localRepoService)
}
