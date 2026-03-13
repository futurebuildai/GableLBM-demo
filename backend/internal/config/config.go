package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWKSURL     string
	AuthIssuer  string

	// Run Payments Gateway
	RunPaymentsAPIKey      string
	RunPaymentsPublicKey   string
	RunPaymentsBaseURL     string
	RunPaymentsEnvironment string // "sandbox" or "production"

	// Avalara Sales Tax
	AvalaraAccountID   string
	AvalaraLicenseKey  string
	AvalaraEnvironment string // "sandbox" or "production"
	AvalaraCompanyCode string

	// Google Maps
	GoogleMapsAPIKey string

	// Twilio SMS
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string

	// Anthropic (Claude AI — PIM Content Generation)
	AnthropicAPIKey string
	AnthropicModel  string

	// Stability AI (PIM Image Generation)
	StabilityAPIKey string
}

func Load() *Config {
	_ = godotenv.Load() // Load .env if it exists, ignore if not

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://gable_user:gable_password@localhost:5434/gable_db?sslmode=disable"),
		JWKSURL:     getEnv("JWKS_URL", ""),
		AuthIssuer:  getEnv("AUTH_ISSUER", ""),

		// Run Payments — defaults to sandbox mode
		RunPaymentsAPIKey:      getEnv("RUN_PAYMENTS_API_KEY", ""),
		RunPaymentsPublicKey:   getEnv("RUN_PAYMENTS_PUBLIC_KEY", ""),
		RunPaymentsBaseURL:     getEnv("RUN_PAYMENTS_BASE_URL", ""),
		RunPaymentsEnvironment: getEnv("RUN_PAYMENTS_ENV", "sandbox"),

		// Avalara Sales Tax — defaults to sandbox mode
		AvalaraAccountID:   getEnv("AVALARA_ACCOUNT_ID", ""),
		AvalaraLicenseKey:  getEnv("AVALARA_LICENSE_KEY", ""),
		AvalaraEnvironment: getEnv("AVALARA_ENV", "sandbox"),
		AvalaraCompanyCode: getEnv("AVALARA_COMPANY_CODE", ""),

		// Google Maps
		GoogleMapsAPIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),

		// Twilio SMS
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: getEnv("TWILIO_FROM_NUMBER", ""),

		// Anthropic (Claude AI — PIM Content Generation)
		AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
		AnthropicModel:  getEnv("ANTHROPIC_MODEL", "claude-sonnet-4-20250514"),

		// Stability AI (PIM Image Generation)
		StabilityAPIKey: getEnv("STABILITY_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
