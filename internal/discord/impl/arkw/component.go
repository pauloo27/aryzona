package arkw

import (
	"github.com/Pauloo27/aryzona/internal/discord/model"
	dc "github.com/diamondburned/arikawa/v3/discord"
)

func buildComponents(components []model.MessageComponent) []dc.InteractiveComponent {
	builtComponents := make([]dc.InteractiveComponent, len(components))

	for i, component := range components {
		switch t := component.(type) {
		case model.ButtonComponent:
			builtComponents[i] = &dc.ButtonComponent{
				Label:    t.Label,
				CustomID: dc.ComponentID(t.ID),
				Emoji:    &dc.ComponentEmoji{Name: t.Emoji},
				Style:    style(t.Style),
				Disabled: t.Disabled,
			}
		}
	}

	return builtComponents
}

func style(s model.ButtonStyle) dc.ButtonComponentStyle {
	switch s {
	case model.PrimaryButtonStyle:
		return dc.PrimaryButtonStyle()
	case model.SecondaryButtonStyle:
		return dc.SecondaryButtonStyle()
	case model.SuccessButtonStyle:
		return dc.SuccessButtonStyle()
	case model.DangerButtonStyle:
		return dc.DangerButtonStyle()
	}
	return dc.PrimaryButtonStyle()
}
