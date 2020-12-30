package discord

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Session is an alias for discordgo.Session with some additional receiver functions for convenience.
type Session discordgo.Session

// GetSession returns a discord session given proper credentials.
func GetSession(c *Credentials) (*Session, error) {
	if c.Token == "" {
		return nil, errors.New("discord bot token required")
	}

	sess, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %v", err)
	}

	return (*Session)(sess), nil
}

// Send is a convenience function wrapping the discord session functionality.
func (s *Session) Send(messages []string) []error {
	var discordErrors []error

	ds := (*discordgo.Session)(s)
	err := ds.Open()
	if err != nil {
		discordErrors = append(discordErrors, err)
		return discordErrors
	}
	defer ds.Close() // nolint

	// TODO: handle more than 100 servers
	guilds, err := ds.UserGuilds(100, "", "")
	if err != nil {
		discordErrors = append(discordErrors, err)
		return discordErrors
	}

	for _, g := range guilds {
		channels, err := ds.GuildChannels(g.ID)
		if err != nil {
			discordErrors = append(discordErrors, err)
			continue
		}
		for _, ch := range channels {
			if ch.Type == discordgo.ChannelTypeGuildText {

				// log.Debugf("Guild name: %s - Channel: %s", g.Name, ch.Name)
				// TODO: remove hardcoded messages[0]
				msg, err := ds.ChannelMessageSend(ch.ID, messages[0])
				if err != nil {
					discordErrors = append(discordErrors, fmt.Errorf("message send error: %s, msg: %+v", err, msg))
					continue
				}
				// log.Debugf("msg: %+v", msg)
			}
		}
	}

	return discordErrors
}
