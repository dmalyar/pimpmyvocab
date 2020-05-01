package bot

import (
	"encoding/json"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/dmalyar/pimpmyvocab/service"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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
		b.processText(contextLog, chat.ID, msg.MessageID, user.ID, text)
	}
}

func (b *Bot) processCallback(callback *tgbotapi.CallbackQuery) {
	user := callback.From
	msg := callback.Message
	chat := msg.Chat
	rawData := callback.Data
	contextLog := b.logger.WithFields(map[string]interface{}{
		"userID":          user.ID,
		"userName":        user.UserName,
		"chatID":          chat.ID,
		"msgID":           msg.MessageID,
		"callbackRawData": rawData,
	})
	defer logProcessingTime(contextLog, time.Now())
	callbackData := new(Callback)
	err := json.Unmarshal([]byte(rawData), callbackData)
	if err != nil {
		contextLog.Errorf("Error unmarshalling callback data json: %s", err)
		b.send(contextLog, newReply(chat.ID, techErrReply))
		return
	}
	switch callbackData.Command {
	case showFullDesc:
		b.processShowFullDescCommand(contextLog, chat.ID, msg.MessageID, user.ID, callbackData.EntryID)
	case addToVocab:
		b.processAddToVocabCommand(contextLog, chat.ID, msg.MessageID, user.ID, callbackData.EntryID, false)
	case addToVocabFullDesc:
		b.processAddToVocabCommand(contextLog, chat.ID, msg.MessageID, user.ID, callbackData.EntryID, true)
	case removeFromVocab:
		b.processRemoveFromVocabCommand(contextLog, chat.ID, msg.MessageID, user.ID, callbackData.EntryID, false)
	case removeFromVocabFullDesc:
		b.processRemoveFromVocabCommand(contextLog, chat.ID, msg.MessageID, user.ID, callbackData.EntryID, true)
	default:
		contextLog.Info("Received unsupported callback")
	}
	b.answerCallback(contextLog, callback.ID)
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

func (b *Bot) processText(contextLog log.Logger, chatID int64, msgID, userID int, text string) {
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
	inVocab, err := b.vocabService.CheckEntryInUserVocab(entry.ID, userID)
	if err != nil {
		contextLog.Errorf("Error checking if entry is in the user's vocab: %s", err)
		b.send(contextLog, newReply(chatID, techErrReply))
		return
	}
	b.send(contextLog, newReply(chatID, entry.ShortDesc()).withQuote(msgID).withShortDescKeyboard(contextLog, entry.ID, inVocab))
	contextLog.Info("Text processed")
}

func (b *Bot) processShowFullDescCommand(contextLog log.Logger, chatID int64, msgID, userID, entryID int) {
	contextLog.Info("Received show full description callback command")
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
	inVocab, err := b.vocabService.CheckEntryInUserVocab(entry.ID, userID)
	if err != nil {
		contextLog.Errorf("Error checking if entry is in the user's vocab: %s", err)
		b.send(contextLog, newReply(chatID, techErrReply))
		return
	}
	b.send(contextLog, newEditText(chatID, msgID, entry.FullDesc()).WithFullDescKeyboard(contextLog, entry.ID, inVocab))
	contextLog.Info("Processed show full description callback command")
}

func (b *Bot) processAddToVocabCommand(contextLog log.Logger, chatID int64, msgID, userID, entryID int, fullDescShown bool) {
	contextLog.Info("Received add to vocab callback command")
	err := b.vocabService.AddEntryToUserVocab(entryID, userID)
	if err != nil {
		contextLog.Errorf("Error adding entry to vocab: %s", err)
		b.send(contextLog, newReply(chatID, techErrReply))
		return
	}
	if fullDescShown {
		b.send(contextLog, newEditKeyboardMsgFullDesc(contextLog, chatID, msgID, entryID, true))
	} else {
		b.send(contextLog, newEditKeyboardMsgShortDesc(contextLog, chatID, msgID, entryID, true))
	}
	contextLog.Info("Processed add to vocab callback command")
}

func (b *Bot) processRemoveFromVocabCommand(contextLog log.Logger, chatID int64, msgID, userID, entryID int, fullDescShown bool) {
	contextLog.Info("Received remove from vocab callback command")
	err := b.vocabService.RemoveEntryFromUserVocab(entryID, userID)
	if err != nil {
		contextLog.Errorf("Error removing entry to vocab: %s", err)
		b.send(contextLog, newReply(chatID, techErrReply))
		return
	}
	if fullDescShown {
		b.send(contextLog, newEditKeyboardMsgFullDesc(contextLog, chatID, msgID, entryID, false))
	} else {
		b.send(contextLog, newEditKeyboardMsgShortDesc(contextLog, chatID, msgID, entryID, false))
	}
	contextLog.Info("Processed remove from vocab callback command")
}

func (b *Bot) send(contextLog log.Logger, msg tgbotapi.Chattable) {
	contextLog = contextLog.WithField("msgToSend", msg)
	_, err := b.api.Send(msg)
	if err != nil {
		contextLog.Errorf("Error sending message: %s", err)
		return
	}
	contextLog.Debug("Message sent")
}

func (b *Bot) answerCallback(contextLog log.Logger, id string) {
	_, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(id, ""))
	if err != nil {
		contextLog.Errorf("Error answering callback: %s", err)
	}
	contextLog.Debug("Callback answered")
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

func (m *replyMsg) withShortDescKeyboard(contextLog log.Logger, entryID int, inVocab bool) *replyMsg {
	keyboard, err := shortDescKeyboard(entryID, inVocab)
	if err != nil {
		contextLog.Errorf("Error generating short desc keyboard: %s", err)
		return m
	}
	m.ReplyMarkup = keyboard
	m.keyboardFlag = true
	return m
}

func (m *replyMsg) String() string {
	return fmt.Sprintf("New message (with quote = %v; with keyboard = %v): %s", m.quoteFlag, m.keyboardFlag, m.Text)
}

type editTextMsg struct {
	*tgbotapi.EditMessageTextConfig
	keyboardFlag bool
}

func newEditText(chatID int64, msgID int, text string) *editTextMsg {
	msg := tgbotapi.NewEditMessageText(chatID, msgID, text)
	return &editTextMsg{EditMessageTextConfig: &msg}
}

func (m *editTextMsg) WithFullDescKeyboard(contextLog log.Logger, entryID int, inVocab bool) *editTextMsg {
	keyboard, err := fullDescKeyboard(entryID, inVocab)
	if err != nil {
		contextLog.Errorf("Error generating full desc keyboard: %s", err)
		return m
	}
	m.ReplyMarkup = keyboard
	m.keyboardFlag = true
	return m
}

type editKeyboardMsg struct {
	*tgbotapi.EditMessageReplyMarkupConfig
}

func newEditKeyboardMsgShortDesc(contextLog log.Logger, chatID int64, msgID, entryID int, inVocab bool) *editKeyboardMsg {
	keyboard, err := shortDescKeyboard(entryID, inVocab)
	if err != nil {
		contextLog.Errorf("Error generating short desc keyboard: %s", err)
		msg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		return &editKeyboardMsg{EditMessageReplyMarkupConfig: &msg}
	}
	msg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, *keyboard)
	return &editKeyboardMsg{EditMessageReplyMarkupConfig: &msg}
}

