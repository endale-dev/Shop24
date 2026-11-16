package storefront

import (
	"regexp"
	"strings"

	"shop24/internal/domain/models"
)

const Version = 1

const (
	SectionHero         = "hero"
	SectionProductGrid  = "product_grid"
	SectionFeatured     = "featured_products"
	SectionAbout        = "about"
	SectionAnnouncement = "announcement"
	SectionCategories   = "categories"
	SectionNewsletter   = "newsletter"
)

var (
	hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	themes   = map[string]struct{}{
		"editorial": {}, "minimal": {}, "luxury": {}, "bold": {}, "playful": {}, "mono": {}, "dark": {},
	}
	fonts = map[string]struct{}{
		"Instrument Serif": {}, "Fraunces": {}, "Newsreader": {}, "Outfit": {}, "Figtree": {}, "Literata": {},
	}
	targets = map[string]struct{}{"home": {}, "products": {}, "about": {}, "cart": {}}
)

func Default(name, description string) models.StorefrontConfig {
	name = fallback(strings.TrimSpace(name), "Your store")
	description = strings.TrimSpace(description)
	sub := description
	if sub == "" {
		sub = "A smaller catalog, chosen with care."
	}
	about := description
	if about == "" {
		about = name + " is an independent shop. Every product on this page is part of the live catalog."
	}
	cfg := models.StorefrontConfig{
		Version: Version,
		Theme:   Preset("editorial"),
		Header: models.Header{
			Announcement:     "Cash on delivery. Orders are confirmed by the shop.",
			ShowAnnouncement: false,
			Alignment:        "left",
			Nav:              defaultNav(),
		},
		Sections: []models.Section{
			heroSection(name, sub, "lg", true),
			featuredSection(true),
			gridSection(true),
			categoriesSection(true),
			aboutSection(name, about, true),
			newsletterSection(false),
		},
		Footer: models.Footer{
			Text:  "© " + name,
			Links: defaultNav(),
		},
	}
	return Normalize(cfg)
}

func Preset(name string) models.Theme {
	presets := map[string]models.Theme{
		"editorial": {
			Name: "editorial", FontHeading: "Fraunces", FontBody: "Figtree", Radius: "md", ImageRatio: "portrait",
			Colors: models.Colors{Background: "#f6f1e8", Surface: "#fffdf8", Text: "#1f1a17", Muted: "#6f655d", Primary: "#9a3412", OnPrimary: "#fff7ed", Accent: "#b45309", Border: "#e6ddd1"},
		},
		"minimal": {
			Name: "minimal", FontHeading: "Outfit", FontBody: "Figtree", Radius: "sm", ImageRatio: "square",
			Colors: models.Colors{Background: "#ffffff", Surface: "#fafafa", Text: "#171717", Muted: "#737373", Primary: "#171717", OnPrimary: "#ffffff", Accent: "#404040", Border: "#e5e5e5"},
		},
		"mono": {
			Name: "mono", FontHeading: "Outfit", FontBody: "Figtree", Radius: "sm", ImageRatio: "square",
			Colors: models.Colors{Background: "#ffffff", Surface: "#f4f4f5", Text: "#111111", Muted: "#52525b", Primary: "#111111", OnPrimary: "#ffffff", Accent: "#18181b", Border: "#e4e4e7"},
		},
		"luxury": {
			Name: "luxury", FontHeading: "Instrument Serif", FontBody: "Figtree", Radius: "none", ImageRatio: "portrait",
			Colors: models.Colors{Background: "#100e0c", Surface: "#1a1714", Text: "#f6f1e7", Muted: "#b7ab9a", Primary: "#c6a15b", OnPrimary: "#1a1408", Accent: "#e7d3a1", Border: "#3a342c"},
		},
		"bold": {
			Name: "bold", FontHeading: "Outfit", FontBody: "Figtree", Radius: "lg", ImageRatio: "square",
			Colors: models.Colors{Background: "#fff7ed", Surface: "#ffffff", Text: "#1c1917", Muted: "#78716c", Primary: "#ea580c", OnPrimary: "#ffffff", Accent: "#0f172a", Border: "#fed7aa"},
		},
		"playful": {
			Name: "playful", FontHeading: "Fraunces", FontBody: "Outfit", Radius: "lg", ImageRatio: "square",
			Colors: models.Colors{Background: "#fdf2f8", Surface: "#ffffff", Text: "#3b0764", Muted: "#7e22ce", Primary: "#db2777", OnPrimary: "#ffffff", Accent: "#0891b2", Border: "#fbcfe8"},
		},
		"dark": {
			Name: "dark", FontHeading: "Outfit", FontBody: "Figtree", Radius: "md", ImageRatio: "portrait",
			Colors: models.Colors{Background: "#121212", Surface: "#1c1c1c", Text: "#f5f5f5", Muted: "#a3a3a3", Primary: "#f5f5f5", OnPrimary: "#111111", Accent: "#d4d4d4", Border: "#2e2e2e"},
		},
	}
	preset, ok := presets[name]
	if !ok {
		preset = presets["editorial"]
	}
	return preset
}

