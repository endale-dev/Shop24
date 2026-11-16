package mongorepo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"shop24/internal/domain/models"
)

type colorDoc struct {
	Background string `bson:"background"`
	Surface    string `bson:"surface"`
	Text       string `bson:"text"`
	Muted      string `bson:"muted"`
	Primary    string `bson:"primary"`
	OnPrimary  string `bson:"onPrimary"`
	Accent     string `bson:"accent"`
	Border     string `bson:"border"`
}

type themeDoc struct {
	Name        string   `bson:"name"`
	Colors      colorDoc `bson:"colors"`
	FontHeading string   `bson:"fontHeading"`
	FontBody    string   `bson:"fontBody"`
	Radius      string   `bson:"radius"`
	ImageRatio  string   `bson:"imageRatio"`
}

type navDoc struct {
	Label  string `bson:"label"`
	Target string `bson:"target"`
}

type headerDoc struct {
	Announcement     string   `bson:"announcement"`
	ShowAnnouncement bool     `bson:"showAnnouncement"`
	Alignment        string   `bson:"alignment"`
	Nav              []navDoc `bson:"nav"`
}

type footerDoc struct {
	Text  string   `bson:"text"`
	Links []navDoc `bson:"links"`
}

type heroDoc struct {
	Eyebrow     string `bson:"eyebrow"`
	Headline    string `bson:"headline"`
	Subheadline string `bson:"subheadline"`
	CTALabel    string `bson:"ctaLabel"`
	CTATarget   string `bson:"ctaTarget"`
	ImageURL    string `bson:"imageUrl"`
	Align       string `bson:"align"`
	Size        string `bson:"size"`
}

type gridDoc struct {
	Title   string `bson:"title"`
	Columns int    `bson:"columns"`
	Limit   int    `bson:"limit"`
}

type featuredDoc struct {
	Title string `bson:"title"`
	Limit int    `bson:"limit"`
}

type aboutDoc struct {
	Title    string `bson:"title"`
	Body     string `bson:"body"`
	ImageURL string `bson:"imageUrl"`
}

type announcementDoc struct {
	Text string `bson:"text"`
}

type categoriesDoc struct {
	Title string `bson:"title"`
}

type newsletterDoc struct {
	Title string `bson:"title"`
	Text  string `bson:"text"`
}

type sectionDoc struct {
	ID           string           `bson:"id"`
	Type         string           `bson:"type"`
	Visible      bool             `bson:"visible"`
	Hero         *heroDoc         `bson:"hero,omitempty"`
	ProductGrid  *gridDoc         `bson:"productGrid,omitempty"`
	Featured     *featuredDoc     `bson:"featured,omitempty"`
	About        *aboutDoc        `bson:"about,omitempty"`
	Announcement *announcementDoc `bson:"announcement,omitempty"`
	Categories   *categoriesDoc   `bson:"categories,omitempty"`
	Newsletter   *newsletterDoc   `bson:"newsletter,omitempty"`
}

type storefrontDoc struct {
	Version  int          `bson:"version"`
	Theme    themeDoc     `bson:"theme"`
	Header   headerDoc    `bson:"header"`
	Sections []sectionDoc `bson:"sections"`
	Footer   footerDoc    `bson:"footer"`
}

type designNoteDoc struct {
	Provider string `bson:"provider"`
	Mode     string `bson:"mode"`
	Prompt   string `bson:"prompt"`
	Summary  string `bson:"summary"`
}

type storeDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID     string             `bson:"ownerId"`
	Name        string             `bson:"name"`
	Slug        string             `bson:"slug"`
	Description string             `bson:"description"`
	LogoURL     string             `bson:"logoUrl"`
	Currency    string             `bson:"currency"`
	Published   bool               `bson:"published"`
	Storefront  storefrontDoc      `bson:"storefront"`
	Draft       storefrontDoc      `bson:"draft"`
	DesignNote  *designNoteDoc     `bson:"designNote,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

func navToDocs(items []models.NavItem) []navDoc {
	out := make([]navDoc, 0, len(items))
	for _, item := range items {
		out = append(out, navDoc{Label: item.Label, Target: item.Target})
	}
	return out
}

func navFromDocs(items []navDoc) []models.NavItem {
	out := make([]models.NavItem, 0, len(items))
	for _, item := range items {
		out = append(out, models.NavItem{Label: item.Label, Target: item.Target})
	}
	return out
}

func storefrontToDoc(cfg models.StorefrontConfig) storefrontDoc {
	sections := make([]sectionDoc, 0, len(cfg.Sections))
	for _, section := range cfg.Sections {
		doc := sectionDoc{ID: section.ID, Type: section.Type, Visible: section.Visible}
		if section.Hero != nil {
			hero := *section.Hero
			doc.Hero = &heroDoc{
				Eyebrow: hero.Eyebrow, Headline: hero.Headline, Subheadline: hero.Subheadline,
				CTALabel: hero.CTALabel, CTATarget: hero.CTATarget, ImageURL: hero.ImageURL,
				Align: hero.Align, Size: hero.Size,
			}
		}
		if section.ProductGrid != nil {
			grid := *section.ProductGrid
			doc.ProductGrid = &gridDoc{Title: grid.Title, Columns: grid.Columns, Limit: grid.Limit}
		}
		if section.Featured != nil {
			featured := *section.Featured
			doc.Featured = &featuredDoc{Title: featured.Title, Limit: featured.Limit}
		}
		if section.About != nil {
			about := *section.About
			doc.About = &aboutDoc{Title: about.Title, Body: about.Body, ImageURL: about.ImageURL}
		}
		if section.Announcement != nil {
			doc.Announcement = &announcementDoc{Text: section.Announcement.Text}
		}
		if section.Categories != nil {
			doc.Categories = &categoriesDoc{Title: section.Categories.Title}
		}
		if section.Newsletter != nil {
			doc.Newsletter = &newsletterDoc{Title: section.Newsletter.Title, Text: section.Newsletter.Text}
		}
		sections = append(sections, doc)
	}
	return storefrontDoc{
		Version: cfg.Version,
		Theme: themeDoc{
			Name: cfg.Theme.Name,
			Colors: colorDoc{
				Background: cfg.Theme.Colors.Background, Surface: cfg.Theme.Colors.Surface,
				Text: cfg.Theme.Colors.Text, Muted: cfg.Theme.Colors.Muted, Primary: cfg.Theme.Colors.Primary,
				OnPrimary: cfg.Theme.Colors.OnPrimary, Accent: cfg.Theme.Colors.Accent, Border: cfg.Theme.Colors.Border,
			},
			FontHeading: cfg.Theme.FontHeading, FontBody: cfg.Theme.FontBody,
			Radius: cfg.Theme.Radius, ImageRatio: cfg.Theme.ImageRatio,
		},
		Header: headerDoc{
			Announcement: cfg.Header.Announcement, ShowAnnouncement: cfg.Header.ShowAnnouncement,
			Alignment: cfg.Header.Alignment, Nav: navToDocs(cfg.Header.Nav),
		},
		Sections: sections,
		Footer:   footerDoc{Text: cfg.Footer.Text, Links: navToDocs(cfg.Footer.Links)},
	}
}

func storefrontFromDoc(doc storefrontDoc) models.StorefrontConfig {
	sections := make([]models.Section, 0, len(doc.Sections))
	for _, section := range doc.Sections {
		item := models.Section{ID: section.ID, Type: section.Type, Visible: section.Visible}
		if section.Hero != nil {
			item.Hero = &models.HeroProps{
				Eyebrow: section.Hero.Eyebrow, Headline: section.Hero.Headline, Subheadline: section.Hero.Subheadline,
				CTALabel: section.Hero.CTALabel, CTATarget: section.Hero.CTATarget, ImageURL: section.Hero.ImageURL,
				Align: section.Hero.Align, Size: section.Hero.Size,
			}
		}
		if section.ProductGrid != nil {
			item.ProductGrid = &models.ProductGridProps{Title: section.ProductGrid.Title, Columns: section.ProductGrid.Columns, Limit: section.ProductGrid.Limit}
		}
		if section.Featured != nil {
			item.Featured = &models.FeaturedProps{Title: section.Featured.Title, Limit: section.Featured.Limit}
		}
		if section.About != nil {
			item.About = &models.AboutProps{Title: section.About.Title, Body: section.About.Body, ImageURL: section.About.ImageURL}
		}
		if section.Announcement != nil {
			item.Announcement = &models.AnnouncementProps{Text: section.Announcement.Text}
		}
		if section.Categories != nil {
			item.Categories = &models.CategoriesProps{Title: section.Categories.Title}
		}
		if section.Newsletter != nil {
			item.Newsletter = &models.NewsletterProps{Title: section.Newsletter.Title, Text: section.Newsletter.Text}
		}
		sections = append(sections, item)
	}
	return models.StorefrontConfig{
		Version: doc.Version,
		Theme: models.Theme{
			Name: doc.Theme.Name,
			Colors: models.Colors{
				Background: doc.Theme.Colors.Background, Surface: doc.Theme.Colors.Surface,
				Text: doc.Theme.Colors.Text, Muted: doc.Theme.Colors.Muted, Primary: doc.Theme.Colors.Primary,
				OnPrimary: doc.Theme.Colors.OnPrimary, Accent: doc.Theme.Colors.Accent, Border: doc.Theme.Colors.Border,
			},
			FontHeading: doc.Theme.FontHeading, FontBody: doc.Theme.FontBody,
			Radius: doc.Theme.Radius, ImageRatio: doc.Theme.ImageRatio,
		},
		Header: models.Header{
			Announcement: doc.Header.Announcement, ShowAnnouncement: doc.Header.ShowAnnouncement,
			Alignment: doc.Header.Alignment, Nav: navFromDocs(doc.Header.Nav),
		},
		Sections: sections,
		Footer:   models.Footer{Text: doc.Footer.Text, Links: navFromDocs(doc.Footer.Links)},
	}
}

func noteToDoc(note models.DesignNote) designNoteDoc {
	return designNoteDoc{Provider: note.Provider, Mode: note.Mode, Prompt: note.Prompt, Summary: note.Summary}
}

func noteFromDoc(doc *designNoteDoc) *models.DesignNote {
	if doc == nil {
		return nil
	}
	return &models.DesignNote{Provider: doc.Provider, Mode: doc.Mode, Prompt: doc.Prompt, Summary: doc.Summary}
}

func (d storeDoc) model() models.Store {
	return models.Store{
		ID: d.ID.Hex(), OwnerID: d.OwnerID, Name: d.Name, Slug: d.Slug, Description: d.Description,
		LogoURL: d.LogoURL, Currency: d.Currency, Published: d.Published,
		Storefront: storefrontFromDoc(d.Storefront), Draft: storefrontFromDoc(d.Draft),
		DesignNote: noteFromDoc(d.DesignNote), CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
	}
}
