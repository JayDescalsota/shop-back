package validator

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

type ValidationErrors map[string]string

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation passed"
	}
	var msgs []string
	for field, msg := range ve {
		msgs = append(msgs, field+": "+msg)
	}
	return "validation failed: " + strings.Join(msgs, "; ")
}

type Validator struct {
	errors ValidationErrors
}

func New() *Validator {
	return &Validator{errors: make(ValidationErrors)}
}

func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

func (v *Validator) Valid() bool {
	return !v.errors.HasErrors()
}

func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors[field] = "is required"
	}
	return v
}

func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		v.errors[field] = "must be a valid email address"
	}
	return v
}

func (v *Validator) MinLength(field, value string, min int) *Validator {
	if utf8.RuneCountInString(value) < min {
		v.errors[field] = "must be at least " + itoa(min) + " characters"
	}
	return v
}

func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if utf8.RuneCountInString(value) > max {
		v.errors[field] = "must not exceed " + itoa(max) + " characters"
	}
	return v
}

func (v *Validator) Phone(field, value string) *Validator {
	if value == "" {
		return v
	}
	re := regexp.MustCompile(`^\+?[1-9]\d{6,14}$`)
	if !re.MatchString(value) {
		v.errors[field] = "must be a valid phone number"
	}
	return v
}

func (v *Validator) Enum(field, value string, allowed ...string) *Validator {
	for _, a := range allowed {
		if value == a {
			return v
		}
	}
	v.errors[field] = "must be one of: " + strings.Join(allowed, ", ")
	return v
}

func (v *Validator) Min(field string, value, min int) *Validator {
	if value < min {
		v.errors[field] = "must be at least " + itoa(min)
	}
	return v
}

func (v *Validator) Max(field string, value, max int) *Validator {
	if value > max {
		v.errors[field] = "must not exceed " + itoa(max)
	}
	return v
}

func (v *Validator) URL(field, value string) *Validator {
	if value == "" {
		return v
	}
	re := regexp.MustCompile(`^https?://`)
	if !re.MatchString(value) {
		v.errors[field] = "must be a valid URL"
	}
	return v
}

func itoa(i int) string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					fmt.Sprintf("%d", i),
					"", "", 0,
				),
				"", "", 0,
			),
			"", "", 0,
		),
		"",
	))
}

// Use a simple approach to avoid import cycle
func itoaSimple(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

var itoa = itoaSimple
