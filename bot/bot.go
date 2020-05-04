package bot

import (
	"encoding/json"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/dmalyar/pimpmyvocab/service"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"sort"
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

// Run processes incoming updates from bot api updates channel.
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

func (b *Bot) processMessage(in *tgbotapi.Message) {
	msg := &message{
		id:       in.MessageID,
		chatID:   in.Chat.ID,
		userID:   in.From.ID,
		userName: in.From.UserName,
		text:     in.Text,
	}
	logger := b.logger.WithField("message", msg)
	defer logProcessingTime(logger, time.Now())
	text := strings.ToLower(msg.text)
	switch {
	case text == startCommand:
		b.processStartCommand(logger, msg)
	case text == listCommand:
		b.processListCommand(logger, msg)
	case text == clearCommand:
		b.processClearCommand(logger, msg)
	case text == repeatCommand:
		b.processRepeatCommand(logger, msg)
	case text == quizCommand:
		b.processQuizCommand(logger, msg)
	case strings.HasPrefix(text, "/"):
		logger.Info("Received unsupported command")
	default:
		b.processText(logger, msg)
	}
}

func (b *Bot) processCallback(in *tgbotapi.CallbackQuery) {
	callbackMsg := &callbackMessage{
		id:       in.ID,
		msgID:    in.Message.MessageID,
		chatID:   in.Message.Chat.ID,
		userID:   in.From.ID,
		userName: in.From.UserName,
		data:     new(CallbackData),
	}
	logger := b.logger.WithField("callbackMessage", callbackMsg)
	defer logProcessingTime(logger, time.Now())
	defer b.answerCallback(logger, callbackMsg.id)
	err := json.Unmarshal([]byte(in.Data), callbackMsg.data)
	if err != nil {
		logger.Errorf("Error unmarshalling callback data json: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	switch callbackMsg.data.Command {
	case showFullDescCallbackCmd:
		b.processShowFullDescCommand(logger, callbackMsg)
	case addToVocabCallbackCmd:
		b.processAddToVocabCommand(logger, callbackMsg, false)
	case addToVocabFullDescCallbackCmd:
		b.processAddToVocabCommand(logger, callbackMsg, true)
	case rmFromVocabCallbackCmd:
		b.processRemoveFromVocabCommand(logger, callbackMsg, false)
	case rmFromVocabFullDescCallbackCmd:
		b.processRemoveFromVocabCommand(logger, callbackMsg, true)
	case clearVocabAcceptCallbackCmd:
		b.processClearVocabAnswerCommand(logger, callbackMsg, true)
	case clearVocabDeclineCallbackCmd:
		b.processClearVocabAnswerCommand(logger, callbackMsg, false)
	case repeatCallbackCmd:
		b.processRepeatCallbackCommand(logger, callbackMsg)
	case continueQuizCallbackCmd:
		b.processContinueQuizCommand(logger, callbackMsg)
	case showAnswerCallbackCmd:
		b.processShowAnswerCommand(logger, callbackMsg)
	default:
		logger.Info("Received unsupported callback")
	}
}

func (b *Bot) processStartCommand(logger log.Logger, msg *message) {
	logger.Info("Received /start command")
	_, err := b.vocabService.CreateVocab(msg.userID)
	if err != nil {
		logger.Errorf("Error creating vocab: %s", err)
		b.send(logger, newReply(msg.chatID, techErrReply))
		return
	}
	b.send(logger, newReply(msg.chatID, startReply))
	logger.Info("Processed /start command")
}

func (b *Bot) processListCommand(logger log.Logger, msg *message) {
	logger.Info("Received /list command")
	entries, err := b.vocabService.GetEntriesFromUserVocab(msg.userID)
	if err != nil {
		logger.Errorf("Error getting entries: %s", err)
		b.send(logger, newReply(msg.chatID, techErrReply))
		return
	}
	if len(entries) == 0 {
		logger.Info("Processed /list command (no entries)")
		b.send(logger, newReply(msg.chatID, emptyVocabReply))
		return
	}
	b.send(logger, newReply(msg.chatID, createListReply(entries)))
	logger.Info("Processed /list command")
}

func createListReply(entries []*domain.VocabEntry) string {
	builder := new(strings.Builder)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Text < entries[j].Text
	})
	for _, entry := range entries {
		builder.WriteString(fmt.Sprintf("%s – %s\n", entry.Text, entry.MainTranslation))
	}
	return builder.String()
}

