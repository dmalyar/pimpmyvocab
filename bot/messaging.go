package bot

import (
	"encoding/json"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type message struct {
	id       int
	chatID   int64
	userID   int
	userName string
	text     string
}

func (m *message) String() string {
	return fmt.Sprintf("id: %v; chatID: %v; userID: %v; userName: %v; text: %s",
		m.id, m.chatID, m.userID, m.userName, m.text)
}

type callbackMessage struct {
	id       string
	msgID    int
	chatID   int64
	userID   int
	userName string
	data     *CallbackData
}

func (m *callbackMessage) String() string {
	return fmt.Sprintf("id: %s; msgID: %v; chatID: %v; userID: %v; userName: %v; data: {%s}",
		m.id, m.msgID, m.chatID, m.userID, m.userName, m.data)
}

type CallbackData struct {
	Command CallbackCommand
	EntryID int
}

func (c *CallbackData) String() string {
	return fmt.Sprintf("Command: %v; EntryID: %v", c.Command, c.EntryID)
}

func (b *Bot) send(logger log.Logger, msg tgbotapi.Chattable) {
	logger = logger.WithField("msgToSend", msg)
	_, err := b.api.Send(msg)
	if err != nil {
		logger.Errorf("Error sending message: %s", err)
		return
	}
	logger.Debug("Message sent")
}

func (b *Bot) answerCallback(logger log.Logger, id string) {
	_, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(id, ""))
	if err != nil {
		logger.Errorf("Error answering callback: %s", err)
	}
	logger.Debug("Callback answered")
}

type replyMsg struct {
	*tgbotapi.MessageConfig
	quoteFlag, keyboardFlag bool
}

func (m *replyMsg) String() string {
	return fmt.Sprintf("New message (with quote = %v; with keyboard = %v): %s", m.quoteFlag, m.keyboardFlag, m.Text)
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

func (m *replyMsg) withShortDescKeyboard(logger log.Logger, entryID int, inVocab bool) *replyMsg {
	keyboard, err := shortDescKeyboard(entryID, inVocab)
	if err != nil {
		logger.Errorf("Error generating short desc keyboard: %s", err)
		return m
	}
	m.ReplyMarkup = keyboard
	m.keyboardFlag = true
	return m
}

func (m *replyMsg) withClearConfirmationKeyboard(logger log.Logger) *replyMsg {
	clearVocabAcceptCallback, err := json.Marshal(CallbackData{Command: clearVocabAcceptCallbackCmd})
	if err != nil {
		logger.Errorf("Error generating clear confirmation keyboard: %s", err)
		return m
	}
	clearVocabDeclineCallback, err := json.Marshal(CallbackData{Command: clearVocabDeclineCallbackCmd})
	if err != nil {
		logger.Errorf("Error generating clear confirmation keyboard: %s", err)
		return m
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(yesButton, string(clearVocabAcceptCallback)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(noButton, string(clearVocabDeclineCallback)),
		),
	)
	m.ReplyMarkup = keyboard
	m.keyboardFlag = true
	return m
}

func (m *replyMsg) withRepeatKeyboard(logger log.Logger, entryID int) *replyMsg {
	callback, err := json.Marshal(CallbackData{
		Command: repeatCallbackCmd,
		EntryID: entryID,
	})
	if err != nil {
		logger.Errorf("Error generating repeat keyboard: %s", err)
		return m
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newWordButton, string(callback)),
		),
	)
	m.ReplyMarkup = keyboard
	m.keyboardFlag = true
	return m
}

