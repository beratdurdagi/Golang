package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingService struct {
	next PriceFetcher
}

func NewLogging(next PriceFetcher) *loggingService {
	return &loggingService{
		next: next,
	}
}

func (s *loggingService) FetchPrice(ctx context.Context, ticker string) (price float64, err error) {

	defer func(begin time.Time) {
		reqID := ctx.Value("requestID")
		logrus.WithFields(logrus.Fields{
			"requestID": reqID,
			"took":      time.Since(begin),
			"err":       err,
			"price":     price,
		}).Info("fetchPrice")

	}(time.Now())

	return s.next.FetchPrice(ctx, ticker)
}
