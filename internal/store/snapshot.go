package store

import (
	"encoding/json"
	"os"
)

func (s *KVStore) SaveSnapshot(filepath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(s.items)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

func (s *KVStore) LoadSnapshot(filepath string) error {

	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	err = json.Unmarshal(data, &s.items)
	return err
}
