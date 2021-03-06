package main

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func handleDiscord(params map[string]string, messages []string) error {
	session, err := getDiscordSession(params["discord_bot_token"])
	if err != nil {
		return err
	}

	err = session.Open()
	if err != nil {
		return err
	}
	defer session.Close() // nolint

	// TODO: handle more than 100 servers
	guilds, err := session.UserGuilds(100, "", "")
	if err != nil {
		return fmt.Errorf("problem retrieving discord servers: %v", err)
	}
	log.Debugf("guild server count: %d", len(guilds))

	log.Debugf("messages to send: %s", messages)

	for _, g := range guilds {
		channels, err := session.GuildChannels(g.ID)
		if err != nil {
			log.Debugf("channel retrieval error: %v", err)
		}

		for _, ch := range channels {
			if ch.Type == discordgo.ChannelTypeGuildText {
				log.Debugf("sending messages to Guild name: %s - Channel: %s", g.Name, ch.Name)
				for _, m := range messages {
					_, err := session.ChannelMessageSend(ch.ID, m)
					if err != nil {
						log.Debugf("channel send error: %v", err)
					}
				}
			}
		}
	}
	return nil
}

// getDiscordSession returns a discord session given proper credentials.
func getDiscordSession(token string) (*discordgo.Session, error) {
	if token == "" {
		return nil, errors.New("discord bot token required")
	}

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	return sess, nil
}
