package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5355283427:AAEV88bU1yf6qgzAFOWjboxMdbhCjUPGdo0") //токен телеграм бота

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Авторизация на боте под названием: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	flagOfWaitingProblems := false

	log.Println(getTextInfoFromFileTxt())

	for update := range updates {
		if update.Message == nil { // Если не получили сообщение - продолжаем итерацию до возникновения
			continue
		}
		//если полученное сообщение  является фото, аудио или голосовым сообщением - отправляем пользователю сообщение о невозможности прочитать такие файлы
		if update.Message.Photo != nil {
			fmt.Println("Отправлена фотография!")
			fmt.Println(update.Message.Caption)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, но я не умею читать фотографии! Наверно, она очень красивая :("))
		} else if update.Message.Audio != nil {
			fmt.Println("Отправлен аудиофайл!")
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, но я не умею понимать музыку! Хотелось бы узнать, какой он на слух - Чайковский... :("))
		} else if update.Message.Voice != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, но я не умею понимать голос! Но он у вас очень красивый!"))
		} else {
			commands(update, bot, &flagOfWaitingProblems) //если сообщение является текстом - запускаем функцию обработки сообщения с флагом ожидания (для отправки обращений)
		}
	}

}
