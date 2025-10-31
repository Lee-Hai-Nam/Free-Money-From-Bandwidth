package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// AppCredentials represents stored credentials for an app
type AppCredentials struct {
	AppID       string
	DeviceName  string
	Credentials map[string]string
}

// CredentialStore manages encrypted storage of app credentials
type CredentialStore struct {
	filePath string
	key      []byte
	mu       sync.RWMutex
}

// NewCredentialStore creates a new credential store
func NewCredentialStore() *CredentialStore {
	return &CredentialStore{
		filePath: "app_credentials.json.enc",
		key:      getOrCreateKey(),
	}
}

// SaveCredentials saves credentials for an app
func (cs *CredentialStore) SaveCredentials(creds *AppCredentials) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Load existing credentials
	allCreds, err := cs.loadDecrypted()
	if err != nil {
		allCreds = make(map[string]*AppCredentials)
	}

	// Add or update credentials
	allCreds[creds.AppID] = creds

	// Save encrypted
	return cs.saveEncrypted(allCreds)
}

// LoadCredentials loads credentials for an app
func (cs *CredentialStore) LoadCredentials(appID string) (*AppCredentials, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	allCreds, err := cs.loadDecrypted()
	if err != nil {
		return nil, err
	}

	creds, exists := allCreds[appID]
	if !exists {
		return nil, fmt.Errorf("credentials not found for app: %s", appID)
	}

	return creds, nil
}

// LoadAllCredentials loads all stored credentials
func (cs *CredentialStore) LoadAllCredentials() (map[string]*AppCredentials, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	return cs.loadDecrypted()
}

// DeleteCredentials removes credentials for an app
func (cs *CredentialStore) DeleteCredentials(appID string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	allCreds, err := cs.loadDecrypted()
	if err != nil {
		return err
	}

	delete(allCreds, appID)
	return cs.saveEncrypted(allCreds)
}

// Load all app instances (for proxy auto-deployment)
func (cs *CredentialStore) GetAllConfiguredApps() ([]string, error) {
	allCreds, err := cs.LoadAllCredentials()
	if err != nil {
		return []string{}, nil
	}

	appIDs := make([]string, 0, len(allCreds))
	for appID := range allCreds {
		appIDs = append(appIDs, appID)
	}

	return appIDs, nil
}

// Helper functions for encryption
func (cs *CredentialStore) loadDecrypted() (map[string]*AppCredentials, error) {
	encryptedData, err := os.ReadFile(cs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]*AppCredentials), nil
		}
		return nil, err
	}

	if len(encryptedData) == 0 {
		return make(map[string]*AppCredentials), nil
	}

	decrypted, err := decrypt(encryptedData, cs.key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt credentials: %w", err)
	}

	var result map[string]*AppCredentials
	if err := json.Unmarshal(decrypted, &result); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return result, nil
}

func (cs *CredentialStore) saveEncrypted(creds map[string]*AppCredentials) error {
	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	encrypted, err := encrypt(data, cs.key)
	if err != nil {
		return fmt.Errorf("failed to encrypt credentials: %w", err)
	}

	if err := os.WriteFile(cs.filePath, encrypted, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

func getOrCreateKey() []byte {
	// In production, use a proper key derivation
	// For now, use a simple derived key
	const secret = "bandwidth-income-manager-secret-key-2024"
	hash := sha256.Sum256([]byte(secret))
	return hash[:]
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
