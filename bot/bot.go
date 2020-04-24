package bot

import (
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/dmalyar/pimpmyvocab/service"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	startCommand = "/start"
)

type Bot struct {
	logger       log.Logger
	api          *tgbotapi.BotAPI
	vocabService service.Vocab
}

func New(logger log.Logger, api *tgbotapi.BotAPI, vocabService service.Vocab) *Bot {
	return &Bot{
		logger:       logger,
		api:          api,
		vocabService: vocabService,
	}
}

// Run processes incoming updates from api updates channel.
func (b *Bot) Run() {
	updates, err := b.api.GetUpdatesChan(
		tgbotapi.UpdateConfig{
			Offset:  0,
			Timeout: 60,
		},
	)
	if err != nil {
		b.logger.Panicf("Error getting bot updates chan: %s", err)
	}
	b.logger.Info("Successfully got updates channel and start processing messages")
	for update := range updates {
		go b.process(update)
	}
}

func (b *Bot) process(update tgbotapi.Update) {
	msg := update.Message
	if msg == nil {
		b.logger.Debug("Received update with nil message")
		return
	}
	user := msg.From
	chat := msg.Chat
	contextLog := b.logger.WithFields(map[string]interface{}{
		"userID":   user.ID,
		"userName": user.UserName,
		"chatID":   chat.ID,
		"msgID":    msg.MessageID,
	})
	text := msg.Text
	switch text {
	case startCommand:
		b.processStartCommand(contextLog, chat.ID, user.ID)
	default:
		contextLog.Debug("Received unsupported message")
	}
}

func (b *Bot) processStartCommand(contextLog log.Logger, chatID int64, userID int) {
	contextLog.Info("Received /start command")
	replyText := startReplySucc
	_, err := b.vocabService.CreateVocab(userID)
	if err != nil {
		contextLog.Errorf("Error processing start command: %s", err)
		replyText = startReplyErr
	}
	b.reply(contextLog, chatID, replyText, err)
	contextLog.Info("Processed /start command")
}

func (b *Bot) reply(contextLog log.Logger, chatID int64, replyMsg string, err error) {
	reply := tgbotapi.NewMessage(chatID, replyMsg)
	_, err = b.api.Send(reply)
	if err != nil {
		contextLog.Errorf("Error sending reply: %s", err)
	}
}
