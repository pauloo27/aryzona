package discord

import (
	"fmt"

	"github.com/pauloo27/aryzona/internal/config"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

func AsMention(userID string) string {
	return fmt.Sprintf("<@%s>", userID)
}

func OpenChatWithOwner() (model.TextChannel, error) {
	return Bot.OpenChannelWithUser(config.Config.OwnerID)
}

func DisableButtons(components []model.MessageComponent, selectedIndex int) []model.MessageComponent {
	disabledComponents := make([]model.MessageComponent, len(components))
	for i, component := range components {
		buttonComponent := component.(model.ButtonComponent)
		buttonComponent.Disabled = true
		if i != selectedIndex {
			buttonComponent.Style = model.SecondaryButtonStyle
		}
		disabledComponents[i] = buttonComponent
	}
	return disabledComponents
}
