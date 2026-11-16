package input

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func Slug(raw string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	var b strings.Builder
	prevHyphen := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevHyphen = false
			continue
		}
		if r == ' ' || r == '_' || r == '-' {
			if b.Len() > 0 && !prevHyphen {
				b.WriteByte('-')
				prevHyphen = true
			}
		}
	}
	slug := strings.Trim(b.String(), "-")
	if !ValidSlug(slug) {
		return "", errors.New("Use 3 to 40 letters, numbers, and hyphens")
	}
	return slug, nil
}

func ValidSlug(slug string) bool {
	return slugPattern.MatchString(slug) && len(slug) >= 3 && len(slug) <= 40
}

func Email(raw string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if len(value) > 254 {
		return "", errors.New("Enter a valid email address")
	}
	addr, err := mail.ParseAddress(value)
	if err != nil || addr.Address != value || !strings.Contains(value, ".") {
		return "", errors.New("Enter a valid email address")
	}
	return value, nil
}

func Text(raw, label string, min, max int) (string, error) {
	value := strings.TrimSpace(raw)
	if len([]rune(value)) < min {
		return "", errors.New(label + " is too short")
	}
	if len([]rune(value)) > max {
		return "", errors.New(label + " is too long")
	}
	return value, nil
}

func Optional(raw string, max int) (string, error) {
	value := strings.TrimSpace(raw)
	if len([]rune(value)) > max {
		return "", errors.New("That text is too long")
	}
	return value, nil
}

func URL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", nil
	}
	if len(value) > 400 || (!strings.HasPrefix(value, "https://") && !strings.HasPrefix(value, "http://")) {
		return "", errors.New("Use an http or https URL")
	}
	return value, nil
}
