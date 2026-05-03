// backend/shared/models/theme.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ThemeColors represents the color palette for a storefront
type ThemeColors struct {
	Primary       string `bson:"primary" json:"primary"`             // Main brand color
	Secondary     string `bson:"secondary" json:"secondary"`         // Secondary brand color
	Accent        string `bson:"accent" json:"accent"`               // Accent/CTA color
	Background    string `bson:"background" json:"background"`       // Page background
	Surface       string `bson:"surface" json:"surface"`             // Card/surface background
	TextPrimary   string `bson:"textPrimary" json:"textPrimary"`     // Main text color
	TextSecondary string `bson:"textSecondary" json:"textSecondary"` // Secondary text
	Border        string `bson:"border" json:"border"`               // Border color
	Success       string `bson:"success" json:"success"`             // Success state
	Warning       string `bson:"warning" json:"warning"`             // Warning state
	Error         string `bson:"error" json:"error"`                 // Error state
	Info          string `bson:"info" json:"info"`                   // Info state
}

// ThemeTypography represents font settings
type ThemeTypography struct {
	HeadingFont     string `bson:"headingFont" json:"headingFont"`         // Heading font family
	BodyFont        string `bson:"bodyFont" json:"bodyFont"`               // Body font family
	HeadingWeight   string `bson:"headingWeight" json:"headingWeight"`     // Heading font weight
	BodyWeight      string `bson:"bodyWeight" json:"bodyWeight"`           // Body font weight
	HeadingSizeBase string `bson:"headingSizeBase" json:"headingSizeBase"` // Base heading size (rem)
	BodySizeBase    string `bson:"bodySizeBase" json:"bodySizeBase"`       // Base body size (rem)
	LineHeight      string `bson:"lineHeight" json:"lineHeight"`           // Line height
	LetterSpacing   string `bson:"letterSpacing" json:"letterSpacing"`     // Letter spacing
}

// ThemeLayout represents layout configuration
type ThemeLayout struct {
	HeaderStyle     string `bson:"headerStyle" json:"headerStyle"`         // Header style: "classic", "centered", "minimal"
	ContainerWidth  string `bson:"containerWidth" json:"containerWidth"`   // Container: "boxed", "full"
	ProductGridCols int    `bson:"productGridCols" json:"productGridCols"` // Product grid columns (2-4)
	BorderRadius    string `bson:"borderRadius" json:"borderRadius"`       // Border radius (rem)
	Spacing         string `bson:"spacing" json:"spacing"`                 // Spacing scale: "compact", "normal", "relaxed"
	CardStyle       string `bson:"cardStyle" json:"cardStyle"`             // Card style: "flat", "bordered", "elevated"
	ButtonStyle     string `bson:"buttonStyle" json:"buttonStyle"`         // Button style: "rounded", "square", "pill"
	ShowBreadcrumbs bool   `bson:"showBreadcrumbs" json:"showBreadcrumbs"` // Show breadcrumbs
	ShowProductGrid bool   `bson:"showProductGrid" json:"showProductGrid"` // Show grid view option
}

// ThemeAssets represents uploaded assets
type ThemeAssets struct {
	Logo       string `bson:"logo,omitempty" json:"logo,omitempty"`             // Logo URL
	LogoWidth  int    `bson:"logoWidth,omitempty" json:"logoWidth,omitempty"`   // Logo width (px)
	LogoHeight int    `bson:"logoHeight,omitempty" json:"logoHeight,omitempty"` // Logo height (px)
	Favicon    string `bson:"favicon,omitempty" json:"favicon,omitempty"`       // Favicon URL
	HeroImage  string `bson:"heroImage,omitempty" json:"heroImage,omitempty"`   // Hero image URL
}

// Theme represents a complete storefront theme configuration
type Theme struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShopID      primitive.ObjectID `bson:"shop_id" json:"shop_id"` // Tenant isolation
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Version     string             `bson:"version" json:"version"`
	IsActive    bool               `bson:"isActive" json:"isActive"`

	// Theme configuration
	Colors     ThemeColors     `bson:"colors" json:"colors"`
	Typography ThemeTypography `bson:"typography" json:"typography"`
	Layout     ThemeLayout     `bson:"layout" json:"layout"`
	Assets     ThemeAssets     `bson:"assets" json:"assets"`

	// Dynamic Layout Blocks (Header/Footer)
	HeaderBlocks []Block `bson:"headerBlocks,omitempty" json:"headerBlocks,omitempty"`
	FooterBlocks []Block `bson:"footerBlocks,omitempty" json:"footerBlocks,omitempty"`

	// Legacy support (for backward compatibility)
	Tokens            map[string]interface{}            `bson:"tokens,omitempty" json:"tokens,omitempty"`
	ComponentDefaults map[string]map[string]interface{} `bson:"componentDefaults,omitempty" json:"componentDefaults,omitempty"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy string    `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
}

// DefaultTheme creates a theme with sensible defaults
func DefaultTheme(shopID primitive.ObjectID) *Theme {
	now := time.Now()
	return &Theme{
		ShopID:      shopID,
		Name:        "Default Theme",
		Description: "Default storefront theme",
		Version:     "1.0.0",
		IsActive:    true,
		Colors: ThemeColors{
			Primary:       "#10B981", // Green
			Secondary:     "#6366F1", // Indigo
			Accent:        "#F59E0B", // Amber
			Background:    "#FFFFFF", // White
			Surface:       "#F9FAFB", // Gray-50
			TextPrimary:   "#111827", // Gray-900
			TextSecondary: "#6B7280", // Gray-500
			Border:        "#E5E7EB", // Gray-200
			Success:       "#10B981", // Green
			Warning:       "#F59E0B", // Amber
			Error:         "#EF4444", // Red
			Info:          "#3B82F6", // Blue
		},
		Typography: ThemeTypography{
			HeadingFont:     "Inter",
			BodyFont:        "Inter",
			HeadingWeight:   "600",
			BodyWeight:      "400",
			HeadingSizeBase: "2rem",
			BodySizeBase:    "1rem",
			LineHeight:      "1.5",
			LetterSpacing:   "normal",
		},
		Layout: ThemeLayout{
			HeaderStyle:     "classic",
			ContainerWidth:  "boxed",
			ProductGridCols: 3,
			BorderRadius:    "0.5rem",
			Spacing:         "normal",
			CardStyle:       "elevated",
			ButtonStyle:     "rounded",
			ShowBreadcrumbs: true,
			ShowProductGrid: true,
		},
		Assets:    ThemeAssets{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
