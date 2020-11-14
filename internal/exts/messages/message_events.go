package messages

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnMessage(s *discordgo.Session, message *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in OnMessage")
		}
	}()
	if message.Author.Bot {
		return
	}
	content := message.Content
	amazonRegexString := `(http|https):\/\/((?:[\w-_\.]+)amazon(?:\.[\w\-_]+)+)([\w\/\-\.,@?^=%&amp;~\+#]*[\w\-\@?^=%&amp;/~\+#])?`
	amazonRegex := regexp.MustCompile(amazonRegexString)
	urls := amazonRegex.FindAllString(message.Content, -1)
	if len(urls) == 0 {
		return
	}
	for _, url := range urls {
		if strings.Contains(url, "ref=") || strings.Contains(url, "?") {
			parts := strings.Split(url, "/")
			new := strings.Join(parts[:len(parts)-1], "/")
			content = strings.ReplaceAll(content, url, new)
		}
	}
	webhook, err := s.WebhookCreate(message.ChannelID, message.ID, "")
	if err != nil {
		return
	}
	defer s.WebhookDelete(webhook.ID)
	params := &discordgo.WebhookParams{
		Content:   content,
		Username:  message.Author.Username,
		AvatarURL: message.Author.AvatarURL(""),
	}
	_, err = s.WebhookExecute(webhook.ID, webhook.Token, true, params)
	if err != nil {
		return
	}
	s.ChannelMessageDelete(message.ChannelID, message.ID)
}
