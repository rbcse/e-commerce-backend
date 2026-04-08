package otprepository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPRepository interface {
	SaveOTP(identifier, otp string) error
}

type otpRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewOTPRepository(client *redis.Client, ttl time.Duration) OTPRepository {
	return &otpRepository{
		client: client,
		ttl:    ttl,
	}
}

func (or *otpRepository) SaveOTP(identifier, otp string) error {

	ctx := context.Background()
	key := "otp:" + identifier

	pipe := or.client.Pipeline()
	pipe.HSet(ctx, key, map[string]interface{}{
		"identifier": identifier,
		"otp":        otp,
		"attempts":   3,
	})

	pipe.Expire(ctx, key, or.ttl)
	_, err := pipe.Exec(ctx)
	return err

}
