package ai

import (
	"context"
	"regexp"
	"strings"

	"shop24/internal/domain/models"
	"shop24/internal/domain/storefront"
)

// Mock is a deterministic local designer. It is not a language model.
type Mock struct{}

func NewMock() *Mock { return &Mock{} }

func (m *Mock) Name() string { return "mock" }

func (m *Mock) Generate(ctx context.Context, prompt string, store StoreContext) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	cfg := storefront.Default(store.Name, store.Description)
	notes := applyPrompt(&cfg, strings.ToLower(prompt), true)
	if quote := quoted(prompt); quote != "" {
		setHeadline(&cfg, quote)
		notes = append(notes, "a headline taken from the prompt")
	}
	cfg = storefront.Normalize(cfg)
	return Result{Config: cfg, Provider: m.Name(), Summary: summarize(notes, true)}, nil
}

func (m *Mock) Refine(ctx context.Context, prompt string, current models.StorefrontConfig, store StoreContext) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	cfg := storefront.Normalize(current)
	notes := applyPrompt(&cfg, strings.ToLower(prompt), false)
	if quote := quoted(prompt); quote != "" {
		setHeadline(&cfg, quote)
		notes = append(notes, "an updated headline")
	}
	cfg = storefront.Normalize(cfg)
	return Result{Config: cfg, Provider: m.Name(), Summary: summarize(notes, false)}, nil
}

func applyPrompt(cfg *models.StorefrontConfig, prompt string, generate bool) []string {
	notes := []string{}
	if style := detectStyle(prompt); style != "" {
		cfg.Theme = storefront.Preset(style)
		if style == "mono" && strings.Contains(prompt, "minimal") {
			cfg.Theme.FontHeading = "Outfit"
			cfg.Theme.FontBody = "Figtree"
			cfg.Theme.Radius = "sm"
		}
		notes = append(notes, "a "+style+" theme")
	} else if generate {
		notes = append(notes, "an editorial theme")
	}
	if generate {
		setHeroSize(cfg, "md")
	}
	if strings.Contains(prompt, "large product") || strings.Contains(prompt, "large image") || strings.Contains(prompt, "big image") {
		cfg.Theme.ImageRatio = "portrait"
		setColumns(cfg, 2)
		notes = append(notes, "large portrait product images")
	}
	if strings.Contains(prompt, "hero") && (strings.Contains(prompt, "larger") || strings.Contains(prompt, "bigger")) {
		setHeroSize(cfg, "lg")
		notes = append(notes, "a larger hero")
	}
	if strings.Contains(prompt, "hero") && (strings.Contains(prompt, "smaller") || strings.Contains(prompt, "shorter")) {
		setHeroSize(cfg, "sm")
		notes = append(notes, "a smaller hero")
	}
	if strings.Contains(prompt, "darker") || strings.Contains(prompt, "dark background") {
		darken(cfg)
		notes = append(notes, "a darker background")
	}
	if strings.Contains(prompt, "lighter") || strings.Contains(prompt, "light background") {
		lighten(cfg)
		notes = append(notes, "a lighter background")
	}
	if generate && !strings.Contains(prompt, "newsletter") {
		setVisible(cfg, storefront.SectionNewsletter, false)
	}
	if moving, anchor, ok := moveRequest(prompt); ok {
		moveAbove(cfg, moving, anchor)
		notes = append(notes, "reordered sections")
	}
	for _, kind := range []string{
		storefront.SectionHero, storefront.SectionProductGrid, storefront.SectionFeatured,
		storefront.SectionAbout, storefront.SectionAnnouncement, storefront.SectionCategories, storefront.SectionNewsletter,
	} {
		label := sectionLabel(kind)
		if strings.Contains(prompt, "hide "+label) || strings.Contains(prompt, "remove "+label) || strings.Contains(prompt, "without "+label) {
			setVisible(cfg, kind, false)
			notes = append(notes, "hid the "+label+" section")
		}
		if strings.Contains(prompt, "show "+label) {
			setVisible(cfg, kind, true)
			notes = append(notes, "showed the "+label+" section")
		}
	}
	if strings.Contains(prompt, "featured") {
		setVisible(cfg, storefront.SectionFeatured, true)
	}
	notes = append(notes, applyLayout(cfg, prompt)...)
	return notes
}

