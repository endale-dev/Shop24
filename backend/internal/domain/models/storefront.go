package models

type StorefrontConfig struct {
	Version  int       `json:"version"`
	Theme    Theme     `json:"theme"`
	Header   Header    `json:"header"`
	Sections []Section `json:"sections"`
	Footer   Footer    `json:"footer"`
}

type Theme struct {
	Name           string `json:"name"`
	Colors         Colors `json:"colors"`
	FontHeading    string `json:"fontHeading"`
	FontBody       string `json:"fontBody"`
	Radius         string `json:"radius"`
	ImageRatio     string `json:"imageRatio"`
	Density        string `json:"density"`
	Buttons        string `json:"buttons"`
	Cards          string `json:"cards"`
	Collections    string `json:"collections"`
	HeroLayout     string `json:"heroLayout"`
	ProductDisplay string `json:"productDisplay"`
}

type Colors struct {
	Background string `json:"background"`
	Surface    string `json:"surface"`
	Text       string `json:"text"`
	Muted      string `json:"muted"`
	Primary    string `json:"primary"`
	OnPrimary  string `json:"onPrimary"`
	Accent     string `json:"accent"`
	Border     string `json:"border"`
}

type Header struct {
	Announcement     string    `json:"announcement"`
	ShowAnnouncement bool      `json:"showAnnouncement"`
	Alignment        string    `json:"alignment"`
	Nav              []NavItem `json:"nav"`
}

type NavItem struct {
	Label  string `json:"label"`
	Target string `json:"target"`
}

type Footer struct {
	Text  string    `json:"text"`
	Links []NavItem `json:"links"`
}

type Section struct {
	ID           string             `json:"id"`
	Type         string             `json:"type"`
	Visible      bool               `json:"visible"`
	Hero         *HeroProps         `json:"hero,omitempty"`
	ProductGrid  *ProductGridProps  `json:"productGrid,omitempty"`
	Featured     *FeaturedProps     `json:"featured,omitempty"`
	About        *AboutProps        `json:"about,omitempty"`
	Announcement *AnnouncementProps `json:"announcement,omitempty"`
	Categories   *CategoriesProps   `json:"categories,omitempty"`
	Newsletter   *NewsletterProps   `json:"newsletter,omitempty"`
}

type HeroProps struct {
	Eyebrow     string `json:"eyebrow"`
	Headline    string `json:"headline"`
	Subheadline string `json:"subheadline"`
	CTALabel    string `json:"ctaLabel"`
	CTATarget   string `json:"ctaTarget"`
	ImageURL    string `json:"imageUrl"`
	Align       string `json:"align"`
	Size        string `json:"size"`
}

type ProductGridProps struct {
	Title   string `json:"title"`
	Columns int    `json:"columns"`
	Limit   int    `json:"limit"`
}

type FeaturedProps struct {
	Title string `json:"title"`
	Limit int    `json:"limit"`
}

type AboutProps struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ImageURL string `json:"imageUrl"`
}

type AnnouncementProps struct {
	Text string `json:"text"`
}

type CategoriesProps struct {
	Title string `json:"title"`
}

type NewsletterProps struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
