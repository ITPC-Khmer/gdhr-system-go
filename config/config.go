package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppPort    string
	AppEnv     string
	CorsOrigin string
	StaticDir  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret      string
	JWTExpireHours int

	// External GDHR API
	GDHRBaseURL string
	GDHRUser    string
	GDHRPass    string
	GDHRKey     string

	// Scheduled sync
	SyncEnabled            bool
	SyncTimezone           string
	SyncWindowStart        string
	SyncWindowEnd          string
	SyncInstitutesInterval time.Duration
	SyncStaffsInterval     time.Duration
	SyncHTTPTimeout        time.Duration

	// Redis (durable page cursor for the sync)
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Image proxy — allowlist of remote hosts the proxy may fetch images from
	ImageProxyHosts []string
}

// Cfg is the global config instance.
var Cfg *Config

// Load reads the .env file (if present) and populates Cfg.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, falling back to system environment variables")
	}

	Cfg = &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		CorsOrigin: getEnv("CORS_ORIGIN", "http://localhost:5173"),
		StaticDir:  getEnv("STATIC_DIR", "./frontend/dist"),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "admin_system"),

		JWTSecret:      getEnv("JWT_SECRET", "insecure-default-secret-change-me"),
		JWTExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 72),

		GDHRBaseURL: getEnv("GDHR_BASE_URL", "https://apiatd-gdhr.interior.gov.kh/api/v2/gdhr"),
		GDHRUser:    getEnv("GDHR_USER", ""),
		GDHRPass:    getEnv("GDHR_PASS", ""),
		GDHRKey:     getEnv("GDHR_KEY", ""),

		SyncEnabled:            getEnvBool("SYNC_ENABLED", false),
		SyncTimezone:           getEnv("SYNC_TIMEZONE", "Asia/Phnom_Penh"),
		SyncWindowStart:        getEnv("SYNC_WINDOW_START", "18:00"),
		SyncWindowEnd:          getEnv("SYNC_WINDOW_END", "05:00"),
		SyncInstitutesInterval: getEnvDuration("SYNC_INSTITUTES_INTERVAL", 20*time.Second),
		SyncStaffsInterval:     getEnvDuration("SYNC_STAFFS_INTERVAL", 10*time.Second),
		SyncHTTPTimeout:        getEnvDuration("SYNC_HTTP_TIMEOUT", 30*time.Second),

		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		ImageProxyHosts: getEnvList("IMAGE_PROXY_HOSTS", "police-gdhr.interior.gov.kh"),
	}

	return Cfg
}

// getEnvList reads a comma-separated env value into a trimmed, non-empty slice.
func getEnvList(key, fallback string) []string {
	v := getEnv(key, fallback)
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
