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
	key = validateAndPadKey(key)

	return &AESJWTUtil{
		jwtSecret: jwtSecret,
		aesKey:    key,
	}
}

func validateAndPadKey(key []byte) []byte {
	validLengths := []int{16, 24, 32}
	if isValidKeyLength(key, validLengths) {
		return key
	}

	return padKeyToValidLength(key)
}

func isValidKeyLength(key []byte, validLengths []int) bool {
	for _, length := range validLengths {
		if len(key) == length {
			return true
		}
	}
	return false
}

func padKeyToValidLength(key []byte) []byte {
	targetLength := getTargetLength(len(key))

	for len(key) < targetLength {
		key = append(key, 0)
	}

	return key[:targetLength]
}

func getTargetLength(currentLength int) int {
	switch {
	case currentLength < 16:
		return 16
	case currentLength < 24:
		return 24
	case currentLength < 32:
		return 32
	default:
		return 32
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

	// Buat dan return token
	return a.createSignedToken(encryptedData)
}

func (a *AESJWTUtil) createSignedToken(encryptedData string) (string, error) {
	claims := a.createClaims(encryptedData)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.jwtSecret))
}

func (a *AESJWTUtil) createClaims(encryptedData string) *EncryptedClaims {
	expirationTime := time.Now().Add(24 * time.Hour)
	return &EncryptedClaims{
		EncryptedData: encryptedData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "HOVARLAY-BE",
		},
	}
}

func (a *AESJWTUtil) VerifyToken(tokenString string) (*UserData, error) {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Decrypt and unmarshal data
	return a.decryptAndUnmarshalUserData(claims.EncryptedData)
}

func (a *AESJWTUtil) parseToken(tokenString string) (*EncryptedClaims, error) {
	claims := &EncryptedClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, a.getSigningKeyFunc())
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (a *AESJWTUtil) getSigningKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	}
}

func (a *AESJWTUtil) decryptAndUnmarshalUserData(encryptedData string) (*UserData, error) {
	// Decrypt data
	decryptedData, err := a.decrypt(encryptedData)
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
