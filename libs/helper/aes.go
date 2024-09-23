package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const keyString = "GakUsahDengerin!"

var key = []byte(keyString)

// Fungsi untuk mengenkripsi data dengan AES
func EncryptAES(text []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Buat byte array untuk hasil enkripsi
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]

	// Mengisi IV dengan data acak
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)

	// Encode ke dalam Base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Fungsi untuk mendekripsi data dengan AES
func DecryptAES(cryptoText string) (string, error) {
	ciphertext, _ := base64.StdEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Ambil IV dari ciphertext
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext terlalu pendek")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// Dekripsi data
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
