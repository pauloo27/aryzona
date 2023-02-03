package tools

import (
	"fmt"
	"regexp"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/aryzona/internal/providers/doc"
)

var cnpjMaskRe = regexp.MustCompile(`^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$`)

var CNPJCommand = command.Command{
	Name: "cnpj", Description: "Generate a CNPJ",
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandCNPJ)

		cnpj := doc.GenerateCNPJ()
		components := cnpjMaskRe.FindStringSubmatch(cnpj)
		maskedCNPJ := fmt.Sprintf(
			"%s.%s.%s/%s-%s",
			components[1], components[2], components[3], components[4], components[5],
		)

		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle(t.Title.Str()).
				WithField(t.WithoutMask.Str(), cnpj).
				WithField(t.WithoutMask.Str(), maskedCNPJ),
		)
	},
}
