package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"shop24/internal/domain/models"
	"shop24/internal/domain/storefront"
)

const systemPrompt = `You design storefront configuration for a small online shop.
Return one JSON object with version, theme, header, sections, and footer. Do not return HTML, JavaScript, Vue, CSS, or any executable code.
Section types: hero, product_grid, featured_products, about, announcement, categories, newsletter.
Each section has id, type, visible, and exactly one matching object: hero, productGrid, featured, about, announcement, categories, or newsletter.
Theme name: editorial, minimal, luxury, bold, playful, mono, or dark.
Fonts: Instrument Serif, Fraunces, Newsreader, Outfit, Figtree, Literata.
Colors are #rrggbb. Radius: none, sm, md, lg. Image ratio: square, portrait, landscape.
Theme layout, each a string: density tight|regular|roomy, buttons fill|outline|link, cards flat|border|shadow, collections strip|feature|hidden, heroLayout stack|split|bleed, productDisplay grid|feature|list.
Use those layout fields when the prompt asks for a different structure. A full-bleed hero is heroLayout bleed. Products in a list is productDisplay list. Hiding the collections block is collections hidden. Outline buttons are buttons outline. More space is density roomy.
Header alignment is left or center. Nav targets are home, products, about, or cart.
Hero size is sm, md, or lg. Hero align is left or center. Hero ctaTarget uses the nav targets.
Do not invent products or a working email form. Product sections are filled from the live catalog.
When revising, change only what the prompt asks for and keep the rest.
history, when present, is the earlier customization conversation. Keep those choices unless this prompt changes them.
Write concise shop copy that fits on a phone. version, columns, and limit must be JSON numbers. visible and showAnnouncement must be JSON booleans.
Use these field names exactly. theme is an object with name, not a string. A hero uses hero.headline and hero.subheadline. An about section uses about.title and about.body.
{"version":1,"theme":{"name":"minimal","density":"tight","buttons":"link","cards":"flat","collections":"strip","heroLayout":"stack","productDisplay":"grid"},"header":{"announcement":"","showAnnouncement":false,"alignment":"left","nav":[{"label":"Shop","target":"products"}]},"sections":[{"id":"hero","type":"hero","visible":true,"hero":{"eyebrow":"Independent shop","headline":"Shop name","subheadline":"Short line about the shop.","ctaLabel":"Shop","ctaTarget":"products","imageUrl":"","align":"left","size":"md"}},{"id":"grid","type":"product_grid","visible":true,"productGrid":{"title":"The collection","columns":2,"limit":8}}],"footer":{"text":"© Shop name","links":[]}}`

type OpenAI struct {
	key     string
	baseURL string
	model   string
	client  *http.Client
}

func NewOpenAI(key, baseURL, model string, timeout time.Duration) *OpenAI {
	return &OpenAI{
		key: key, baseURL: strings.TrimRight(baseURL, "/"), model: model,
		client: &http.Client{Timeout: timeout},
	}
}

func (o *OpenAI) Name() string {
	if strings.Contains(o.baseURL, "openrouter.ai") {
		return "openrouter"
	}
	return "openai"
}

func (o *OpenAI) Generate(ctx context.Context, prompt string, store StoreContext) (Result, error) {
	return o.complete(ctx, prompt, nil, store, "generate")
}

func (o *OpenAI) Refine(ctx context.Context, prompt string, current models.StorefrontConfig, store StoreContext) (Result, error) {
	return o.complete(ctx, prompt, &current, store, "refine")
}

func (o *OpenAI) complete(ctx context.Context, prompt string, current *models.StorefrontConfig, store StoreContext, mode string) (Result, error) {
	user, err := json.Marshal(map[string]any{
		"mode": mode, "prompt": prompt, "store": store, "current": current,
	})
	if err != nil {
		return Result{}, err
	}
	body, err := json.Marshal(map[string]any{
		"model":           o.model,
		"temperature":     0.4,
		"response_format": map[string]string{"type": "json_object"},
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": string(user)},
		},
	})
	if err != nil {
		return Result{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.key)
	if strings.Contains(o.baseURL, "openrouter.ai") {
		req.Header.Set("X-Title", "Sable")
	}
	res, err := o.client.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer res.Body.Close()
	payload, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return Result{}, err
	}
	if res.StatusCode >= 300 {
		return Result{}, &ProviderError{Detail: providerDetail(res.StatusCode, payload)}
	}
	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return Result{}, err
	}
	if len(parsed.Choices) == 0 || strings.TrimSpace(parsed.Choices[0].Message.Content) == "" {
		return Result{}, fmt.Errorf("design provider returned an empty response")
	}
	raw, err := prepareConfig(extractJSON(parsed.Choices[0].Message.Content))
	if err != nil {
		return Result{}, &ProviderError{Detail: "The designer returned a layout that could not be used. Try again."}
	}
	var cfg models.StorefrontConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return Result{}, &ProviderError{Detail: "The designer returned a layout that could not be used. Try again."}
	}
	return Result{
		Config:   storefront.Normalize(cfg),
		Provider: o.Name(),
		Summary:  "Designed by " + o.Name() + " and checked against the storefront schema.",
	}, nil
}

func providerDetail(status int, payload []byte) string {
	var body struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.Unmarshal(payload, &body)
	message := strings.TrimSpace(body.Error.Message)
	if message == "" {
		return fmt.Sprintf("The design provider returned status %d", status)
	}
	if strings.Contains(strings.ToLower(message), "sk-") {
		return fmt.Sprintf("The design provider returned status %d", status)
	}
	return message
}

func prepareConfig(content string) ([]byte, error) {
	var raw any
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return nil, err
	}
	fixConfigTypes(raw)
	return json.Marshal(raw)
}

func fixConfigTypes(value any) {
	switch node := value.(type) {
	case map[string]any:
		if name, ok := node["theme"].(string); ok {
			node["theme"] = map[string]any{"name": name}
		}
		copyField(node, "headline", "text", "title", "heading")
		copyField(node, "subheadline", "subtitle", "description")
		copyField(node, "body", "text", "description")
		for key, item := range node {
			if key == "version" || key == "columns" || key == "limit" {
				if number, ok := asInt(item); ok {
					node[key] = number
					continue
				}
			}
			if key == "visible" || key == "showAnnouncement" {
				if flag, ok := asBool(item); ok {
					node[key] = flag
					continue
				}
			}
			fixConfigTypes(item)
		}
	case []any:
		for _, item := range node {
			fixConfigTypes(item)
		}
	}
}

func copyField(node map[string]any, dest string, sources ...string) {
	if text, ok := node[dest].(string); ok && strings.TrimSpace(text) != "" {
		return
	}
	for _, source := range sources {
		text, ok := node[source].(string)
		if ok && strings.TrimSpace(text) != "" && source != dest {
			node[dest] = text
			return
		}
	}
}

func asInt(value any) (int, bool) {
	switch number := value.(type) {
	case float64:
		return int(number), true
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(number))
		return parsed, err == nil
	default:
		return 0, false
	}
}

func asBool(value any) (bool, bool) {
	switch flag := value.(type) {
	case bool:
		return flag, true
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(flag))
		return parsed, err == nil
	default:
		return false, false
	}
}

func extractJSON(content string) string {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "```") {
		return content
	}
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	if end := strings.LastIndex(content, "```"); end >= 0 {
		content = content[:end]
	}
	return strings.TrimSpace(content)
}
