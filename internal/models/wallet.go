package models

// Currency код валюты, например "USD", "EUR", "JPY"
type Currency = string

type AccountWallets = map[Currency]float32
