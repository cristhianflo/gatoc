package embedfixer

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestHasNoFixTag(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{name: "tag lowercase", content: "https://x.com/foo/status/1 #nofix", expected: true},
		{name: "tag uppercase", content: "https://x.com/foo/status/1 #NOFIX", expected: true},
		{name: "tag as token only", content: "#nofix", expected: true},
		{name: "no tag", content: "https://x.com/foo/status/1", expected: false},
		{name: "partial token", content: "https://x.com/foo/status/1 #nofixplease", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasNoFixTag(tc.content); got != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestEmbedFixerCommandIncludesNoFixSubcommand(t *testing.T) {
	found := false
	for _, option := range embedFixerConfigCommand.Metadata.Options {
		if option.Name == "nofix" && option.Type == discordgo.ApplicationCommandOptionSubCommand {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("expected embedfixer command to include nofix subcommand")
	}
}
