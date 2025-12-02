package service

import (
	"context"
	gw_grpc "gw-currency-wallet/internal/pb/exchange"
	"gw-currency-wallet/pkg"

	"go.uber.org/zap"
)

type ExchangeService struct {
	rateCache *pkg.Cacher[pkg.ExchangeRates]
	c         gw_grpc.ExchangeServiceClient
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
