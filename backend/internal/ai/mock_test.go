package ai

import (
	"context"
	"strings"
	"testing"

	"shop24/internal/domain/models"
	"shop24/internal/domain/storefront"
)

func TestMockGenerateMinimalMonochrome(t *testing.T) {
	result, err := NewMock().Generate(context.Background(),
		"Create a modern minimalist clothing store with a black and white theme, large product images, a hero section, featured products, and a clean navigation bar.",
		StoreContext{Name: "North", Description: "Wool coats"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Provider != "mock" || !strings.Contains(result.Summary, "not a live model") {
		t.Fatalf("summary does not identify the mock designer: %s", result.Summary)
	}
	if result.Config.Theme.Name != "mono" || result.Config.Theme.ImageRatio != "portrait" {
		t.Fatalf("theme = %+v", result.Config.Theme)
	}
	grid := findSection(t, result.Config, storefront.SectionProductGrid)
	if grid.ProductGrid == nil || grid.ProductGrid.Columns != 2 {
		t.Fatalf("grid = %+v", grid.ProductGrid)
	}
	if !findSection(t, result.Config, storefront.SectionFeatured).Visible {
		t.Fatal("featured section hidden")
	}
}

func TestMockRefineKeepsCopy(t *testing.T) {
	mock := NewMock()
	current, err := mock.Generate(context.Background(), "A calm editorial shop for linen goods.", StoreContext{Name: "North", Description: "Linen"})
	if err != nil {
		t.Fatal(err)
	}
	headline := heroHeadline(current.Config)
	darker, err := mock.Refine(context.Background(), "Make the background darker.", current.Config, StoreContext{Name: "North"})
	if err != nil {
		t.Fatal(err)
	}
	if heroHeadline(darker.Config) != headline {
		t.Fatal("refine replaced the headline")
	}
	if darker.Config.Theme.Colors.Background == current.Config.Theme.Colors.Background {
		t.Fatal("background did not change")
	}
	same, err := mock.Refine(context.Background(), "Please make it nicer", darker.Config, StoreContext{Name: "North"})
	if err != nil {
		t.Fatal(err)
	}
	if same.Config.Theme.Colors.Background != darker.Config.Theme.Colors.Background || heroHeadline(same.Config) != headline {
		t.Fatal("unrecognized refine changed the design")
	}
	if !strings.Contains(same.Summary, "left as-is") {
		t.Fatal(same.Summary)
	}
}

func TestMockChangesStructure(t *testing.T) {
	cfg := storefront.Default("North", "Wool")
	result, err := NewMock().Refine(context.Background(),
		"Use a full-bleed hero, show products as a list, outline buttons, and hide the collections.",
		cfg, StoreContext{Name: "North"},
	)
	if err != nil {
		t.Fatal(err)
	}
	theme := result.Config.Theme
	if theme.HeroLayout != "bleed" || theme.ProductDisplay != "list" || theme.Buttons != "outline" || theme.Collections != "hidden" {
		t.Fatalf("layout = %+v", theme)
	}
	if !strings.Contains(result.Summary, "full-bleed") {
		t.Fatal(result.Summary)
	}
}

func TestMockMovesFeaturedAboveAbout(t *testing.T) {
	cfg := storefront.Default("North", "Wool")
	result, err := NewMock().Refine(context.Background(), "Move the featured products above the about section.", cfg, StoreContext{Name: "North"})
	if err != nil {
		t.Fatal(err)
	}
	featuredAt, aboutAt := -1, -1
	for i, section := range result.Config.Sections {
		if section.Type == storefront.SectionFeatured {
			featuredAt = i
		}
		if section.Type == storefront.SectionAbout {
			aboutAt = i
		}
	}
	if featuredAt < 0 || aboutAt < 0 || featuredAt > aboutAt {
		t.Fatalf("featured=%d about=%d", featuredAt, aboutAt)
	}
}

func findSection(t *testing.T, cfg models.StorefrontConfig, kind string) models.Section {
	t.Helper()
	for _, section := range cfg.Sections {
		if section.Type == kind {
			return section
		}
	}
	t.Fatalf("missing section %s", kind)
	return models.Section{}
}

func heroHeadline(cfg models.StorefrontConfig) string {
	for _, section := range cfg.Sections {
		if section.Hero != nil {
			return section.Hero.Headline
		}
	}
	return ""
}