func (b *Bot) processClearCommand(logger log.Logger, msg *message) {
	logger.Info("Received /clear command")
	b.send(logger, newReply(msg.chatID, clearVocabConfirmationReply).withClearConfirmationKeyboard(logger))
	logger.Info("Processed /clear command")
}

func (b *Bot) processRepeatCommand(logger log.Logger, msg *message) {
	logger.Info("Received /repeat command")
	entry, err := b.vocabService.GetRandomEntryFromUserVocab(msg.userID, -1)
	if err != nil {
		logger.Errorf("Error getting random entry: %s", err)
		b.send(logger, newReply(msg.chatID, techErrReply))
		return
	}
	if entry == nil {
		logger.Info("Processed /repeat command (no entries)")
		b.send(logger, newReply(msg.chatID, emptyVocabReply))
		return
	}
	b.send(logger, newReply(msg.chatID, entry.FullDesc(true)).withRepeatKeyboard(logger, entry.ID))
	logger.Info("Processed /repeat command")
}

func (b *Bot) processQuizCommand(logger log.Logger, msg *message) {
	logger.Info("Received /quiz command")
	entry, err := b.vocabService.GetRandomEntryFromUserVocab(msg.userID, -1)
	if err != nil {
		logger.Errorf("Error getting random entry: %s", err)
		b.send(logger, newReply(msg.chatID, techErrReply))
		return
	}
	if entry == nil {
		logger.Info("Processed /quiz command (no entries)")
		b.send(logger, newReply(msg.chatID, emptyVocabReply))
		return
	}
	b.send(logger, newReply(msg.chatID, entry.Text).withQuizKeyboard(logger, entry.ID))
	logger.Info("Processed /quiz command")
}

func (b *Bot) processText(logger log.Logger, msg *message) {
	logger.Info("Received text")
	entry, err := b.vocabService.GetVocabEntryByText(strings.ToLower(msg.text))
	if err != nil {
		logger.Errorf("Error getting vocab entry: %s", err)
		b.send(logger, newReply(msg.chatID, techErrReply))
		return
	}
	if entry == nil {
		logger.Info("Text processed (not found)")
		b.send(logger, newReply(msg.chatID, wordNotFoundReply).withQuote(msg.id))
		return
	}
	inVocab, err := b.vocabService.CheckEntryInUserVocab(entry.ID, msg.userID)
	if err != nil {
		logger.Errorf("Error checking if entry is in the user's vocab: %s", err)
		b.send(logger, newReply(msg.chatID, techErrReply))
		return
	}
	b.send(
		logger,
		newReply(msg.chatID, entry.ShortDesc()).withQuote(msg.id).withShortDescKeyboard(logger, entry.ID, inVocab),
	)
	logger.Info("Text processed")
}