func applyLayout(cfg *models.StorefrontConfig, prompt string) []string {
	notes := []string{}
	switch {
	case containsAny(prompt, "full bleed", "full-bleed", "banner hero"):
		cfg.Theme.HeroLayout = "bleed"
		notes = append(notes, "a full-bleed hero")
	case containsAny(prompt, "split hero"):
		cfg.Theme.HeroLayout = "split"
		notes = append(notes, "a split hero")
	case containsAny(prompt, "stacked hero", "hero stacked"):
		cfg.Theme.HeroLayout = "stack"
		notes = append(notes, "a stacked hero")
	}
	switch {
	case containsAny(prompt, "product list", "as a list", "list the products"):
		cfg.Theme.ProductDisplay = "list"
		notes = append(notes, "products in a list")
	case containsAny(prompt, "feature the first", "large first product"):
		cfg.Theme.ProductDisplay = "feature"
		notes = append(notes, "a large first product")
	}
	switch {
	case containsAny(prompt, "outline button"):
		cfg.Theme.Buttons = "outline"
		notes = append(notes, "outline buttons")
	case containsAny(prompt, "text button", "link button", "underline button"):
		cfg.Theme.Buttons = "link"
		notes = append(notes, "text buttons")
	}
	if containsAny(prompt, "hide collection", "hide the collections", "without collections") {
		cfg.Theme.Collections = "hidden"
		notes = append(notes, "collections hidden")
	}
	if containsAny(prompt, "more space", "more room", "roomy", "spacious") {
		cfg.Theme.Density = "roomy"
		notes = append(notes, "more space")
	}
	if containsAny(prompt, "tighter", "less space", "compact layout") {
		cfg.Theme.Density = "tight"
		notes = append(notes, "a tighter layout")
	}
	if containsAny(prompt, "card border", "bordered card", "framed card") {
		cfg.Theme.Cards = "border"
		notes = append(notes, "bordered cards")
	}
	if strings.Contains(prompt, "card shadow") || strings.Contains(prompt, "soft shadow") {
		cfg.Theme.Cards = "shadow"
		notes = append(notes, "soft card shadows")
	}
	return notes
}

func detectStyle(prompt string) string {
	switch {
	case containsAny(prompt, "black and white", "black & white", "monochrome"):
		return "mono"
	case containsAny(prompt, "luxury", "elegant", "gold"):
		return "luxury"
	case containsAny(prompt, "playful", "pastel"):
		return "playful"
	case containsAny(prompt, "bold", "vibrant"):
		return "bold"
	case strings.Contains(prompt, "minimal"):
		return "minimal"
	case strings.Contains(prompt, "editorial"):
		return "editorial"
	case strings.Contains(strings.ReplaceAll(prompt, "darker", ""), "dark"):
		return "dark"
	default:
		return ""
	}
}

func containsAny(value string, parts ...string) bool {
	for _, part := range parts {
		if strings.Contains(value, part) {
			return true
		}
	}
	return false
}

func sectionLabel(kind string) string {
	switch kind {
	case storefront.SectionProductGrid:
		return "product grid"
	case storefront.SectionFeatured:
		return "featured products"
	default:
		return strings.ReplaceAll(kind, "_", " ")
	}
}

var movePattern = regexp.MustCompile(`move (?:the )?([a-z ]+?) above (?:the )?([a-z ]+)`)
var quotePattern = regexp.MustCompile(`"([^"]{1,80})"|“([^”]{1,80})”`)

func moveRequest(prompt string) (string, string, bool) {
	match := movePattern.FindStringSubmatch(prompt)
	if match == nil {
		return "", "", false
	}
	moving := alias(match[1])
	anchor := alias(match[2])
	if moving == "" || anchor == "" || moving == anchor {
		return "", "", false
	}
	return moving, anchor, true
}

func alias(value string) string {
	value = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(value), " section"))
	switch value {
	case "hero":
		return storefront.SectionHero
	case "featured", "featured products", "featured product":
		return storefront.SectionFeatured
	case "products", "product grid", "grid", "collection":
		return storefront.SectionProductGrid
	case "about":
		return storefront.SectionAbout
	case "announcement", "banner":
		return storefront.SectionAnnouncement
	case "categories", "category":
		return storefront.SectionCategories
	case "newsletter":
		return storefront.SectionNewsletter
	default:
		return ""
	}
}

func quoted(prompt string) string {
	match := quotePattern.FindStringSubmatch(prompt)
	if match == nil {
		return ""
	}
	if match[1] != "" {
		return strings.TrimSpace(match[1])
	}
	return strings.TrimSpace(match[2])
}