// Normalize repairs a storefront document so the renderer only receives known sections and safe values.
func Normalize(in models.StorefrontConfig) models.StorefrontConfig {
	base := Preset("editorial")
	themeName := strings.ToLower(strings.TrimSpace(in.Theme.Name))
	if _, ok := themes[themeName]; ok {
		base = Preset(themeName)
	}
	if in.Theme.FontHeading != "" || in.Theme.Colors.Background != "" || in.Theme.Radius != "" {
		base = overlayTheme(base, in.Theme)
	}
	header := models.Header{
		Announcement:     clip(in.Header.Announcement, 140),
		ShowAnnouncement: in.Header.ShowAnnouncement && strings.TrimSpace(in.Header.Announcement) != "",
		Alignment:        oneOf(in.Header.Alignment, "left", "left", "center"),
		Nav:              cleanNav(in.Header.Nav),
	}
	sections := cleanSections(in.Sections)
	footer := models.Footer{
		Text:  fallback(clip(in.Footer.Text, 180), "Thank you for visiting."),
		Links: cleanNav(in.Footer.Links),
	}
	return models.StorefrontConfig{
		Version:  Version,
		Theme:    base,
		Header:   header,
		Sections: sections,
		Footer:   footer,
	}
}

func overlayTheme(base, in models.Theme) models.Theme {
	if _, ok := fonts[in.FontHeading]; ok {
		base.FontHeading = in.FontHeading
	}
	if _, ok := fonts[in.FontBody]; ok {
		base.FontBody = in.FontBody
	}
	base.Radius = oneOf(in.Radius, base.Radius, "none", "sm", "md", "lg")
	base.ImageRatio = oneOf(in.ImageRatio, base.ImageRatio, "square", "portrait", "landscape")
	base.Colors = overlayColors(base.Colors, in.Colors)
	return base
}

func overlayColors(base, in models.Colors) models.Colors {
	assign := func(dst *string, src string) {
		if color, ok := normalizeColor(src); ok {
			*dst = color
		}
	}
	assign(&base.Background, in.Background)
	assign(&base.Surface, in.Surface)
	assign(&base.Text, in.Text)
	assign(&base.Muted, in.Muted)
	assign(&base.Primary, in.Primary)
	assign(&base.OnPrimary, in.OnPrimary)
	assign(&base.Accent, in.Accent)
	assign(&base.Border, in.Border)
	return base
}

func cleanNav(items []models.NavItem) []models.NavItem {
	if len(items) == 0 {
		return defaultNav()
	}
	out := make([]models.NavItem, 0, len(items))
	seen := map[string]bool{}
	for _, item := range items {
		target := strings.ToLower(strings.TrimSpace(item.Target))
		if _, ok := targets[target]; !ok || seen[target] {
			continue
		}
		label := clip(item.Label, 24)
		if label == "" {
			continue
		}
		seen[target] = true
		out = append(out, models.NavItem{Label: label, Target: target})
		if len(out) == 6 {
			break
		}
	}
	if len(out) == 0 {
		return defaultNav()
	}
	return out
}

func cleanSections(in []models.Section) []models.Section {
	out := make([]models.Section, 0, 8)
	seen := map[string]bool{}
	for _, section := range in {
		kind := strings.TrimSpace(section.Type)
		if seen[kind] || !knownSection(kind) {
			continue
		}
		seen[kind] = true
		out = append(out, cleanSection(section, kind))
	}
	if !seen[SectionHero] {
		out = append([]models.Section{heroSection("Welcome", "Goods from this shop.", "md", true)}, out...)
	}
	if !seen[SectionProductGrid] {
		out = append(out, gridSection(true))
	}
	if len(out) > 8 {
		out = out[:8]
	}
	return out
}

