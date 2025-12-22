package middleware

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// ใช้สำหรับอ่าน secret จาก ENV: JWT_SECRET
var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// default ชั่วคราว กรณีไม่ได้ตั้งค่า ENV (ควรตั้งค่าใน .env.config)
		secret = "CHANGE_ME_JWT_SECRET"
	}
	return secret
}

// GenerateRefreshToken: ยังใช้ UUID แบบเดิมสำหรับ refresh token
func GenerateRefreshToken() string {
	return uuid.NewString()
}

// GenerateToken: สร้าง JWT จริง ๆ สำหรับ access token
// exp: 15 นาที
func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