func newEditKeyboardMsgFullDesc(contextLog log.Logger, chatID int64, msgID, entryID int, inVocab bool) *editKeyboardMsg {
	keyboard, err := fullDescKeyboard(entryID, inVocab)
	if err != nil {
		contextLog.Errorf("Error generating full desc keyboard: %s", err)
		msg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		return &editKeyboardMsg{EditMessageReplyMarkupConfig: &msg}
	}
	msg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, *keyboard)
	return &editKeyboardMsg{EditMessageReplyMarkupConfig: &msg}
}

func shortDescKeyboard(entryID int, inVocab bool) (*tgbotapi.InlineKeyboardMarkup, error) {
	var vocabActionButton string
	var vocabActionCallback []byte
	var err error
	if inVocab {
		vocabActionButton = removeFromVocabButton
		vocabActionCallback, err = json.Marshal(Callback{
			EntryID: entryID,
			Command: removeFromVocab,
		})
	} else {
		vocabActionButton = addToVocabButton
		vocabActionCallback, err = json.Marshal(Callback{
			EntryID: entryID,
			Command: addToVocab,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("error marshalling vocab action callback json for short desc keyboard: %s", err)
	}
	showFullDescCallback, err := json.Marshal(Callback{
		EntryID: entryID,
		Command: showFullDesc,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling show full desc action callback json for short desc keyboard: %s", err)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(showFullDescButton, string(showFullDescCallback)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocabActionButton, string(vocabActionCallback)),
		),
	)
	return &keyboard, nil
}

func fullDescKeyboard(entryID int, inVocab bool) (*tgbotapi.InlineKeyboardMarkup, error) {
	var vocabActionButton string
	var vocabActionCallback []byte
	var err error
	if inVocab {
		vocabActionButton = removeFromVocabButton
		vocabActionCallback, err = json.Marshal(Callback{
			EntryID: entryID,
			Command: removeFromVocabFullDesc,
		})
	} else {
		vocabActionButton = addToVocabButton
		vocabActionCallback, err = json.Marshal(Callback{
			EntryID: entryID,
			Command: addToVocabFullDesc,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("marshalling callback json for full desc keyboard: %s", err)
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(vocabActionButton, string(vocabActionCallback)),
		),
	)
	return &keyboard, nil
}

type CallbackCommand int

const (
	showFullDesc CallbackCommand = iota
	addToVocab
	addToVocabFullDesc
	removeFromVocab
	removeFromVocabFullDesc
)

type Callback struct {
	EntryID int
	Command CallbackCommand
}

func (m *editTextMsg) String() string {
	return fmt.Sprintf("Edit message text (msgID = %v; with keyboard = %v): %s", m.MessageID, m.keyboardFlag, m.Text)
}

func logProcessingTime(contextLog log.Logger, start time.Time) {
	contextLog.Debugf("Processing time: %s", time.Since(start))
}
