package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Ссылка на документацию", "https://go-telegram-bot-api.dev/examples/inline-keyboard.html")),
) //inline кнопка для перехода к документации на GitHub

const AdminID int64 = 313944272  //id-chat администратора
const TextFile = "text/info.txt" //файл с информационным блоком для отправки
type File struct {               //структура для отправки сообщения с файлом
	path string
	name string
}

type Info struct { //структура для хранения сообщений команд
	start                  string
	help                   string
	inform                 string
	getID                  string
	noneCommand            string
	wrongFormat            string
	hasWhiteSpaceInCommand string
	textFAQ                string
}

func photoInitialise() []string { //инициализация фото для отправки
	photoPath := [...]string{
		"1.png", //get_id рисунок 1
		"2.png", //wAP кнопка
		"3.png", //форма регистрации

	}

	return photoPath[:]
}

func display(file []*File) { //команда для дебагга фотографий
	for i := range photoInitialise() {
		log.Println(file[i].path + file[i].name)
	}
	log.Println("Количество адресов фотографий: " + strconv.Itoa(len(photoInitialise())))
}

func pickTheImage(id int, file []*File) tgbotapi.FileBytes { //возвращает фотографию в виде пакета Byte данных для отправки в телеграмм
	var photoBytes []byte
	var err error
	for i := range file {
		if i == id {
			photoBytes, err = ioutil.ReadFile(file[i].path + file[i].name)
			if err != nil {
				panic(err)
			}
		}
	}

	return tgbotapi.FileBytes{
		Name:  "Picture",
		Bytes: photoBytes,
	}
}

