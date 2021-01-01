package discord

import "fmt"

// Credentials stores the credentials necessary for authentication with the Discord API.
type Credentials struct {
	Token string
}

// ChannelMessageSendError is an implementation of the error interface used to signify an error
// sending a message to a channel.
type ChannelMessageSendError struct {
	message string
}

// NewChannelMessageSendError returns a new instance of ChannelMessageSendError.
func NewChannelMessageSendError(e error) *ChannelMessageSendError {
	return &ChannelMessageSendError{
		message: fmt.Sprintf("error sending message: %s", e),
	}
}

// Error implements the error interface for ChannelMessageSendError type.
func (e *ChannelMessageSendError) Error() string {
	return e.message
}
