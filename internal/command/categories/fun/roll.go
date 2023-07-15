package fun

import (
	"errors"

	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/core/f"
	"github.com/pauloo27/aryzona/internal/core/rnd"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
	"github.com/pauloo27/aryzona/internal/providers/dice"
	"github.com/pauloo27/logger"
)

const (
	gif = "https://c.tenor.com/KzfDyYrsLngAAAAM/dice-roll.gif"
)

var RollCommand = command.Command{
	Name: "roll", Aliases: []string{"dice"},
	Parameters: []*command.Parameter{
		{Name: "faces", Required: false, Type: diceNotation},
	},
	Handler: func(ctx *command.Context) {
		t := ctx.T.(*i18n.CommandRoll)

		var d *dice.DiceNotation

		if len(ctx.Args) == 1 {
			d = ctx.Args[0].(*dice.DiceNotation)
		} else {
			d = dice.DefaultDice
		}

		numbers := make([]int, d.Dices)
		result := 0

		for i := 0; i < d.Dices; i++ {
			luckyNumber, err := rnd.Rnd(d.Faces)
			if err != nil {
				ctx.Error(t.SomethingWentWrong.Str())
				logger.Error(err)
				return
			}
			// +1 since the dice starts at 1
			luckyNumber++
			result += luckyNumber
			numbers[i] = luckyNumber
		}

		embed := model.NewEmbed().
			WithTitle(t.Title.Str(":game_die:", result)).
			WithDescription(
				t.Description.Str(
					d.String(),
					d.Dices, f.Pluralize(d.Dices, t.Dice.Str(), t.Dices.Str()),
					d.Faces, f.Pluralize(d.Faces, t.Face.Str(), t.Faces.Str()),
					numbers, result,
				),
			).
			WithImage(gif)

		ctx.SuccessEmbed(embed)
	},
}

var (
	// FIXME: i18n this
	diceNotation = &command.ParameterType{
		BaseType: parameters.TypeString,
		Name:     "dice notation",
		Parser: func(ctx *command.Context, index int, args []string) (any, error) {
			d, err := dice.ParseNotation(args[index])
			if err != nil {
				return nil, errors.New("invalid notation")
			}
			// discord limits the embed size, so we cant send a really
			// large thing (by large, i am not talking about the result,
			// but about the list of all rolled dices)
			if d.Dices > 200 {
				return nil, errors.New("dices cannot be more than 200")
			}
			if d.Faces > 200 {
				return nil, errors.New("faces cannot be more than 200")
			}
			return d, err
		},
	}
)
