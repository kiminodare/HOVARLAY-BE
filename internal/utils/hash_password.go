package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Params menyimpan parameter konfigurasi Argon2
type Argon2Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultArgon2Params parameter default yang direkomendasikan
var DefaultArgon2Params = Argon2Params{
	Time:    3,         // 3 iterasi
	Memory:  64 * 1024, // 64 MB RAM
	Threads: 2,         // 2 threads
	KeyLen:  32,        // 32 bytes output
}

var ErrMismatchedHashAndPassword = errors.New("hashedPassword is not the hash of the given password")

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		DefaultArgon2Params.Time,
		DefaultArgon2Params.Memory,
		DefaultArgon2Params.Threads,
		DefaultArgon2Params.KeyLen,
	)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		DefaultArgon2Params.Memory,
		DefaultArgon2Params.Time,
		DefaultArgon2Params.Threads,
		hex.EncodeToString(salt),
		hex.EncodeToString(hash),
	)

	return encoded, nil
}

func ComparePassword(password, hashedPassword string) error {
	params, salt, hash, err := parseHash(hashedPassword)
	if err != nil {
		return err
	}

	// Generate hash dari password input dengan parameter yang sama
	inputHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLen,
	)

	if hex.EncodeToString(inputHash) != hex.EncodeToString(hash) {
		return ErrMismatchedHashAndPassword
	}

	return nil
}

func parseHash(encodedHash string) (Argon2Params, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return Argon2Params{}, nil, nil, fmt.Errorf("invalid hash format")
	}

	var params Argon2Params

	paramStr := vals[3]
	_, err := fmt.Sscanf(paramStr, "m=%d,t=%d,p=%d",
		&params.Memory, &params.Time, &params.Threads)
	if err != nil {
		return Argon2Params{}, nil, nil, err
	}

	params.KeyLen = 32
	salt, err := hex.DecodeString(vals[4])
	if err != nil {
		return Argon2Params{}, nil, nil, err
	}

	hash, err := hex.DecodeString(vals[5])
	if err != nil {
		return Argon2Params{}, nil, nil, err
	}

	return params, salt, hash, nil
}
