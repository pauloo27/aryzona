package tools

import (
	"fmt"
	"regexp"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/aryzona/internal/providers/doc"
)

var cpfMaskRe = regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})$`)

var CPFCommand = command.Command{
	Name: "cpf", Description: "Generate a CPF",
	Handler: func(ctx *command.CommandContext) {
		t := ctx.T.(*i18n.CommandCPF)

		cpf := doc.GenerateCPF()
		components := cpfMaskRe.FindStringSubmatch(cpf)
		maskedCPF := fmt.Sprintf(
			"%s.%s.%s-%s",
			components[1], components[2], components[3], components[4],
		)

		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle(t.Title.Str()).
				WithField("Without mask", cpf).
				WithField("With mask", maskedCPF),
		)
	},
}
