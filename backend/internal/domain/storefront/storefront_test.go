package storefront

import (
	"strings"
	"testing"

	"shop24/internal/domain/models"
)

func TestNormalizeDropsUnsupportedSections(t *testing.T) {
	cfg := models.StorefrontConfig{
		Theme: models.Theme{Name: "luxury", Colors: models.Colors{Background: "red"}},
		Sections: []models.Section{
			{Type: "script", Visible: true},
			{Type: SectionHero, Visible: true, Hero: &models.HeroProps{Headline: "Wool"}},
		},
	}
	out := Normalize(cfg)
	for _, section := range out.Sections {
		if section.Type == "script" {
			t.Fatal("unsupported section was kept")
		}
	}
	if !has(out, SectionHero) || !has(out, SectionProductGrid) {
		t.Fatal("required sections missing")
	}
	if out.Theme.Colors.Background != Preset("luxury").Colors.Background {
		t.Fatalf("invalid color was kept: %s", out.Theme.Colors.Background)
	}
}

func TestDefaultIsStable(t *testing.T) {
	cfg := Default("North", "Wool coats")
	again := Normalize(cfg)
	if cfg.Theme.Name != again.Theme.Name || len(cfg.Sections) != len(again.Sections) {
		t.Fatal("normalize changed a default storefront")
	}
	if strings.Contains(cfg.Footer.Text, "<script") {
		t.Fatal("default footer looks executable")
	}
}

func has(cfg models.StorefrontConfig, kind string) bool {
	for _, section := range cfg.Sections {
		if section.Type == kind {
			return true
		}
	}
	return false
}
