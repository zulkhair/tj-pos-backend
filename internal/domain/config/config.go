package configdomain

import "sync"

type Config struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type ConfigCache struct {
	sync.RWMutex
	DataMap map[string]*Config
}

const LOGIN_URL = "LOGIN_URL"
const FORBIDDEN_URL = "FORBIDDEN_URL"
const UNAUTHORIZED_URL = "UNAUTHORIZED_URL"
const SESSION_TIMEOUT_MINUTE = "SESSION_TIMEOUT_MINUTE"
