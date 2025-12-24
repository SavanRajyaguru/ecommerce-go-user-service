package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/SavanRajyaguru/ecommerce-go-config-service/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	AppPort   string
	JWTSecret string
	DB        DBConfig
	Redis     RedisConfig
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

var AppConfig *Config

func LoadConfig() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	AppConfig = &Config{
		AppPort:   getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", ""),
	}

	if AppConfig.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required in environment")
	}

	// 2. Fetch from Config Service
	configServiceURL := getEnv("CONFIG_SERVICE_URL", "localhost:50051")
	log.Printf("Connecting to Config Service at: %s", configServiceURL)
	fetchRemoteConfig(configServiceURL)
}

func fetchRemoteConfig(url string) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	client := pb.NewConfigServiceClient(conn)

	// Retry logic for the RPC call
	var resp *pb.GetConfigResponse

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err = client.GetConfig(ctx, &pb.GetConfigRequest{
			ServiceName: "user-service",
		})
		cancel()

		if err == nil {
			break
		}

		log.Printf("Failed to fetch config (attempt %d/10): %v. Retrying in 2s...", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to fetch config from %s after retries: %v", url, err)
	}

	var serviceCfg ServiceConfigWrapper
	if err := json.Unmarshal([]byte(resp.ConfigJson), &serviceCfg); err != nil {
		log.Fatalf("Failed to unmarshal config json: %v", err)
	}
	// The config service returns the inner struct directly?
	// Based on server.go: data = s.Config.UserService
	// so it is { "db": ..., "redis": ... }
	// My ServiceConfigWrapper should reflect that.

	// Wait, let's verify if I need a wrapper or not.
	// If I unmarshal into AppConfig.DB directly? No, specific structs.
	// Let's use a temporary struct to match JSON.

	var remoteCfg struct {
		DB    DBConfig    `json:"db"`
		Redis RedisConfig `json:"redis"`
	}

	if err := json.Unmarshal([]byte(resp.ConfigJson), &remoteCfg); err != nil {
		log.Fatalf("Failed to unmarshal config json: %v", err)
	}

	AppConfig.DB = remoteCfg.DB
	AppConfig.Redis = remoteCfg.Redis

	fmt.Println("Configuration loaded successfully")
}

type ServiceConfigWrapper struct {
	// This was a placeholder
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
