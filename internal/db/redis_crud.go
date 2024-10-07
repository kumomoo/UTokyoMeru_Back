package db

import (
	"context"
	"time"
)

var ctx = context.Background()

func SetVerificationCode(mail string, code string) (err error) {
	err = rdb.Set(ctx, "email_code:"+mail, code, 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetVerificationCode(mail string) (storedCode string, err error) {
	storedCode, err = rdb.Get(ctx, "email_code:"+mail).Result()
	if err != nil {
		return "", err
	}
	return
}
