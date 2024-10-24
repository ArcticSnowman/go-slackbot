// Package slackbot hopes to ease development of Slack bots by adding helpful
// methods and a mux-router style interface to the github.com/slack-go/slack package.
//
// Incoming Slack RTM events are mapped to a handler in the following form:
//
//	bot.Hear("(?i)how are you(.*)").MessageHandler(HowAreYouHandler)
//
// The package adds Reply and ReplyWithAttachments methods:
//
//	func HowAreYouHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
//		bot.Reply(evt, "A bit tired. You get it? A bit?", slackbot.WithTyping)
//	}
//
//	func HowAreYouAttachmentsHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
//		txt := "Beep Beep Boop is a ridiculously simple hosting platform for your Slackbots."
//		attachment := slack.Attachment{
//			Pretext:   "We bring bots to life. :sunglasses: :thumbsup:",
//			Title:     "Host, deploy and share your bot in seconds.",
//			TitleLink: "https://GrantStreetGroup.com/",
//			Text:      txt,
//			Fallback:  txt,
//			ImageURL:  "https://storage.googleapis.com/GrantStreetGroup/_assets/bot-1.22f6fb.png",
//			Color:     "#7CD197",
//		}
//
//		attachments := []slack.Attachment{attachment}
//		bot.ReplyWithAttachments(evt, attachments, slackbot.WithTyping)
//	}
//
// The slackbot package exposes  github.com/slack-go/slack RTM and Client objects
// enabling a consumer to interact with the lower level package directly:
//
//	func HowAreYouHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
//		bot.RTM.NewOutgoingMessage("Hello", "#random")
//	}
//
// Project home and samples: https://github.com/GrantStreetGroup/go-slackbot
package slackbot

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

const (
	WithTyping    bool = true
	WithoutTyping bool = false

	maxTypingSleep time.Duration = time.Millisecond * 2000
)

// New constructs a new Bot using the slackToken to authorize against the Slack service.
func New(slackToken string, options ...slack.Option) *Bot {
	b := &Bot{Client: slack.New(slackToken, options...)}
	b.EventMode = "RTM"
	return b
}

type Bot struct {
	SimpleRouter
	// Routes to be matched, in order.
	routes []*Route
	// Slack UserID of the bot UserID
	botUserID string
	// Slack API
	Client    *slack.Client
	RTM       *slack.RTM
	Socket    *socketmode.Client
	Debug     bool
	EventMode string
}

// Run listeners for incoming slack events via RTM or Socketmode, matching them to an appropriate handler.
func (b *Bot) Run() {

	if b.EventMode == "RTM" {
		b.RunRTM()
	} else {
		b.RunSocketMode()
	}
}

// Run listeners for incoming slack RTM events, matching them to an appropriate handler.
func (b *Bot) RunRTM() {
	b.RTM = b.Client.NewRTM()
	go b.RTM.ManageConnection()
	for {
		select {
		case msg := <-b.RTM.IncomingEvents:
			ctx := context.Background()
			ctx = AddBotToContext(ctx, b)
			if b.Debug {
				ctx = SetDebug(ctx)
			}
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Printf("Connected: %#v\n", ev.Info.User)
				b.SetBotID(ev.Info.User.ID)
			case *slack.MessageEvent:
				// ignore messages from the current user, the bot user
				if b.botUserID == ev.User {
					continue
				}
				// ignore message_replies subtypes as we will get the full message event.
				if ev.SubType == slack.MsgSubTypeMessageReplied {
					continue
				}

				// Ignore thread broadcast messages
				if ev.SubType == slack.MsgSubTypeThreadBroadcast {
					continue
				}

				// Ignore bot_messages
				if ev.SubType == slack.MsgSubTypeBotMessage {
					continue
				}
				ctx = AddMessageToContext(ctx, ev)
				var match RouteMatch
				if matched, ctx := b.Match(ctx, &match); matched {
					match.Handler(ctx)
				}

			case *slack.ReactionAddedEvent:
				// Handle reaction events
				if b.botUserID == ev.User {
					continue
				}

				ctx = AddReactionAddedToContext(ctx, ev)
				var match RouteMatch
				if matched, ctx := b.Match(ctx, &match); matched {
					match.Handler(ctx)
				}
			case *slack.ReactionRemovedEvent:
				// Handle reaction events
				if b.botUserID == ev.User {
					continue
				}

				ctx = AddReactionRemovedToContext(ctx, ev)
				var match RouteMatch
				if matched, ctx := b.Match(ctx, &match); matched {
					match.Handler(ctx)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break

			default:
				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}

// Reply replies to a message event with a simple message.
func (b *Bot) Reply(evt *slack.MessageEvent, msg string, typing bool) {
	if b.EventMode == "RTM" {
		if typing {
			b.Type(evt, msg)
		}
		b.RTM.SendMessage(b.RTM.NewOutgoingMessage(msg, evt.Channel))
	} else {
		b.Client.PostMessage(evt.Channel, slack.MsgOptionText(msg, true), slack.MsgOptionAsUser(true))
	}
}

// ReplyWithAttachments replys to a message event with a Slack Attachments message.
func (b *Bot) ReplyWithAttachments(evt *slack.MessageEvent, attachments []slack.Attachment, typing bool) {
	//	params := slack.PostMessageParameters{AsUser: true}
	//	params.Attachments = attachments

	b.Client.PostMessage(evt.Msg.Channel, slack.MsgOptionAttachments(attachments...), slack.MsgOptionAsUser(true))
}

// Type sends a typing message and simulates delay (max 2000ms) based on message size.
func (b *Bot) Type(evt *slack.MessageEvent, msg interface{}) {
	msgLen := msgLen(msg)

	sleepDuration := time.Minute * time.Duration(msgLen) / 3000
	if sleepDuration > maxTypingSleep {
		sleepDuration = maxTypingSleep
	}

	b.RTM.SendMessage(b.RTM.NewTypingMessage(evt.Channel))
	time.Sleep(sleepDuration)
}

// Fetch the botUserID.
func (b *Bot) BotUserID() string {
	return b.botUserID
}

func (b *Bot) SetBotID(ID string) {
	b.botUserID = ID
}

// msgLen gets lenght of message and attachment messages. Unsupported types return 0.
func msgLen(msg interface{}) (msgLen int) {
	switch m := msg.(type) {
	case string:
		msgLen = len(m)
	case []slack.Attachment:
		msgLen = len(fmt.Sprintf("%#v", m))
	}
	return
}