func setHeadline(cfg *models.StorefrontConfig, headline string) {
	for i := range cfg.Sections {
		if cfg.Sections[i].Hero != nil {
			cfg.Sections[i].Hero.Headline = headline
		}
	}
}

func setHeroSize(cfg *models.StorefrontConfig, size string) {
	for i := range cfg.Sections {
		if cfg.Sections[i].Hero != nil {
			cfg.Sections[i].Hero.Size = size
		}
	}
}

func setColumns(cfg *models.StorefrontConfig, columns int) {
	for i := range cfg.Sections {
		if cfg.Sections[i].ProductGrid != nil {
			cfg.Sections[i].ProductGrid.Columns = columns
		}
	}
}

func setVisible(cfg *models.StorefrontConfig, kind string, visible bool) {
	for i := range cfg.Sections {
		if cfg.Sections[i].Type == kind {
			cfg.Sections[i].Visible = visible
		}
	}
}

func moveAbove(cfg *models.StorefrontConfig, moving, anchor string) {
	from := -1
	for i, section := range cfg.Sections {
		if section.Type == moving {
			from = i
			break
		}
	}
	if from < 0 {
		return
	}
	item := cfg.Sections[from]
	rest := append([]models.Section{}, cfg.Sections[:from]...)
	rest = append(rest, cfg.Sections[from+1:]...)
	anchorAt := -1
	for i, section := range rest {
		if section.Type == anchor {
			anchorAt = i
			break
		}
	}
	if anchorAt < 0 {
		return
	}
	out := append([]models.Section{}, rest[:anchorAt]...)
	out = append(out, item)
	out = append(out, rest[anchorAt:]...)
	cfg.Sections = out
}

func darken(cfg *models.StorefrontConfig) {
	colors := &cfg.Theme.Colors
	if luminance(colors.Background) > 140 {
		colors.Background = "#171717"
		colors.Surface = "#221f1c"
		colors.Text = "#f5f5f4"
		colors.Muted = "#a8a29e"
		colors.Border = "#3f3a36"
		return
	}
	colors.Background = scale(colors.Background, 0.8)
	colors.Surface = scale(colors.Surface, 0.85)
}

func lighten(cfg *models.StorefrontConfig) {
	colors := &cfg.Theme.Colors
	if luminance(colors.Background) < 80 {
		colors.Background = "#f6f1e8"
		colors.Surface = "#fffdf8"
		colors.Text = "#1f1a17"
		colors.Muted = "#6f655d"
		colors.Border = "#e6ddd1"
		return
	}
	colors.Background = mixWhite(colors.Background, 0.35)
	colors.Surface = mixWhite(colors.Surface, 0.35)
}

func luminance(hex string) float64 {
	r, g, b, ok := parseHex(hex)
	if !ok {
		return 255
	}
	return 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
}

func scale(hex string, factor float64) string {
	r, g, b, ok := parseHex(hex)
	if !ok {
		return hex
	}
	return formatHex(int(float64(r)*factor), int(float64(g)*factor), int(float64(b)*factor))
}

func mixWhite(hex string, amount float64) string {
	r, g, b, ok := parseHex(hex)
	if !ok {
		return hex
	}
	mix := func(channel int) int {
		return int(float64(channel) + (255-float64(channel))*amount)
	}
	return formatHex(mix(r), mix(g), mix(b))
}

func parseHex(hex string) (int, int, int, bool) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0, false
	}
	value := 0
	for _, r := range hex[1:] {
		value <<= 4
		switch {
		case r >= '0' && r <= '9':
			value += int(r - '0')
		case r >= 'a' && r <= 'f':
			value += int(r-'a') + 10
		case r >= 'A' && r <= 'F':
			value += int(r-'A') + 10
		default:
			return 0, 0, 0, false
		}
	}
	return (value >> 16) & 255, (value >> 8) & 255, value & 255, true
}

func formatHex(r, g, b int) string {
	clamp := func(v int) int {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return v
	}
	r, g, b = clamp(r), clamp(g), clamp(b)
	const digits = "0123456789abcdef"
	return "#" + string([]byte{digits[r>>4], digits[r&15], digits[g>>4], digits[g&15], digits[b>>4], digits[b&15]})
}

func summarize(notes []string, generate bool) string {
	if len(notes) == 0 {
		return "Mock designer did not recognize a specific change, so the current design was left as-is. This is a local preview, not a live model."
	}
	verb := "adjusted"
	if generate {
		verb = "applied"
	}
	return "Mock designer " + verb + " " + strings.Join(notes, ", ") + ". This is a local preview, not a live model."
}
