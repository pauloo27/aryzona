package utils

import (
	"math/rand"
	"time"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/utils"
)

var DiceCommand = command.Command{
	Name: "dice", Aliases: []string{"rolar", "dado", "dados", "roll", "rool", "d"},
	Description: "Play a video/song from utube",
	Arguments: []*command.CommandArgument{
		{Name: "Dice sides", Required: false, Type: command.ArgumentInt},
	},
	Handler: func(ctx *command.CommandContext) {
		rand.Seed(time.Now().UnixMilli())
		var sides int
		if len(ctx.Args) != 0 {
			sides = ctx.Args[0].(int)
		}
		if sides == 0 {
			sides = 6
		}

		luckyNumber := rand.Intn(sides) + 1
		embed := utils.NewEmbedBuilder().
			Title(utils.Fmt(":game_die: You got ||%d|| (click in the black box to reveal)", luckyNumber)).
			Description(utils.Fmt("You rolled a %d sides", sides))
		ctx.SuccessEmbed(embed.Build())
	},
}
