package utils

import (
	"fmt"
	"testing"
)

func TestRsa(t *testing.T) {
	// 定义文件路径
	privateKeyPath := "private_key.pem"
	publicKeyPath := "public_key.pem"

	// 加载私钥
	privateKey, err := loadPrivateKeyFromPem(privateKeyPath)
	if err != nil {
		fmt.Printf("Error loading private key: %v\n", err)
		return
	}

	// 加载公钥
	publicKey, err := loadPublicKeyFromPem(publicKeyPath)
	if err != nil {
		fmt.Printf("Error loading public key: %v\n", err)
		return
	}

	// 要加密的消息
	message := "Hello, RSA encryption and decryption!"

	// 使用公钥加密
	ciphertext, err := encryptWithPublicKey(publicKey, message)
	if err != nil {
		fmt.Printf("Error encrypting message: %v\n", err)
		return
	}
	fmt.Printf("Encrypted message (hex): %x\n", ciphertext)

	// 使用私钥解密
	decryptedMessage, err := decryptWithPrivateKey(privateKey, ciphertext)
	if err != nil {
		fmt.Printf("Error decrypting message: %v\n", err)
		return
	}
	fmt.Printf("Decrypted message: %s\n", decryptedMessage)
}
