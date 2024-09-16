package slackbot

import (
	"github.com/slack-go/slack"
	"golang.org/x/net/context"
)

const (
	BOT_CONTEXT      = "__BOT_CONTEXT__"
	MESSAGE_CONTEXT  = "__MESSAGE_CONTEXT__"
	REACTION_CONTEXT = "__REACTION_CONTEXT__"
	REACTION_EVENT   = "__REACTION_EVENT__"
	BOT_DEBUG        = "__BOT_DEBUG__"
	CHANNEL_CONTEXT  = "__CHANNEL_CONTEXT__"
	CHANNEL_JOINED   = "__CHANNEL_JOINED__"
	CHANNEL_LEFT     = "__CHANNEL_LEFT__"
)

func BotFromContext(ctx context.Context) *Bot {
	if result, ok := ctx.Value(BOT_CONTEXT).(*Bot); ok {
		return result
	}
	return nil
}

// AddBotToContext sets the bot reference in context and returns the newly derived context
func AddBotToContext(ctx context.Context, bot *Bot) context.Context {
	return context.WithValue(ctx, BOT_CONTEXT, bot)
}

func MessageFromContext(ctx context.Context) *slack.MessageEvent {
	if result, ok := ctx.Value(MESSAGE_CONTEXT).(*slack.MessageEvent); ok {
		return result
	}
	return nil
}

// AddMessageToContext sets the Slack message event reference in context and returns the newly derived context
func AddMessageToContext(ctx context.Context, msg *slack.MessageEvent) context.Context {
	return context.WithValue(ctx, MESSAGE_CONTEXT, msg)
}

// AddReactionAddedToContext
func AddReactionAddedToContext(ctx context.Context, react *slack.ReactionAddedEvent) context.Context {
	nctx := context.WithValue(ctx, REACTION_CONTEXT, "Added")
	return context.WithValue(nctx, REACTION_EVENT, react)
}

// AddReactionAddedToContext
func AddReactionRemovedToContext(ctx context.Context, react *slack.ReactionRemovedEvent) context.Context {
	nctx := context.WithValue(ctx, REACTION_CONTEXT, "Removed")
	return context.WithValue(nctx, REACTION_EVENT, react)
}

func ReactionTypeFromContext(ctx context.Context) string {
	if result, ok := ctx.Value(REACTION_CONTEXT).(string); ok {
		return result
	}
	return ""
}
func ReactionAddedFromContext(ctx context.Context) *slack.ReactionAddedEvent {
	if result, ok := ctx.Value(REACTION_EVENT).(*slack.ReactionAddedEvent); ok {
		return result
	}
	return nil
}

func ReactionRemovedFromContext(ctx context.Context) *slack.ReactionRemovedEvent {
	if result, ok := ctx.Value(REACTION_EVENT).(*slack.ReactionRemovedEvent); ok {
		return result
	}
	return nil
}

func AddChannelJoinToContext(ctx context.Context, channel *slack.ChannelJoinedEvent) context.Context {
	nctx := context.WithValue(ctx, CHANNEL_CONTEXT, "Joined")
	return context.WithValue(nctx, CHANNEL_JOINED, channel)
}

func AddChannelLeftToContext(ctx context.Context, channel *slack.ChannelLeftEvent) context.Context {
	nctx := context.WithValue(ctx, CHANNEL_CONTEXT, "Left")
	return context.WithValue(nctx, CHANNEL_LEFT, channel)
}

func ChannelTypeFromContext(ctx context.Context) string {
	if result, ok := ctx.Value(CHANNEL_CONTEXT).(string); ok {
		return result
	}
	return ""

}
func ChannelJoinedFromContext(ctx context.Context) *slack.ChannelJoinedEvent {
	if result, ok := ctx.Value(CHANNEL_JOINED).(*slack.ChannelJoinedEvent); ok {
		return result
	}
	return nil
}

func ChannelLeftFromContext(ctx context.Context) *slack.ChannelLeftEvent {
	if result, ok := ctx.Value(CHANNEL_LEFT).(*slack.ChannelLeftEvent); ok {
		return result
	}
	return nil
}

func SetDebug(ctx context.Context) context.Context {
	return context.WithValue(ctx, BOT_DEBUG, true)
}

func IsDebug(ctx context.Context) bool {
	if result, ok := ctx.Value(BOT_DEBUG).(bool); ok {
		return result
	}
	return false
}
