package tools

import (
	"fmt"
	"regexp"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/utils"
)

var cpfMaskRe = regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})$`)

var CPFCommand = command.Command{
	Name: "cpf", Description: "Generate a CPF",
	Handler: func(ctx *command.CommandContext) {
		cpf := utils.GenerateCPF()
		components := cpfMaskRe.FindStringSubmatch(cpf)
		maskedCPF := fmt.Sprintf(
			"%s.%s.%s-%s",
			components[1], components[2], components[3], components[4],
		)

		ctx.SuccessEmbed(
			model.NewEmbed().
				WithTitle("CPF").
				WithField("Without mask", cpf).
				WithField("With mask", maskedCPF),
		)
	},
}
