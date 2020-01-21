package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v7"
)

func consume(ctx context.Context, cfg *Config, client redis.UniversalClient) error {

	xReadArgs := &redis.XReadArgs{
		Streams: []string{cfg.Stream, "$"},
		Block:   time.Duration(cfg.Block) * time.Second,
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		resp, err := client.XRead(xReadArgs).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return fmt.Errorf("XRead failed: %s", err.Error())
		}
		stream := resp[0]

		for _, m := range stream.Messages {
			b, err := json.Marshal(m)
			if err != nil {
				return fmt.Errorf("json.Marshal failed: %s", err.Error())
			}
			fmt.Println(string(b))
			// we should use the $ ID only for the first call to XREAD. Subsequent
			// ID should be the last reported item in the stream, otherwise we
			// could miss all the entries that are added in between.
			xReadArgs.Streams[1] = m.ID
		}
	}
}
