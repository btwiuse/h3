package utils

import "os"

func EnvCERT(fallback string) string {
	if cert, ok := os.LookupEnv("CERT"); ok {
		return cert
	}
	return fallback
}

func EnvKEY(fallback string) string {
	if key, ok := os.LookupEnv("KEY"); ok {
		return key
	}
	return fallback
}

func EnvPORT(fallback string) string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return fallback
}