func (b *Bot) processShowFullDescCommand(logger log.Logger, callbackMsg *callbackMessage) {
	logger.Info("Received show full description callback command")
	entry, err := b.vocabService.GetVocabEntryByID(callbackMsg.data.EntryID)
	if err != nil {
		logger.Errorf("Error getting vocab entry: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
	}
	if entry == nil {
		logger.Errorf("Vocab entry not found")
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	logger.WithField("vocabEntry", entry)
	inVocab, err := b.vocabService.CheckEntryInUserVocab(entry.ID, callbackMsg.userID)
	if err != nil {
		logger.Errorf("Error checking if entry is in the user's vocab: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	b.send(
		logger,
		newEditText(callbackMsg.chatID, callbackMsg.msgID, entry.FullDesc(false)).withFullDescKeyboard(logger, entry.ID, inVocab),
	)
	logger.Info("Processed show full description callback command")
}

func (b *Bot) processAddToVocabCommand(logger log.Logger, callbackMsg *callbackMessage, fullDescShown bool) {
	logger.Info("Received add to vocab callback command")
	callbackData := callbackMsg.data
	err := b.vocabService.AddEntryToUserVocab(callbackData.EntryID, callbackMsg.userID)
	if err != nil {
		logger.Errorf("Error adding entry to vocab: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	if fullDescShown {
		b.send(
			logger,
			newEditKeyboardMsgFullDesc(logger, callbackMsg.chatID, callbackMsg.msgID, callbackData.EntryID, true),
		)
	} else {
		b.send(
			logger,
			newEditKeyboardMsgShortDesc(logger, callbackMsg.chatID, callbackMsg.msgID, callbackData.EntryID, true),
		)
	}
	logger.Info("Processed add to vocab callback command")
}

func (b *Bot) processRemoveFromVocabCommand(logger log.Logger, callbackMsg *callbackMessage, fullDescShown bool) {
	logger.Info("Received remove from vocab callback command")
	callbackData := callbackMsg.data
	err := b.vocabService.RemoveEntryFromUserVocab(callbackData.EntryID, callbackMsg.userID)
	if err != nil {
		logger.Errorf("Error removing entry to vocab: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	if fullDescShown {
		b.send(
			logger,
			newEditKeyboardMsgFullDesc(logger, callbackMsg.chatID, callbackMsg.msgID, callbackData.EntryID, false),
		)
	} else {
		b.send(
			logger,
			newEditKeyboardMsgShortDesc(logger, callbackMsg.chatID, callbackMsg.msgID, callbackData.EntryID, false),
		)
	}
	logger.Info("Processed remove from vocab callback command")
}

func (b *Bot) processClearVocabAnswerCommand(logger log.Logger, callbackMsg *callbackMessage, accepted bool) {
	logger.Infof("Received clear vocab answer callback command (%v)", accepted)
	if !accepted {
		b.send(logger, newEditText(callbackMsg.chatID, callbackMsg.msgID, clearVocabDeclinedReply))
		logger.Info("Processed clear vocab answer callback command")
		return
	}
	err := b.vocabService.ClearUserVocab(callbackMsg.userID)
	if err != nil {
		logger.Errorf("Error clearing the user's vocab: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	b.send(logger, newEditText(callbackMsg.chatID, callbackMsg.msgID, clearVocabAcceptedReply))
	logger.Info("Processed clear vocab answer callback command")
}

func (b *Bot) processRepeatCallbackCommand(logger log.Logger, callbackMsg *callbackMessage) {
	logger.Info("Received repeat callback command")
	entry, err := b.vocabService.GetRandomEntryFromUserVocab(callbackMsg.userID, callbackMsg.data.EntryID)
	if err != nil {
		logger.Errorf("Error getting random entry: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	if entry == nil {
		logger.Info("Processed repeat callback command (no entries)")
		b.send(logger, newReply(callbackMsg.chatID, emptyVocabReply))
		return
	}
	b.send(logger, newReply(callbackMsg.chatID, entry.FullDesc(true)).withRepeatKeyboard(logger, entry.ID))
	logger.Info("Processed repeat callback command")
}

func (b *Bot) processContinueQuizCommand(logger log.Logger, callbackMsg *callbackMessage) {
	logger.Info("Received continue quiz callback command")
	entry, err := b.vocabService.GetRandomEntryFromUserVocab(callbackMsg.userID, callbackMsg.data.EntryID)
	if err != nil {
		logger.Errorf("Error getting random entry: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	if entry == nil {
		logger.Info("Processed continue quiz callback command (no entries)")
		b.send(logger, newReply(callbackMsg.chatID, emptyVocabReply))
		return
	}
	b.send(logger, newReply(callbackMsg.chatID, entry.Text).withQuizKeyboard(logger, entry.ID))
	logger.Info("Processed continue quiz callback command")
}

func (b *Bot) processShowAnswerCommand(logger log.Logger, callbackMsg *callbackMessage) {
	logger.Info("Received show answer callback command")
	entry, err := b.vocabService.GetVocabEntryByID(callbackMsg.data.EntryID)
	if err != nil {
		logger.Errorf("Error getting vocab entry: %s", err)
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
	}
	if entry == nil {
		logger.Errorf("Vocab entry not found")
		b.send(logger, newReply(callbackMsg.chatID, techErrReply))
		return
	}
	logger.WithField("vocabEntry", entry)
	b.send(
		logger,
		newEditText(callbackMsg.chatID, callbackMsg.msgID, entry.FullDesc(true)).withQuizKeyboard(logger, entry.ID),
	)
	logger.Info("Processed show answer callback command")
}

func logProcessingTime(logger log.Logger, start time.Time) {
	logger.Debugf("Processing time: %s", time.Since(start))
}
