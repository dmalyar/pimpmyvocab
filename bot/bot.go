package bot

import (
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/dmalyar/pimpmyvocab/service"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
	"time"
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
		if update.Message != nil {
			go b.processMessage(update.Message)
			continue
		}
		if update.CallbackQuery != nil {
			go b.processCallback(update.CallbackQuery)
			continue
		}
		b.logger.Info("Received update of unsupported type")
	}
}

func (b *Bot) processMessage(msg *tgbotapi.Message) {
	user := msg.From
	chat := msg.Chat
	contextLog := b.logger.WithFields(map[string]interface{}{
		"userID":   user.ID,
		"userName": user.UserName,
		"chatID":   chat.ID,
		"msgID":    msg.MessageID,
		"text":     msg.Text,
	})
	defer logProcessingTime(contextLog, time.Now())
	text := strings.ToLower(msg.Text)
	switch {
	case text == startCommand:
		b.processStartCommand(contextLog, chat.ID, user.ID)
	case strings.HasPrefix(text, "/"):
		contextLog.Info("Received unsupported command")
	default:
		b.processText(contextLog, chat.ID, msg.MessageID, text)
	}
}

func (b *Bot) processCallback(callback *tgbotapi.CallbackQuery) {
	user := callback.From
	msg := callback.Message
	chat := msg.Chat
	data := callback.Data
	contextLog := b.logger.WithFields(map[string]interface{}{
		"userID":       user.ID,
		"userName":     user.UserName,
		"chatID":       chat.ID,
		"msgID":        msg.MessageID,
		"callbackData": data,
	})
	defer logProcessingTime(contextLog, time.Now())
	switch {
	case strings.HasPrefix(data, showFullDescCommand):
		entryID, err := strconv.Atoi(strings.TrimPrefix(data, showFullDescCommand))
		if err != nil {
			contextLog.Errorf("Error parsing callback data: %s", err)
			b.send(contextLog, newReply(chat.ID, techErrReply))
			return
		}
		b.processShowFullDescCommand(contextLog, chat.ID, msg.MessageID, entryID)
	default:
		contextLog.Info("Received unsupported callback")
	}
}

func (b *Bot) processStartCommand(contextLog log.Logger, chatID int64, userID int) {
	contextLog.Info("Received /start command")
	replyText := startReply
	_, err := b.vocabService.CreateVocab(userID)
	if err != nil {
		contextLog.Errorf("Error creating vocab: %s", err)
		replyText = techErrReply
	}
	b.send(contextLog, newReply(chatID, replyText))
	contextLog.Info("Processed /start command")
}

func (b *Bot) processText(contextLog log.Logger, chatID int64, msgID int, text string) {
	contextLog.Info("Received text")
	entry, err := b.vocabService.GetVocabEntryByText(text)
	if err != nil {
		contextLog.Errorf("Error getting vocab entry: %s", err)
		b.send(contextLog, newReply(chatID, techErrReply))
		return
	}
	if entry == nil {
		contextLog.Info("Text processed (not found)")
		b.send(contextLog, newReply(chatID, wordNotFoundReply).withQuote(msgID))
		return
	}
	b.send(contextLog, newReply(chatID, entry.ShortDesc()).withQuote(msgID).withShortDescKeyboard(entry))
	contextLog.Info("Text processed")
}

func (b *Bot) processShowFullDescCommand(contextLog log.Logger, chatID int64, msgID, entryID int) {
	contextLog.Info("Received /showfulldesc callback command")
	entry, err := b.vocabService.GetVocabEntryByID(entryID)
	if err != nil {
		contextLog.Errorf("Error getting vocab entry: %s", err)
		b.send(contextLog, newReply(chatID, techErrReply))
	}
	if entry == nil {
		contextLog.Errorf("Vocab entry not found")
		b.send(contextLog, newReply(chatID, techErrReply))
		return
	}
	contextLog.WithField("vocabEntry", entry)
	b.send(contextLog, newEditText(chatID, msgID, entry.FullDesc()))
	contextLog.Info("Processed /showfulldesc callback command")
}

func (b *Bot) send(contextLog log.Logger, msg tgbotapi.Chattable) {
	contextLog = contextLog.WithField("msgToSend", msg)
	_, err := b.api.Send(msg)
	if err != nil {
		contextLog.Errorf("Error sending message: %s", err)
		return
	}
	contextLog.Debugf("Message sent")
}

type replyMsg struct {
	*tgbotapi.MessageConfig
	quoteFlag, keyboardFlag bool
}

func newReply(chatID int64, text string) *replyMsg {
	msg := tgbotapi.NewMessage(chatID, text)
	return &replyMsg{MessageConfig: &msg}
}

func (m *replyMsg) withQuote(msgID int) *replyMsg {
	m.ReplyToMessageID = msgID
	m.quoteFlag = true
	return m
}

func (m *replyMsg) withShortDescKeyboard(ve *domain.VocabEntry) *replyMsg {
	m.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(showFullDescButton, showFullDescCommand+strconv.Itoa(ve.ID)),
		),
	)
	m.keyboardFlag = true
	return m
}

func (m *replyMsg) String() string {
	return fmt.Sprintf("New message (with quote = %v; with keyboard = %v): %s", m.quoteFlag, m.keyboardFlag, m.Text)
}

type editTextMsg struct {
	*tgbotapi.EditMessageTextConfig
}

func newEditText(chatID int64, msgID int, text string) *editTextMsg {
	msg := tgbotapi.NewEditMessageText(chatID, msgID, text)
	return &editTextMsg{&msg}
}

func (m *editTextMsg) String() string {
	return fmt.Sprintf("Edit message text (msgID = %v): %s", m.MessageID, m.Text)
}

func logProcessingTime(contextLog log.Logger, start time.Time) {
	contextLog.Debugf("Processing time: %s", time.Since(start))
}
