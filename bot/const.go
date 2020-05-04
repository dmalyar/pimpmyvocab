package bot

const (
	startCommand  = "/start"
	listCommand   = "/list"
	clearCommand  = "/clear"
	repeatCommand = "/repeat"

	startReply                  = "Start successful reply" // TODO change text
	techErrReply                = "Кажется, у бота технические проблемы :(\nПопробуйте повторить запрос позже. А мы пока поменяем ему масло."
	emptyVocabReply             = "В вашем словаре пока нет записей.\nНо ведь это легко исправить ;)"
	clearVocabConfirmationReply = "Вы уверены, что хотите удалить все записи из своего словаря?"
	clearVocabDeclinedReply     = "Вот и правильно, отличный же словарь!"
	clearVocabAcceptedReply     = "Готово! Начните с чистого листа!"
	wordNotFoundReply           = "А вы точно продюссер? А это точно английское слово?\nПросто мы по нему ничего не нашли :("

	showFullDescButton    = "Все варианты перевода"
	addToVocabButton      = "Добавить в словарь"
	removeFromVocabButton = "Удалить из словаря"
	yesButton             = "Да"
	noButton              = "Нет"
	newWordButton         = "Новое слово"
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
)
