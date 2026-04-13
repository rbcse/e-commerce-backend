package otprepository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPData struct {
	Identifier string
	OTP string
	Attempts int
}

type OTPRepository interface {
	SaveOTP(identifier, otp string) error
	GetOTP(identifier string) (*OTPData)
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
	fmt.Println("Otp saved to redis");
	return err

}

func (or *otpRepository) GetOTP(identifier string) (*OTPData) {

	ctx := context.Background()
	key := "otp:" + identifier
	result , err := or.client.HGetAll(ctx,key).Result()

	if err != nil {
		return nil
	}

	if len(result) == 0 {
		return nil
	}

	attempts , _ := strconv.Atoi(result["attempts"])
	return &OTPData{
		Identifier: result["identifier"],
		OTP: result["otp"],
		Attempts: attempts,
	}

}