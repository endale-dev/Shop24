package ai

import (
	"encoding/json"
	"testing"

	"shop24/internal/domain/models"
)

func TestPrepareConfigAcceptsStringNumbers(t *testing.T) {
	raw, err := prepareConfig(`{"version":"1","theme":"minimal","sections":[{"type":"hero","visible":"true","hero":{"text":"Wool"}},{"type":"product_grid","visible":"true","productGrid":{"columns":"3","limit":"8"}},{"type":"about","visible":true,"about":{"text":"Made here"}}]}`)
	if err != nil {
		t.Fatal(err)
	}
	var cfg models.StorefrontConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Version != 1 || cfg.Theme.Name != "minimal" || cfg.Sections[0].Hero.Headline != "Wool" || cfg.Sections[2].About.Body != "Made here" {
		t.Fatalf("coercion failed: %+v", cfg)
	}
}
