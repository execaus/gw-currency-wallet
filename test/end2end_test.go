package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gw-currency-wallet/config"
	"gw-currency-wallet/internal/handler"
	gw_grpc "gw-currency-wallet/internal/pb/exchange"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/internal/service"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestE2E(t *testing.T) {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "1234",
			Name:     "gw-wallet",
		},
		Auth: config.AuthConfig{
			SecretKey: "your_jwt_secret_key",
		},
		ExchangeService: config.ExchangeService{
			Host: "localhost",
			Port: "8081",
		},
	}

	ctx := t.Context()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		zap.L().Fatal("error to open connect to database")
	}

	grpcConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.ExchangeService.Host, cfg.ExchangeService.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.L().Fatal("failed to connect to grpc server")
	}

	exchangeClient := gw_grpc.NewExchangeServiceClient(grpcConn)

	r := repository.NewRepository(pool)
	s := service.NewService(ctx, r, &cfg.Auth, exchangeClient)
	h := handler.NewHandler(s)

	router := h.Router()

	// 1. Регистрация пользователя
	username := fmt.Sprintf("user%d", rand.Intn(1000000))
	email := fmt.Sprintf("%s@example.com", username)
	password := "password123"
	regBody := map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	}
	rr := doPost(t, router, "/api/v1/register", regBody, "")
	assert.Equal(t, http.StatusCreated, rr.Code)
	var regResp map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &regResp)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", regResp["message"])

	// 2. Авторизация пользователя (login)
	loginBody := map[string]string{
		"username": username,
		"password": password,
	}
	rr = doPost(t, router, "/api/v1/login", loginBody, "")
	assert.Equal(t, http.StatusOK, rr.Code)
	var loginResp map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	token, ok := loginResp["token"]
	assert.True(t, ok)
	assert.NotEmpty(t, token)

	// 3. Пополнение счета (deposit)
	depositBody := map[string]interface{}{
		"amount":   100.0,
		"currency": "USD",
	}
	rr = doPost(t, router, "/api/v1/wallet/deposit", depositBody, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	var depositResp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &depositResp)
	assert.NoError(t, err)
	assert.Equal(t, "Account topped up successfully", depositResp["message"])

	// 4. Получение баланса (balance)
	rr = doGet(t, router, "/api/v1/balance", token)
	assert.Equal(t, http.StatusOK, rr.Code)
	var balanceResp map[string]map[string]float64
	err = json.Unmarshal(rr.Body.Bytes(), &balanceResp)
	assert.NoError(t, err)
	balance, ok := balanceResp["balance"]
	assert.True(t, ok)
	assert.NotEmpty(t, balance)
	assert.GreaterOrEqual(t, balance["USD"], 100.0)

	// 5. Вывод средств (withdraw)
	withdrawBody := map[string]interface{}{
		"amount":   50.0,
		"currency": "USD",
	}
	rr = doPost(t, router, "/api/v1/wallet/withdraw", withdrawBody, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	var withdrawResp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &withdrawResp)
	assert.NoError(t, err)
	assert.Equal(t, "Withdrawal successful", withdrawResp["message"])

	// 6. Получение курса валют (exchange rates)
	rr = doGet(t, router, "/api/v1/exchange/rates", token)
	assert.Equal(t, http.StatusOK, rr.Code)
	var ratesResp map[string]map[string]float64
	err = json.Unmarshal(rr.Body.Bytes(), &ratesResp)
	assert.NoError(t, err)
	rates, ok := ratesResp["rates"]
	assert.True(t, ok)
	assert.NotEmpty(t, rates)

	// 7. Обмен валют (exchange)
	exchangeAmount := 10.0
	exchangeBody := map[string]interface{}{
		"from_currency": "USD",
		"to_currency":   "EUR",
		"amount":        exchangeAmount,
	}
	rr = doPost(t, router, "/api/v1/exchange", exchangeBody, token)
	assert.Equal(t, http.StatusOK, rr.Code)
	var exchangeResp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &exchangeResp)
	assert.NoError(t, err)
	assert.Equal(t, "Exchange successful", exchangeResp["message"])
	assert.NotEmpty(t, exchangeResp["exchanged_amount"])
	expectedExchanged := exchangeAmount * rates["EUR"] / rates["USD"]
	assert.Equal(t, expectedExchanged, exchangeResp["exchanged_amount"])
}

func doPost(t *testing.T, handler http.Handler, url string, body interface{}, token string) *httptest.ResponseRecorder {
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func doGet(t *testing.T, handler http.Handler, url string, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}
