package utils

import "os"

func JWTSecret() []byte {
    if key := os.Getenv("SECRET_KEY"); key != "" {
        return []byte(key)
    }
    return []byte("MY_DEFAULT_SECRET_KEY")
}