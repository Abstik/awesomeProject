package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// 定义文件路径
const (
	privateKeyPath = "private_key.pem"
	publicKeyPath  = "public_key.pem"
	keySize        = 2048 // RSA 密钥长度
)

// 生成RSA密钥对，并保存到文件
func generateAndSaveRSAKeys() error {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// 编码并保存私钥到文件
	err = savePrivateKeyToFile(privateKey, privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to save private key: %v", err)
	}

	// 提取公钥
	publicKey := &privateKey.PublicKey

	// 编码并保存公钥到文件
	err = savePublicKeyToFile(publicKey, publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to save public key: %v", err)
	}

	return nil
}

// 将私钥保存为 PEM 格式
func savePrivateKeyToFile(privateKey *rsa.PrivateKey, filepath string) error {
	// 将私钥编码为PKCS1格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// 创建PEM块
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// 创建或覆盖文件
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入PEM数据
	err = pem.Encode(file, privateKeyBlock)
	if err != nil {
		return err
	}

	return nil
}

// 将公钥保存为 PEM 格式
func savePublicKeyToFile(publicKey *rsa.PublicKey, filepath string) error {
	// 将公钥编码为PKIX格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	// 创建PEM块
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// 创建或覆盖文件
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入PEM数据
	err = pem.Encode(file, publicKeyBlock)
	if err != nil {
		return err
	}

	return nil
}

// 加载私钥
func loadPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	// 读取文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// 解码PEM块
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	// 解析PKCS1格式的私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// 加载公钥
func loadPublicKey(filepath string) (*rsa.PublicKey, error) {
	// 读取文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// 解码PEM块
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	// 解析PKIX格式的公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 类型断言为 *rsa.PublicKey
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPublicKey, nil
}

func GenerateRsa() {
	// 如果密钥文件不存在，则生成密钥对
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		fmt.Println("Keys not found, generating new keys...")
		err := generateAndSaveRSAKeys()
		if err != nil {
			fmt.Printf("Error generating keys: %v\n", err)
			return
		}
		fmt.Println("Keys generated and saved successfully.")
	} else {
		fmt.Println("Keys already exist, skipping generation.")
	}

	// 加载私钥
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		fmt.Printf("Error loading private key: %v\n", err)
		return
	}
	fmt.Printf("Private Key Loaded: %v\n", privateKey)

	// 加载公钥
	publicKey, err := loadPublicKey(publicKeyPath)
	if err != nil {
		fmt.Printf("Error loading public key: %v\n", err)
		return
	}
	fmt.Printf("Public Key Loaded: %v\n", publicKey)
}

// 加载私钥
func loadPrivateKeyFromPem(filepath string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// 加载公钥
func loadPublicKeyFromPem(filepath string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaPublicKey, nil
}

// 使用公钥加密
func encryptWithPublicKey(publicKey *rsa.PublicKey, message string) ([]byte, error) {
	// 使用 SHA-256 作为 OAEP 的哈希函数
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(
		hash,            // 哈希函数
		rand.Reader,     // 随机数生成器
		publicKey,       // 公钥
		[]byte(message), // 待加密的消息
		nil,             // 可选的 label 参数 (通常为 nil)
	)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// 使用私钥解密
func decryptWithPrivateKey(privateKey *rsa.PrivateKey, ciphertext []byte) (string, error) {
	// 使用 SHA-256 作为 OAEP 的哈希函数
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(
		hash,        // 哈希函数
		rand.Reader, // 随机数生成器
		privateKey,  // 私钥
		ciphertext,  // 加密的消息
		nil,         // 可选的 label 参数 (通常为 nil)
	)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