func getTextInfoFromFileTxt() string { //функция для обработки сообщения информационного блока
	file, err := os.Open(TextFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadFile(TextFile)
	if err != nil {
		panic(err)
	}
	return string(fileContent)
}

func commands(update tgbotapi.Update, bot *tgbotapi.BotAPI, flagOfWaitingProblems *bool) { //обработчик сообщений

	var file []*File
	for i := range photoInitialise() {
		file = append(file, &File{"photo/", photoInitialise()[i]})
	}

	log.Println(len(file))

	info := Info{
		"Приветствую, " + update.Message.Chat.FirstName + "! Ты находишься во вспомогательном боте проекта АСПИД - ASPID_INFO. Прежде всего хочу поблагодарить тебя за использование нашего продукта. Я очень рад, что ты приобрел АСПИД для автоматизированного управления своим домом. Будь уверен, он в надежных и нежных роборуках! " +
			"Для просмотра команд введи /help. Если хочешь сразу узнать, как настроить систему АСПИД для работы - введи /info",
		"Ниже приведен список команд, используемых в этом боте:\n/start - основная команда запуска бота\n/help - команда для отображения всех команд\n/info - " +
			"памятка пользователя, информация по подключению к устройству\n/get_id - команда для отображения вашего UserID\n" +
			"/problems - команда для отправки обращения в службу технической поддержки\n/back - команда отмены отправки сообщения в техническую поддержку (вводится только после /problems)\n",
		getTextInfoFromFileTxt(),
		"Ваш идентификационный номер: " + strconv.FormatInt(update.Message.Chat.ID, 10) + ". Если вы не знаете, что с ним делать - воспользуйтесь командой /info",
		"Такая команда не поддерживается нашим ботом. Пожалуйста, воспользуйтесь командой /help для отображения списка действующих команд\n",
		"Простите, но я умею работать только в формате команд. Пожалуйста, введите сообщение в формате /команда. Для отображения списка команд воспользуйтесь /help",
		"После команды не должно быть пробелов. Повторите попытку снова!", "Пожалуйста, опишите проблему в следующем сообщении: ",
	}
	command := strings.Split(update.Message.Text, " ")

	if *flagOfWaitingProblems && command[0] != "/back" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо! Ваше обращение отправлено, в скором времени вам придет ответ в этот чат, так что советуем не покидать его!"))
		fmt.Println("Пользователь с ником: " + update.Message.Chat.UserName + " и UserID: " + strconv.FormatInt(update.Message.Chat.ID, 10) + " прислал сообщение о проблеме: ")
		fmt.Println(update.Message.Text)

		//код для обращений
		msg := tgbotapi.NewMessage(AdminID, "*Внимание, обращение!*\n")
		msg.ParseMode = "markdown"
		bot.Send(msg)
		msg = tgbotapi.NewMessage(AdminID, "Сообщение №"+"*"+strconv.Itoa(update.Message.MessageID)+"*"+": Пользователь с ником: "+"*"+update.Message.Chat.UserName+"*"+" и UserID: "+"*"+strconv.FormatInt(update.Message.Chat.ID, 10)+"*"+" прислал сообщение о проблеме:\n\n"+"_"+update.Message.Text+"_")
		msg.ParseMode = "markdown"
		bot.Send(msg)
		*flagOfWaitingProblems = false
	} else {

		if command[0][0:1] != "/" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, info.wrongFormat))
		} else if command[0][0:1] == "/" {
			if len(command) == 1 {
				switch command[0] {
				case "/start":
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, info.start))
				case "/help":
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, info.help))
				case "/info":
					informSlice := strings.Split(info.inform, "@")
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, informSlice[0]))
					bot.Send(tgbotapi.NewPhotoUpload(update.Message.Chat.ID, pickTheImage(0, file)))
					time.Sleep(2 * time.Second)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, informSlice[1]))
					bot.Send(tgbotapi.NewPhotoUpload(update.Message.Chat.ID, pickTheImage(1, file)))
					time.Sleep(2 * time.Second)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, informSlice[2]))
					bot.Send(tgbotapi.NewPhotoUpload(update.Message.Chat.ID, pickTheImage(0, file)))
					time.Sleep(2 * time.Second)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, informSlice[3]))
					bot.Send(tgbotapi.NewPhotoUpload(update.Message.Chat.ID, pickTheImage(1, file)))
					time.Sleep(2 * time.Second)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, informSlice[4]))
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, informSlice[5])
					msg.ReplyMarkup = numericKeyboard
					bot.Send(msg)
				case "/get_id":
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, info.getID))
				case "/photo":
					display(file)
					bot.Send(tgbotapi.NewPhotoUpload(update.Message.Chat.ID, pickTheImage(0, file)))
					bot.Send(tgbotapi.NewPhotoUpload(update.Message.Chat.ID, pickTheImage(1, file)))
				case "/problems":
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, info.textFAQ))
					*flagOfWaitingProblems = true
				case "/back":
					if *flagOfWaitingProblems != true {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда /back используется только для отмены отправки обращения"))
					} else {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Обращение отменено. "))
						*flagOfWaitingProblems = false
					}
				case "/reply":
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда /reply является системной, пожалуйста, не используйте ее во избежания проблем."))

				default:
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, info.noneCommand))
				}
			} else if command[0] == "/reply" { // формат команды /reply id_chat id_message "Текст"
				if update.Message.Chat.ID != AdminID {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, друг, не шали с командой /reply, ты делаешь моему роботизированному сердцу неисправимо больно!"))
				} else {
					if len(command) > 4 {
						convertingStringChatID, err := strconv.ParseInt(command[1], 10, 64)
						if err != nil {
							panic(err)
						}

						convertingStringMessageID, err := strconv.Atoi(command[2])
						if err != nil {
							panic(err)
						}

						message := "Вы получили от разработчика следующий ответ: \n\n"

						for i := 3; i < len(command); i++ {
							message += command[i] + " "
						}

						msg := tgbotapi.NewMessage(convertingStringChatID, message)
						msg.ReplyToMessageID = convertingStringMessageID

						bot.Send(msg)
					} else {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неправильный формат команды!"))
					}
				}

			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, info.hasWhiteSpaceInCommand))
			}
		}
	}

}