func (m *replyMsg) withQuizKeyboard(logger log.Logger, entryID int) *replyMsg {
	continueCallback, err := json.Marshal(CallbackData{
		Command: continueQuizCallbackCmd,
		EntryID: entryID,
	})
	if err != nil {
		logger.Errorf("Error generating quiz keyboard: %s", err)
		return m
	}
	showAnswerCallback, err := json.Marshal(CallbackData{
		Command: showAnswerCallbackCmd,
		EntryID: entryID,
	})
	if err != nil {
		logger.Errorf("Error generating quiz keyboard: %s", err)
		return m
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(showAnswerButton, string(showAnswerCallback)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newWordButton, string(continueCallback)),
		),
	)
	m.ReplyMarkup = keyboard
	m.keyboardFlag = true
	return m
}

type editTextMsg struct {
	*tgbotapi.EditMessageTextConfig
	keyboardFlag bool
}

func (m *editTextMsg) String() string {
	return fmt.Sprintf("Edit message text (msgID = %v; with keyboard = %v): %s", m.MessageID, m.keyboardFlag, m.Text)
}

func newEditText(chatID int64, msgID int, text string) *editTextMsg {
	msg := tgbotapi.NewEditMessageText(chatID, msgID, text)
	return &editTextMsg{EditMessageTextConfig: &msg}
}

func (m *editTextMsg) withFullDescKeyboard(logger log.Logger, entryID int, inVocab bool) *editTextMsg {
	keyboard, err := fullDescKeyboard(entryID, inVocab)
	if err != nil {
		logger.Errorf("Error generating full desc keyboard: %s", err)
		return m
	}
	m.ReplyMarkup = keyboard
	m.keyboardFlag = true
	return m
}

func (m *editTextMsg) withQuizKeyboard(logger log.Logger, entryID int) *editTextMsg {
	callback, err := json.Marshal(CallbackData{
		Command: continueQuizCallbackCmd,
		EntryID: entryID,
	})
	if err != nil {
		logger.Errorf("Error generating quiz keyboard: %s", err)
		return m
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(newWordButton, string(callback)),
		),
	)
	m.ReplyMarkup = &keyboard
	m.keyboardFlag = true
	return m
}

type editKeyboardMsg struct {
	*tgbotapi.EditMessageReplyMarkupConfig
}

func (m *editKeyboardMsg) String() string {
	return fmt.Sprintf("Edit keyboard message (msgID = %v)", m.MessageID)
}

func newEditKeyboardMsgShortDesc(logger log.Logger, chatID int64, msgID, entryID int, inVocab bool) *editKeyboardMsg {
	keyboard, err := shortDescKeyboard(entryID, inVocab)
	if err != nil {
		logger.Errorf("Error generating short desc keyboard: %s", err)
		msg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, tgbotapi.InlineKeyboardMarkup{})
		return &editKeyboardMsg{EditMessageReplyMarkupConfig: &msg}
	}
	msg := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, *keyboard)
	return &editKeyboardMsg{EditMessageReplyMarkupConfig: &msg}
}

func newEditKeyboardMsgFullDesc(logger log.Logger, chatID int64, msgID, entryID int, inVocab bool) *editKeyboardMsg {
	keyboard, err := fullDescKeyboard(entryID, inVocab)
	if err != nil {
		logger.Errorf("Error generating full desc keyboard: %s", err)
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
		vocabActionCallback, err = json.Marshal(CallbackData{
			EntryID: entryID,
			Command: rmFromVocabCallbackCmd,
		})
	} else {
		vocabActionButton = addToVocabButton
		vocabActionCallback, err = json.Marshal(CallbackData{
			EntryID: entryID,
			Command: addToVocabCallbackCmd,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("error marshalling vocab action callback json for short desc keyboard: %s", err)
	}
	showFullDescCallback, err := json.Marshal(CallbackData{
		EntryID: entryID,
		Command: showFullDescCallbackCmd,
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
		vocabActionCallback, err = json.Marshal(CallbackData{
			EntryID: entryID,
			Command: rmFromVocabFullDescCallbackCmd,
		})
	} else {
		vocabActionButton = addToVocabButton
		vocabActionCallback, err = json.Marshal(CallbackData{
			EntryID: entryID,
			Command: addToVocabFullDescCallbackCmd,
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
