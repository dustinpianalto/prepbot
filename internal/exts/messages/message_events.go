package messages

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func CleanAmazonURLs(s *discordgo.Session, message *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in CleanAmazonURLs")
		}
	}()
	if message.Author.Bot {
		return
	}
	content := message.Content
	amazonRegexString := `(http|https):\/\/((?:[\w-_\.]*)amazon(?:\.[\w\-_]+)+)([\w\/\-\.,@?^=%&amp;~\+#]*[\w\-\@?^=%&amp;/~\+#])?`
	amazonRegex := regexp.MustCompile(amazonRegexString)
	urls := amazonRegex.FindAllString(message.Content, -1)
	if len(urls) != 0 {

		for _, url := range urls {
			if strings.Contains(url, "ref=") || strings.Contains(url, "?") {
				parts := strings.Split(url, "/")
				new := strings.Join(parts[:len(parts)-1], "/")
				if strings.Contains(new, "ref=") {
					parts = strings.Split(new, "ref=")
					new = parts[0]
				}
				content = strings.ReplaceAll(content, url, new)
			}
		}
		_, err := sendWebhook(s, message, message.ChannelID, content)
		if err == nil {
			s.ChannelMessageDelete(message.ChannelID, message.ID)
		}
	}
	message.Content = content
	moveNewsLinks(s, message)
}

func moveNewsLinks(s *discordgo.Session, message *discordgo.MessageCreate) {
	linkChannel := "784855106239332352"
	chatChannel := "777679317878177804"
	if message.ChannelID == chatChannel {
		_, err := sendWebhook(s, message, linkChannel, message.Content)
		if err != nil {
			log.Println(err)
		}
	}
}

func sendWebhook(s *discordgo.Session, message *discordgo.MessageCreate, channelID, content string) (string, error) {
	webhook, err := s.WebhookCreate(channelID, message.ID, "")
	if err != nil {
		return "", err
	}
	defer s.WebhookDelete(webhook.ID)
	var name string
	if message.Member != nil && message.Member.Nick != "" {
		name = message.Member.Nick
	} else {
		name = message.Author.Username
	}
	params := &discordgo.WebhookParams{
		Content:   content,
		Username:  name,
		AvatarURL: message.Author.AvatarURL(""),
	}
	w, err := s.WebhookExecute(webhook.ID, webhook.Token, true, params)
	if err != nil {
		return "", err
	}
	return w.ID, nil
}
