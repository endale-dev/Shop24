package models

import "time"

type Store struct {
	ID          string           `json:"id"`
	OwnerID     string           `json:"ownerId"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	LogoURL     string           `json:"logoUrl"`
	Currency    string           `json:"currency"`
	Published   bool             `json:"published"`
	Storefront  StorefrontConfig `json:"storefront"`
	Draft       StorefrontConfig `json:"draft"`
	DesignNote  *DesignNote      `json:"designNote,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

type DesignNote struct {
	Provider string `json:"provider"`
	Mode     string `json:"mode"`
	Prompt   string `json:"prompt"`
	Summary  string `json:"summary"`
}

type PublicStore struct {
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	LogoURL     string           `json:"logoUrl"`
	Currency    string           `json:"currency"`
	Storefront  StorefrontConfig `json:"storefront"`
}

func (s Store) Public() PublicStore {
	return PublicStore{
		Name:        s.Name,
		Slug:        s.Slug,
		Description: s.Description,
		LogoURL:     s.LogoURL,
		Currency:    s.Currency,
		Storefront:  s.Storefront,
	}
}
