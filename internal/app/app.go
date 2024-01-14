package app

import(
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/Senyacka/Go_motivate_tg_bot/internal/models"
	"github.com/Senyacka/Go_motivate_tg_bot/internal/config"
)

var gBot *tgbotapi.BotAPI
var gToken string
var gChatId int64

var gUsersInChat models.Users

func getUserFromUpdate(update *tgbotapi.Update) (user *models.User, found bool){
	if update.CallbackQuery == nil {
		return
	}

	userId := update.CallbackQuery.From.ID
	for _, userInChat := range gUsersInChat {
		if userInChat.Id == userId {
			return userInChat, true
		}
	}
	return 
}

func storeUserFromUpdate(update *tgbotapi.Update) (user *models.User, found bool){
	if update.CallbackQuery == nil {
		return
	}

	userId := update.CallbackQuery.From.ID
	userName := update.CallbackQuery.From.UserName
	user = &models.User{Id: userId, Name: userName}
	gUsersInChat = append(gUsersInChat, user)
	return user, true
}

func delay(delayInSec uint8) {
	time.Sleep(time.Second * time.Duration(delayInSec))
}

func isCallbackQuery(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data != ""
}

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == "start"
}

func printSystemMessageWithDelay(delayInSec uint8, message string) {
	gBot.Send(tgbotapi.NewMessage(gChatId, message))
	delay(delayInSec)
}

func askToPrintIntro() {
	msg := tgbotapi.NewMessage(gChatId, "Хочешь чтобы я рассказал что может этот бот, и для чего он нужен?")
	keyboard := tgbotapi.NewInlineKeyboardRow(
		getKeyboardButton(config.BUTTON_TEXT_PRINT_INTRO, config.BUTTON_CODE_PRINT_INTRO),
		getKeyboardButton(config.BUTTON_TEXT_SKIP_INTRO, config.BUTTON_CODE_SKIP_INTRO),
	)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard)
	gBot.Send(msg)
}

func printIntro(update *tgbotapi.Update) {
	printSystemMessageWithDelay(1, "Привет! Я ХХХХХХ.")
	printSystemMessageWithDelay(2, "Моя задача - начислять тебе монеты за выполнение повседневных задач.")
	printSystemMessageWithDelay(3, "К сожалению я никак не смогу проверить, сделана задача или нет, так что надеюсь на твою честность!")
	printSystemMessageWithDelay(2, "Нажми на задачу ниже, чтобы отметить её выполненной!")
	showMenu()
}

func getKeyboardButton(buttonText, buttonCode string) tgbotapi.InlineKeyboardButton{
	return tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode)
} 

func showMenu() {
	msg := tgbotapi.NewMessage(gChatId, "Что ты сделал?")
	keyboard := tgbotapi.NewInlineKeyboardRow(
		getKeyboardButton(config.BUTTON_TOOTHS, config.BUTTON_CODE_TOOTHS),
		getKeyboardButton(config.BUTTON_SPORT, config.BUTTON_CODE_SPORT),
		getKeyboardButton(config.BUTTON_WALK, config.BUTTON_CODE_WALK),
		getKeyboardButton(config.BUTTON_BALANCE, config.BUTTON_CODE_BALANCE),
)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard)
	gBot.Send(msg)
}

func init() {
	var err error
	_ = os.Setenv(config.TOKEN_NAME_IN_OS, "6952704844:AAEx0mqHsLsv4SIUyHupC6koI5FyBubr5nQ")

	if gToken = os.Getenv(config.TOKEN_NAME_IN_OS); gToken == "" {
		panic(fmt.Errorf(`failed to load environment variable "%s"`, config.TOKEN_NAME_IN_OS))
	}
	gBot, err = tgbotapi.NewBotAPI(gToken)
	if err != nil {
		log.Panic(err)
	}

	gBot.Debug = true
}

func updateProcessing(update *tgbotapi.Update) {
	user, found := getUserFromUpdate(update)
	if !found {
		if user, found = storeUserFromUpdate(update); !found {
			gBot.Send(tgbotapi.NewMessage(gChatId, "Что-то пошло не так, попробуй еще раз"))
			return
		}
	}

	choiceCode := update.CallbackQuery.Data
	log.Printf("[%T] %s", time.Now(), choiceCode)

	switch choiceCode {
	case config.BUTTON_CODE_PRINT_INTRO:
		printIntro(update)
	case config.BUTTON_CODE_SKIP_INTRO:
		showMenu()
	case config.BUTTON_CODE_TOOTHS:
		printSystemMessageWithDelay(1, "Молодец! Умываться - очень важно! Ты заработал(а) 1 монету " + config.EMOJI_COIN)
		user.Points += 1
		showMenu()
	case config.BUTTON_CODE_SPORT:
		printSystemMessageWithDelay(1, "Молодец! Спорт - ! Ты заработал(а) 2 монеты " + config.EMOJI_COIN)
		user.Points += 2
		showMenu()
	case config.BUTTON_CODE_VITAMINS:
		printSystemMessageWithDelay(1, "Молодец! Пить витамины - важно! Ты заработал(а) 3 монеты " + config.EMOJI_COIN)
		user.Points += 3
		showMenu()
	case config.BUTTON_CODE_WALK:
		printSystemMessageWithDelay(1, "Молодец! Прогулка - важно! Ты заработал(а) 4 монеты " + config.EMOJI_COIN)
		user.Points += 4
		showMenu()
	case config.BUTTON_CODE_BALANCE:
		printSystemMessageWithDelay(1, fmt.Sprintf("Ты заработал(а) %d монет", user.Points))
		showMenu()
	default:
		log.Printf("Unknown button code: %s", &choiceCode)
		showMenu()
	}
}

func Start() {
	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(1)
	updateConfig.Timeout = config.UPDATE_CONFIG_TIMEOUT

	for update := range gBot.GetUpdatesChan(updateConfig) {
		if isCallbackQuery(&update) {
			updateProcessing(&update)
		} else if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			gChatId = update.Message.Chat.ID
			askToPrintIntro()
		}
	}
}