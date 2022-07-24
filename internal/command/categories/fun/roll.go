package fun

import (
	"errors"
	"fmt"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/providers/dice"
	"github.com/Pauloo27/aryzona/internal/utils"
	"github.com/Pauloo27/aryzona/internal/utils/rnd"
)

const (
	gif = "https://c.tenor.com/KzfDyYrsLngAAAAM/dice-roll.gif"
)

var RollCommand = command.Command{
	Name: "roll", Aliases: []string{"rolar", "dado", "dados", "dice", "rool", "d"},
	Description: "Roll a dice",
	Parameters: []*command.CommandParameter{
		{Name: "sides", Description: "dice sides", Required: false, Type: diceNotation},
	},
	Handler: func(ctx *command.CommandContext) {
		var d *dice.DiceNotation

		if len(ctx.Args) == 1 {
			d = ctx.Args[0].(*dice.DiceNotation)
		} else {
			d = dice.DefaultDice
		}

		numbers := make([]int, d.Dices)
		result := 0

		for i := 0; i < d.Dices; i++ {
			luckyNumber, err := rnd.Rnd(d.Sides)
			if err != nil {
				ctx.Error("something went wrong =(")
				return
			}
			// +1 since the dice starts at 1
			luckyNumber++
			result += luckyNumber
			numbers[i] = luckyNumber
		}

		embed := discord.NewEmbed().
			WithTitle(fmt.Sprintf(":game_die: You got %d", result)).
			WithDescription(
				fmt.Sprintf("You rolled %s (%d %s with %d %s)\n%v -> %d\n_Gif by [Tenor](https://tenor.com/)_",
					d.String(),
					d.Dices, utils.Pluralize(d.Dices, "dice", "dices"),
					d.Sides, utils.Pluralize(d.Sides, "side", "sides"),
					numbers, result),
			).
			WithImage(gif)

		ctx.SuccessEmbed(embed)
	},
}

var (
	diceNotation = &command.CommandParameterType{
		BaseType: parameters.TypeString,
		Name:     "dice notation",
		Parser: func(index int, args []string) (interface{}, error) {
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
			if d.Sides > 200 {
				return nil, errors.New("sides cannot be more than 200")
			}
			return d, err
		},
	}
)
