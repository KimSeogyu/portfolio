package redisutils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, addrs ...string) (redis.UniversalClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:          addrs,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		DialTimeout:    1 * time.Second,
		PoolSize:       50,
		MinIdleConns:   10,
		MaxRetries:     5,
		ReadOnly:       false, // 읽기 전용 모드 비활성화
		RouteRandomly:  true,  // 클러스터 노드에 랜덤하게 요청 분산
		RouteByLatency: true,  // 지연 시간이 낮은 노드 선호
	})

	status := client.Ping(ctx)
	if status.Err() != nil {
		return nil, fmt.Errorf("failed to connect to Redis cluster: %w", status.Err())
	}

	return client, nil
}