func cleanSection(section models.Section, kind string) models.Section {
	cleaned := models.Section{ID: kind, Type: kind, Visible: section.Visible}
	switch kind {
	case SectionHero:
		props := models.HeroProps{}
		if section.Hero != nil {
			props = *section.Hero
		}
		cleaned.Hero = &models.HeroProps{
			Eyebrow:     clip(props.Eyebrow, 40),
			Headline:    fallback(clip(props.Headline, 80), "Welcome"),
			Subheadline: clip(props.Subheadline, 180),
			CTALabel:    fallback(clip(props.CTALabel, 24), "Shop"),
			CTATarget:   oneOf(props.CTATarget, "products", "home", "products", "about", "cart"),
			ImageURL:    cleanURL(props.ImageURL),
			Align:       oneOf(props.Align, "left", "left", "center"),
			Size:        oneOf(props.Size, "md", "sm", "md", "lg"),
		}
	case SectionProductGrid:
		props := models.ProductGridProps{Columns: 3, Limit: 12}
		if section.ProductGrid != nil {
			props = *section.ProductGrid
		}
		cleaned.ProductGrid = &models.ProductGridProps{
			Title:   fallback(clip(props.Title, 60), "All products"),
			Columns: clamp(props.Columns, 2, 4, 3),
			Limit:   clamp(props.Limit, 1, 24, 12),
		}
	case SectionFeatured:
		props := models.FeaturedProps{Limit: 4}
		if section.Featured != nil {
			props = *section.Featured
		}
		cleaned.Featured = &models.FeaturedProps{
			Title: fallback(clip(props.Title, 60), "Featured"),
			Limit: clamp(props.Limit, 1, 8, 4),
		}
	case SectionAbout:
		props := models.AboutProps{}
		if section.About != nil {
			props = *section.About
		}
		cleaned.About = &models.AboutProps{
			Title:    fallback(clip(props.Title, 60), "About"),
			Body:     fallback(clip(props.Body, 800), "This shop is run by an independent seller."),
			ImageURL: cleanURL(props.ImageURL),
		}
	case SectionAnnouncement:
		text := ""
		if section.Announcement != nil {
			text = section.Announcement.Text
		}
		cleaned.Announcement = &models.AnnouncementProps{Text: fallback(clip(text, 140), "Welcome in.")}
	case SectionCategories:
		title := ""
		if section.Categories != nil {
			title = section.Categories.Title
		}
		cleaned.Categories = &models.CategoriesProps{Title: fallback(clip(title, 60), "Categories")}
	case SectionNewsletter:
		props := models.NewsletterProps{}
		if section.Newsletter != nil {
			props = *section.Newsletter
		}
		cleaned.Newsletter = &models.NewsletterProps{
			Title: fallback(clip(props.Title, 60), "Notes from the shop"),
			Text:  fallback(clip(props.Text, 180), "New pieces are added to this page as they are published."),
		}
	}
	return cleaned
}

func knownSection(kind string) bool {
	switch kind {
	case SectionHero, SectionProductGrid, SectionFeatured, SectionAbout, SectionAnnouncement, SectionCategories, SectionNewsletter:
		return true
	default:
		return false
	}
}

func heroSection(headline, sub, size string, visible bool) models.Section {
	return models.Section{
		ID: SectionHero, Type: SectionHero, Visible: visible,
		Hero: &models.HeroProps{
			Eyebrow: "Independent shop", Headline: headline, Subheadline: sub,
			CTALabel: "Shop the collection", CTATarget: "products", Align: "left", Size: size,
		},
	}
}

func gridSection(visible bool) models.Section {
	return models.Section{
		ID: SectionProductGrid, Type: SectionProductGrid, Visible: visible,
		ProductGrid: &models.ProductGridProps{Title: "The collection", Columns: 3, Limit: 12},
	}
}

func featuredSection(visible bool) models.Section {
	return models.Section{
		ID: SectionFeatured, Type: SectionFeatured, Visible: visible,
		Featured: &models.FeaturedProps{Title: "Featured", Limit: 4},
	}
}

func aboutSection(name, body string, visible bool) models.Section {
	return models.Section{
		ID: SectionAbout, Type: SectionAbout, Visible: visible,
		About: &models.AboutProps{Title: "About " + name, Body: body},
	}
}

func categoriesSection(visible bool) models.Section {
	return models.Section{
		ID: SectionCategories, Type: SectionCategories, Visible: visible,
		Categories: &models.CategoriesProps{Title: "Shop by category"},
	}
}

func newsletterSection(visible bool) models.Section {
	return models.Section{
		ID: SectionNewsletter, Type: SectionNewsletter, Visible: visible,
		Newsletter: &models.NewsletterProps{Title: "From the shop", Text: "This note is part of the storefront design. It does not collect email addresses."},
	}
}

func defaultNav() []models.NavItem {
	return []models.NavItem{
		{Label: "Home", Target: "home"},
		{Label: "Shop", Target: "products"},
		{Label: "About", Target: "about"},
		{Label: "Cart", Target: "cart"},
	}
}

func normalizeColor(raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	if len(value) == 4 && value[0] == '#' {
		value = "#" + strings.Repeat(value[1:2], 2) + strings.Repeat(value[2:3], 2) + strings.Repeat(value[3:4], 2)
	}
	if !hexColor.MatchString(value) {
		return "", false
	}
	return strings.ToLower(value), true
}

func cleanURL(raw string) string {
	value := strings.TrimSpace(raw)
	if len(value) > 400 {
		return ""
	}
	if strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://") {
		return value
	}
	return ""
}

func clip(value string, max int) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return strings.TrimSpace(string(runes[:max]))
}

func fallback(value, backup string) string {
	if strings.TrimSpace(value) == "" {
		return backup
	}
	return value
}

func oneOf(value, backup string, allowed ...string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, item := range allowed {
		if value == item {
			return value
		}
	}
	return backup
}

func clamp(value, min, max, backup int) int {
	if value < min || value > max {
		if value == 0 {
			return backup
		}
		if value < min {
			return min
		}
		return max
	}
	return value
}
