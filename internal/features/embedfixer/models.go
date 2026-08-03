package embedfixer

import "gorm.io/gorm"

type EmbedFixerDomainOverride struct {
	gorm.Model
	GuildID      string `gorm:"not null;index:idx_embedfixer_guild_platform,unique"`
	Platform     string `gorm:"not null;index:idx_embedfixer_guild_platform,unique"`
	CustomDomain string `gorm:"not null"`
}
