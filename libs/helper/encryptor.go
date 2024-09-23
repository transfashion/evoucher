package helper

import (
	"encoding/base64"
	"log"
	"net/url"
)

func Encrypt(str string) (string, error) {
	// Enkripsi data yang sudah dikompres
	encryptedData, err := EncryptAES([]byte(str))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	// Compress data hasil encrypt
	// compressedData, err := Compress([]byte(encryptedData))
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return "", err
	// }

	compressedData := []byte(encryptedData)

	// Encode dengan Base64
	encodedData := base64.StdEncoding.EncodeToString(compressedData)

	// URL encode hasilnya
	urlEncodedData := url.QueryEscape(encodedData)
	return urlEncodedData, nil
}

func Decrypt(str string) (string, error) {
	// Decode dengan Base64
	decodedData, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	// Decompress data
	// decompressedData, err := Decompress(decodedData)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return "", err
	// }

	decompressedData := decodedData

	// Decrypt data
	decryptedData, err := DecryptAES(string(decompressedData))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return decryptedData, nil
}
