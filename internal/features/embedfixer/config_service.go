package embedfixer

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bachacode/gatoc/internal/database"
	"gorm.io/gorm"
)

func normalizeDomainHost(domain string) (string, error) {
	value := strings.TrimSpace(strings.ToLower(domain))
	if value == "" {
		return "", fmt.Errorf("domain is required")
	}

	if strings.Contains(value, "://") {
		parsed, err := url.Parse(value)
		if err != nil {
			return "", fmt.Errorf("invalid domain")
		}
		if parsed.Path != "" && parsed.Path != "/" {
			return "", fmt.Errorf("domain must not include path, query, or fragment")
		}
		if parsed.RawQuery != "" || parsed.Fragment != "" {
			return "", fmt.Errorf("domain must not include path, query, or fragment")
		}

		value = parsed.Hostname()
	} else if strings.Contains(value, "/") || strings.Contains(value, "?") || strings.Contains(value, "#") {
		return "", fmt.Errorf("domain must not include path, query, or fragment")
	}

	value = strings.TrimPrefix(value, "www.")
	if value == "" || !strings.Contains(value, ".") || strings.Contains(value, " ") {
		return "", fmt.Errorf("invalid domain host")
	}

	return value, nil
}

func getCustomDomain(db *gorm.DB, guildID string, platformKey string) (string, bool, error) {
	if db == nil {
		return "", false, nil
	}

	var override database.EmbedFixerDomainOverride
	err := db.Where("guild_id = ? AND platform = ?", guildID, platformKey).First(&override).Error
	if err == nil {
		return override.CustomDomain, true, nil
	}

	if err == gorm.ErrRecordNotFound {
		return "", false, nil
	}

	return "", false, err
}

func setCustomDomain(db *gorm.DB, guildID string, platformKey string, customDomain string) error {
	if db == nil {
		return fmt.Errorf("database is not configured")
	}

	override := database.EmbedFixerDomainOverride{
		GuildID:      guildID,
		Platform:     platformKey,
		CustomDomain: customDomain,
	}

	return db.Where("guild_id = ? AND platform = ?", guildID, platformKey).
		Assign(database.EmbedFixerDomainOverride{CustomDomain: customDomain}).
		FirstOrCreate(&override).Error
}

func resetCustomDomain(db *gorm.DB, guildID string, platformKey string) error {
	if db == nil {
		return fmt.Errorf("database is not configured")
	}

	return db.Where("guild_id = ? AND platform = ?", guildID, platformKey).Delete(&database.EmbedFixerDomainOverride{}).Error
}

func activeDomain(platform platformConfig, customDomain string, isCustom bool) (string, bool) {
	if isCustom {
		return customDomain, true
	}

	return platform.DefaultDomain, false
}

func resolveReplacementDomain(db *gorm.DB, guildID string, platform platformConfig) (string, bool, error) {
	customDomain, isCustom, err := getCustomDomain(db, guildID, platform.Key)
	if err != nil {
		return "", false, err
	}

	domain, usingCustom := activeDomain(platform, customDomain, isCustom)
	return domain, usingCustom, nil
}
