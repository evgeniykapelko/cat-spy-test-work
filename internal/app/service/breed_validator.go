package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"spy_cat/internal/app/model"
	"sync"
	"time"
)

type BreedValidator struct {
	mu        sync.RWMutex
	breeds    map[string]bool
	lastFetch time.Time
}

func NewBreedValidator() *BreedValidator {
	return &BreedValidator{
		breeds: make(map[string]bool),
	}
}

func (v *BreedValidator) fetchBreeds() error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch breeds from TheCatAPI")
	}

	var breeds []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return err
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	v.breeds = make(map[string]bool)
	for _, breed := range breeds {
		v.breeds[breed.Name] = true
	}
	v.lastFetch = time.Now()
	return nil
}

func (v *BreedValidator) Validate(cat *model.Cat) error {
	v.mu.RLock()
	if time.Since(v.lastFetch) > 24*time.Hour {
		v.mu.RUnlock()
		if err := v.fetchBreeds(); err != nil {
			return err
		}
		v.mu.RLock()
	}
	defer v.mu.RUnlock()

	if v.breeds[cat.Breed] {
		return nil
	}
	return errors.New("invalid breed")
}
