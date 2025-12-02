package pkg

// Currency код валюты, например "USD", "EUR", "JPY"
type Currency = string

// Rate курс валюты относительно базовой валюты
type Rate = float32

// ExchangeRates — карта валютных курсов
type ExchangeRates = map[Currency]Rate

// AccountWallets - состояние баланса кошельков
type AccountWallets = map[Currency]float32
