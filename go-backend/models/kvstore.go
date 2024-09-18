package models

import (
	"gorm.io/gorm"
)

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

// KVStore represents the key-value store
type KVStore struct {
	db *gorm.DB
}

var kvStore *KVStore

// InitializeKVStore initializes the key-value store
func InitializeKVStore(db *gorm.DB) {
	kvStore = &KVStore{db: db}
}

// Put stores a key-value pair
func PutKV(key, value string) error {
	kv := KeyValue{Key: key, Value: value}
	return kvStore.db.Save(&kv).Error
}

// Get retrieves the value for a given key
func GetKV(key string) (string, error) {
	var kv KeyValue
	if err := kvStore.db.First(&kv, "key = ?", key).Error; err != nil {
		return "", err
	}
	return kv.Value, nil
}

// Delete removes a key-value pair
func DeleteKV(key string) error {
	return kvStore.db.Delete(&KeyValue{}, "key = ?", key).Error
}
