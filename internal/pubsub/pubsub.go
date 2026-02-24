package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	ChannelPostCreated    = "post:created"
	ChannelCommentCreated = "comment:created:%s"
	ChannelLikeAdded      = "like:added:%s"
	ChannelLikeRemoved    = "like:removed:%s"
	ChannelUserFollowed   = "user:followed:%s"
	ChannelUserUnfollowed = "user:unfollowed:%s"
	ChannelUserUpdated    = "user:updated:%s"
)

type PubSub struct {
	client *redis.Client
}

func New(client *redis.Client) *PubSub {
	return &PubSub{client: client}
}

func (ps *PubSub) Publish(ctx context.Context, channel string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	return ps.client.Publish(ctx, channel, string(data)).Err()
}

func (ps *PubSub) Subscribe(ctx context.Context, channel string) <-chan string {
	msgCh := make(chan string, 10)
	sub := ps.client.Subscribe(ctx, channel)

	go func() {
		defer close(msgCh)
		defer sub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-sub.Channel():
				if !ok {
					return
				}
				select {
				case msgCh <- msg.Payload:
				default:
					log.Printf("pubsub: dropped message on channel %s (buffer full)", channel)
				}
			}
		}
	}()

	return msgCh
}
