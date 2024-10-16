package slackbot

import (
	"golang.org/x/net/context"

	"github.com/slack-go/slack"
)

type MessageType string

type Reactions string

const (
	DirectMessage MessageType = "direct_message"
	DirectMention MessageType = "direct_mention"
	Mention       MessageType = "mention"
	Ambient       MessageType = "ambient"
)

type Handler func(context.Context)
type MessageHandler func(ctx context.Context, bot *Bot, msg *slack.MessageEvent)
type ReactionHandler func(ctx context.Context, bot *Bot, added *slack.ReactionAddedEvent, removed *slack.ReactionRemovedEvent)

type Preprocessor func(context.Context) context.Context

// Matcher type for matching message routes
type Matcher interface {
	Match(context.Context) (bool, context.Context)
	SetBotID(botID string)
}
