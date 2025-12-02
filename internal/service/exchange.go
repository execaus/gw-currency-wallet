package service

import (
	"context"
	"errors"
	gw_grpc "gw-currency-wallet/internal/pb/exchange"
	"gw-currency-wallet/pkg"
	"time"

	"go.uber.org/zap"
)

const (
	maxAttempts    = 3
	retryInterval  = 500 * time.Millisecond
	requestTimeout = 5 * time.Second
)

type ExchangeService struct {
	rateCache *pkg.Cacher[pkg.ExchangeRates]
	c         gw_grpc.ExchangeServiceClient
}

func (s *ExchangeService) GetRate(ctx context.Context, from, to pkg.Currency) (pkg.Rate, error) {
	rates, err := s.rateCache.GetData(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return -1, err
	}

	rateCh := make(chan float32, 1)
	timeoutContext, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	go func() {
		defer close(rateCh)

		var fromRate, toRate float32
		var fromOK, toOK bool

		for i := 0; i < maxAttempts; i++ {
			fromRate, fromOK = rates[from]
			toRate, toOK = rates[to]

			if fromOK && toOK {
				rateCh <- toRate / fromRate
				return
			}

			select {
			case <-timeoutContext.Done():
				return
			case <-time.After(retryInterval):
				newRates, err := s.rateCache.ForceSync(timeoutContext)
				if err != nil {
					zap.L().Error(err.Error())
				} else {
					rates = newRates
				}
			}
		}

		fromRate, fromOK = rates[from]
		toRate, toOK = rates[to]
		if fromOK && toOK {
			rateCh <- toRate / fromRate
			return
		}
	}()

	select {
	case <-timeoutContext.Done():
		zap.L().Error(ErrTimeout)
		return -1, errors.New(ErrTimeout)
	case rate, ok := <-rateCh:
		if !ok {
			return -1, errors.New(ErrGetRate)
		}
		return rate, nil
	}

}

func (s *ExchangeService) GetRates(ctx context.Context) (pkg.ExchangeRates, error) {
	rates, err := s.rateCache.GetData(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return rates, nil
}

func (s *ExchangeService) IsExistCurrency(ctx context.Context, currency pkg.Currency) (bool, error) {
	rates, err := s.rateCache.GetData(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return false, err
	}

	_, ok := rates[currency]
	return ok, nil
}

func NewExchangeService(ctx context.Context, client gw_grpc.ExchangeServiceClient) *ExchangeService {
	s := &ExchangeService{
		c: client,
	}

	cache, err := s.getRateCacher(ctx)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	s.rateCache = cache

	return s
}

func (s *ExchangeService) getRateCacher(ctx context.Context) (*pkg.Cacher[pkg.ExchangeRates], error) {
	return pkg.NewCacher[pkg.ExchangeRates](ctx, func(c context.Context) (pkg.ExchangeRates, error) {
		resp, err := s.c.GetExchangeRates(c, nil)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}

		return resp.Rates, nil
	}, pkg.DefaultCacherTTL)
}
