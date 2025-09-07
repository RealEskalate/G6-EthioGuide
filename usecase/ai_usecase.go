package usecase

import (
	"EthioGuide/domain"
	"context"
	"fmt"
	"strings"
	"time"
)

var supportedLangs = map[string]bool{
	"en": true, // English
	"am": true, // Amharic
	// Add other supported languages here
}

// nonTranslatableKeys defines a set of JSON keys whose string values should NOT be translated.
// This is more reliable than guessing based on content.
var nonTranslatableKeys = map[string]bool{
	// Identifiers
	"id":                true,
	"ID":                true,
	"groupId":           true,
	"GroupID":           true,
	"organizationId":    true,
	"OrganizationID":    true,
	"organization_id":   true,
	"parent_id":         true,
	"userId":            true,
	"UserID":            true,
	"user_id":           true,
	"procedure_id":      true,
	"ProcedureID":       true,
	"user_procedure_id": true,
	"checklist_id":      true,
	"noticeIds":         true,

	// User account & auth fields
	"username":        true,
	"email":           true,
	"phone":           true,
	"password":        true,
	"old_password":    true,
	"new_password":    true,
	"preferredLang":   true,
	"provider":        true,
	"code":            true,
	"role":            true,
	"access_token":    true,
	"refresh_token":   true,
	"resetToken":      true,
	"activationToken": true,

	// URLs and file paths
	"profile_picture": true,
	"profilePicURL":   true,
	"profile_pic_url": true,
	"website":         true,

	// Enumerated types or statuses
	"status":   true,
	"type":     true,
	"currency": true,

	// Timestamps and technical metadata
	"created_at": true,
	"CreatedAt":  true,
	"updated_at": true,
	"UpdatedAt":  true,
	"source":     true,
}

type geminiUseCase struct {
	geminiServices domain.IAIService
	contextTimeout time.Duration
}

func NewGeminiUsecase(geminiServices domain.IAIService, timeOut time.Duration) domain.IGeminiUseCase {
	return &geminiUseCase{
		geminiServices: geminiServices,
		contextTimeout: timeOut,
	}
}

// TranslateJSON translates all eligible string values within a JSON object based on a key denylist.
func (g *geminiUseCase) TranslateJSON(ctx context.Context, data map[string]interface{}, targetLang string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, g.contextTimeout)
	defer cancel()

	if !supportedLangs[targetLang] {
		return nil, fmt.Errorf("%w: %s", domain.ErrUnsupportedLanguage, targetLang)
	}

	// === PASS 1: Collect all strings from keys that are allowed to be translated ===
	stringMap := make(map[string]bool)
	collectStringsByKey(data, stringMap)

	if len(stringMap) == 0 {
		return data, nil // Nothing to translate
	}

	var originals []string
	for s := range stringMap {
		originals = append(originals, s)
	}

	// === Translate in a single batch call ===
	const separator = "<!--EthioGuideTranslationSeparator-->"
	contentToTranslate := strings.Join(originals, separator)

	prompt := fmt.Sprintf(`
You are a machine that translates text segments.
You will be given a block of text containing one or more segments separated by "%s".
Translate each segment into the language with the code '%s'.
Your response MUST contain the exact same number of "%s" separators as the input.
Do not add or remove separators.
Do not add any introductory text, explanations, markdown, or any text other than the translated segments and their separators.

Example:
Input Text: "Hello world|||How are you?"
Your Response for target language 'fr': "Bonjour le monde|||Comment ça va?"

Now, perform the translation for the following text:
%s
`, separator, targetLang, separator, separator, contentToTranslate)

	translatedBlock, err := g.geminiServices.GenerateCompletion(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("gemini service failed to generate completion: %w", err)
	}

	translatedParts := strings.Split(translatedBlock, separator)
	if len(translatedParts) != len(originals) {
		return nil, fmt.Errorf("%w: expected %d parts, got %d. AI failed to follow instructions", domain.ErrTranslationMismatch, len(originals), len(translatedParts))
	}

	translationMap := make(map[string]string)
	for i, originalStr := range originals {
		translationMap[originalStr] = strings.TrimSpace(translatedParts[i])
	}

	// === PASS 2: Replace strings in the original structure ===
	translatedData := replaceStringsByKey(data, translationMap).(map[string]interface{})

	return translatedData, nil
}

// collectStringsByKey recursively traverses the JSON, collecting strings from non-excluded keys.
func collectStringsByKey(node interface{}, stringMap map[string]bool) {
	switch v := node.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if nonTranslatableKeys[key] {
				continue // Skip this entire branch
			}
			// For maps within maps (like 'steps'), we don't pass the parent key.
			// The check is on the immediate parent key of the string value.
			collectStringsByKey(value, stringMap)
		}
	case []interface{}:
		for _, element := range v {
			collectStringsByKey(element, stringMap)
		}
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			stringMap[trimmed] = true
		}
	}
}

// replaceStringsByKey recursively traverses the JSON, replacing strings from non-excluded keys.
func replaceStringsByKey(node interface{}, translationMap map[string]string) interface{} {
	switch v := node.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{}, len(v))
		for key, value := range v {
			if nonTranslatableKeys[key] {
				newMap[key] = value // Keep original value and skip recursion
				continue
			}
			newMap[key] = replaceStringsByKey(value, translationMap)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for i, element := range v {
			newSlice[i] = replaceStringsByKey(element, translationMap)
		}
		return newSlice
	case string:
		trimmed := strings.TrimSpace(v)
		if translatedStr, ok := translationMap[trimmed]; ok {
			return translatedStr // Return the translation
		}
		return v // Return original string if no translation is found (should be rare)
	default:
		return v // For numbers, booleans, null, etc.
	}
}
