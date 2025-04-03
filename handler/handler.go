package handler

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/StupidYoshiaki/FastTakeshitake/downloader"
	"github.com/bwmarrin/discordgo"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	found := make(chan string, 1)
	var wg sync.WaitGroup

	for key := range downloader.S3FilePaths {
		k := key
		wg.Add(1)
		go func(k string) {
			name := strings.ReplaceAll(k, ".png", "")
			defer wg.Done()
			if strings.Contains(message, name) {
				select {
				case found <- k:
				default:
				}
			}
		}(k)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case key, ok := <-found:
		if ok {
			log.Printf("Matched key: %s", key)

			fileName := key
			filePath := downloader.S3FilePaths[fileName]

			file, err := os.Open(filePath)
			if err != nil {
				log.Printf("Failed to open file: %v", err)
				return
			}
			defer file.Close()
			msg := &discordgo.MessageSend{
				Files: []*discordgo.File{
					{
						Name:   fileName,
						Reader: file,
					},
				},
			}
			_, err = s.ChannelMessageSendComplex(m.ChannelID, msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}
			log.Printf("Sent file: %s", fileName)
		}
	case <-done:
		log.Println("No match found")
	}
}
