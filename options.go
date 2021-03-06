package twtxt

import "time"

const (
	// DefaultData is the default data directory for storage
	DefaultData = "./data"

	// DefaultStore is the default data store used for accounts, sessions, etc
	DefaultStore = "bitcask://twtxt.db"

	// DefaultBaseURL is the default Base URL for the app used to construct feed URLs
	DefaultBaseURL = "http://0.0.0.0:8000"

	// DefaultName is the default instance name
	DefaultName = "twtxt.net"

	// DefaultRegister is the default user registration flag
	DefaultRegister = false

	// DefaultCookieSecret is the server's default cookie secret
	DefaultCookieSecret = "PLEASE_CHANGE_ME!!!"

	// DefaultSessionExpiry is the server's default session expiry time
	DefaultSessionExpiry = 24 * time.Hour
)

func NewConfig() *Config {
	return &Config{
		Data:    DefaultData,
		Store:   DefaultStore,
		BaseURL: DefaultBaseURL,
	}
}

// Option is a function that takes a config struct and modifies it
type Option func(*Config) error

// WithData sets the data directory to use for storage
func WithData(data string) Option {
	return func(cfg *Config) error {
		cfg.Data = data
		return nil
	}
}

// WithStore sets the store to use for accounts, sessions, etc.
func WithStore(store string) Option {
	return func(cfg *Config) error {
		cfg.Store = store
		return nil
	}
}

// WithBaseURL sets the Base URL used for constructing feed URLs
func WithBaseURL(baseURL string) Option {
	return func(cfg *Config) error {
		cfg.BaseURL = baseURL
		return nil
	}
}

// WithName sets the instance's name
func WithName(name string) Option {
	return func(cfg *Config) error {
		cfg.Name = name
		return nil
	}
}

// WithRegister sets the user registration flag
func WithRegister(register bool) Option {
	return func(cfg *Config) error {
		cfg.Register = register
		return nil
	}
}

// WithCookieSecret sets the server's cookie secret
func WithCookieSecret(secret string) Option {
	return func(cfg *Config) error {
		cfg.CookieSecret = secret
		return nil
	}
}

// WithSessionExpiry sets the server's session expiry time
func WithSessionExpiry(expiry time.Duration) Option {
	return func(cfg *Config) error {
		cfg.SessionExpiry = expiry
		return nil
	}
}
