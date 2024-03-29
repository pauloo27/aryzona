package arkw

import (
	"fmt"

	dc "github.com/diamondburned/arikawa/v3/discord"
	"github.com/pauloo27/aryzona/internal/discord/model"
)

func buildRows(rows []model.MessageComponentRow) dc.ContainerComponents {
	builtRows := make(dc.ContainerComponents, len(rows))

	for i, row := range rows {
		rawComponents := buildComponents(row.Components)
		row := dc.ActionRowComponent(rawComponents)
		builtRows[i] = &row
	}

	return builtRows
}

func buildComponents(components []model.MessageComponent) []dc.InteractiveComponent {
	builtComponents := make([]dc.InteractiveComponent, len(components))

	for i, component := range components {
		switch t := component.(type) {
		case model.ButtonComponent:
			var emoji *dc.ComponentEmoji
			if t.Emoji != "" {
				emoji = &dc.ComponentEmoji{Name: t.Emoji}
			}
			id := fmt.Sprintf("%s:%s", t.BaseID, t.ID)
			builtComponents[i] = &dc.ButtonComponent{
				Label:    t.Label,
				CustomID: dc.ComponentID(id),
				Emoji:    emoji,
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
