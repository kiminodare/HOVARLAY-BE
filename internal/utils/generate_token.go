package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AESJWTUtil struct {
	jwtSecret string
	aesKey    []byte
}

func NewAESJWTUtil(jwtSecret string, aesKey string) *AESJWTUtil {
	key := []byte(aesKey)

	validLengths := []int{16, 24, 32}
	valid := false
	for _, length := range validLengths {
		if len(key) == length {
			valid = true
			break
		}
	}

	if !valid {
		if len(key) < 16 {
			for len(key) < 16 {
				key = append(key, 0)
			}
			key = key[:16]
		} else if len(key) < 24 {
			for len(key) < 24 {
				key = append(key, 0)
			}
			key = key[:24]
		} else if len(key) < 32 {
			for len(key) < 32 {
				key = append(key, 0)
			}
			key = key[:32]
		} else {
			key = key[:32]
		}
	}

	return &AESJWTUtil{
		jwtSecret: jwtSecret,
		aesKey:    key,
	}
}

type EncryptedClaims struct {
	EncryptedData string `json:"data"` // Data terenkripsi
	jwt.RegisteredClaims
}

type UserData struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

func (a *AESJWTUtil) encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(a.aesKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (a *AESJWTUtil) decrypt(ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(a.aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertextBytes := data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (a *AESJWTUtil) GenerateToken(userID uuid.UUID, email string) (string, error) {
	// Buat user data
	userData := UserData{
		UserID: userID,
		Email:  email,
	}

	// Marshal ke JSON
	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		return "", err
	}

	// Enkripsi data
	encryptedData, err := a.encrypt(userDataJSON)
	if err != nil {
		return "", err
	}

	// Buat claims dengan data terenkripsi
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &EncryptedClaims{
		EncryptedData: encryptedData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "HOVARLAY-BE",
		},
	}

	// Buat dan sign JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AESJWTUtil) VerifyToken(tokenString string) (*UserData, error) {
	claims := &EncryptedClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Decrypt data
	decryptedData, err := a.decrypt(claims.EncryptedData)
	if err != nil {
		return nil, err
	}

	// Unmarshal ke UserData
	var userData UserData
	err = json.Unmarshal(decryptedData, &userData)
	if err != nil {
		return nil, err
	}

	return &userData, nil
}

func (a *AESJWTUtil) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	userData, err := a.VerifyToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	return userData.UserID, nil
}
