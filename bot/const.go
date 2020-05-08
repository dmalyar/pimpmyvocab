package bot

const (
	startCommand  = "/start"
	helpCommand   = "/help"
	listCommand   = "/list"
	clearCommand  = "/clear"
	repeatCommand = "/repeat"
	quizCommand   = "/quiz"

	helpReply = "Теперь у вас в телеграме есть личный словарь для изучения английского языка!\n\n" +
		"Пришлите боту английское слово, чтобы получить по нему краткую словарную статью " +
		"с возможностью добавить её в свой словарь.\n\n" +
		"Используйте команду /list для просмотра словаря.\n\n" +
		"Команда /repeat поможет вам закрепить знания.\n\n" +
		"Команду /quiz используйте для проверки своих знаний.\n\n" +
		"Команда /clear – очистка словаря. Не волнуйтесь, бот уточнит ваше намерение начать всё с чистого листа.\n\n" +
		"Если слово не отображается в словаре после добавления, повторно отправьте команду /start"
	techErrReply = "Кажется, у бота технические проблемы :(\n" +
		"Попробуйте повторить запрос позже. А мы пока поменяем ему масло."
	offlineReply = "Наверное, вы заметили, что какое-то время наш бот отдыхал и не мог обрабатывать ваши запросы.\n" +
		"Теперь он снова в строю!"
	emptyVocabReply             = "В вашем словаре пока нет записей.\nНо ведь это легко исправить ;)"
	clearVocabConfirmationReply = "Вы уверены, что хотите удалить все записи из своего словаря?"
	clearVocabDeclinedReply     = "Вот и правильно, отличный же словарь!"
	clearVocabAcceptedReply     = "Готово! Начните с чистого листа!"
	wordNotFoundReply           = "А вы точно продюсер? А это точно английское слово?\n" +
		"Просто бот по нему ничего не нашёл :("

	showFullDescButton    = "Все варианты перевода"
	addToVocabButton      = "Добавить в словарь"
	removeFromVocabButton = "Удалить из словаря"
	yesButton             = "Да"
	noButton              = "Нет"
	newWordButton         = "Новое слово"
	showAnswerButton      = "Показать перевод"
)

type CallbackCommand int

const (
	showFullDescCallbackCmd CallbackCommand = iota
	addToVocabCallbackCmd
	addToVocabFullDescCallbackCmd
	rmFromVocabCallbackCmd
	rmFromVocabFullDescCallbackCmd
	clearVocabAcceptCallbackCmd
	clearVocabDeclineCallbackCmd
	repeatCallbackCmd
	continueQuizCallbackCmd
	showAnswerCallbackCmd
)
