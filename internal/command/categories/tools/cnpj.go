package tools

import (
	"fmt"
	"regexp"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/doc"
)

var cnpjMaskRe = regexp.MustCompile(`^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$`)

var CNPJCommand = command.Command{
	Name: "cnpj",
	Handler: func(ctx *command.Context) command.Result {
		t := ctx.T.(*i18n.CommandCNPJ)

		cnpj := doc.GenerateCNPJ()
		components := cnpjMaskRe.FindStringSubmatch(cnpj)
		maskedCNPJ := fmt.Sprintf(
			"%s.%s.%s/%s-%s",
			components[1], components[2], components[3], components[4], components[5],
		)

		return ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle(t.Title.Str()).
				WithField(t.WithoutMask.Str(), cnpj).
				WithField(t.WithMask.Str(), maskedCNPJ),
		)
	},
}
