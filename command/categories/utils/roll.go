package utils

import (
	"crypto/rand"
	"math/big"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/utils"
)

var gif = "https://c.tenor.com/KzfDyYrsLngAAAAM/dice-roll.gif"

var RollCommand = command.Command{
	Name: "roll", Aliases: []string{"rolar", "dado", "dados", "dice", "rool", "d"},
	Description: "Play a video/song from utube",
	Arguments: []*command.CommandArgument{
		{Name: "sides", Description: "dice sides", Required: false, Type: command.ArgumentInt},
	},
	Handler: func(ctx *command.CommandContext) {
		var sides int
		if len(ctx.Args) != 0 {
			sides = ctx.Args[0].(int)
		}
		if sides <= 0 {
			sides = 6
		}

		/*
			ATTENTION | ATENÇÃO | ATENCIÓN | ATTENZIONE | ATENTO | ANIMADVERSIO

			UGLY CODE AHEAD!
		*/

		bigLuckyNumber, err := rand.Int(rand.Reader, big.NewInt(int64(sides)))
		if err != nil {
			ctx.Error("something went wrong =(")
		}
		luckyNumber := bigLuckyNumber.Int64() + 1

		embed := utils.NewEmbedBuilder().
			Title(utils.Fmt(":game_die: You got ||  %d  || (click in the black box to reveal)", luckyNumber)).
			Description(utils.Fmt("You rolled a %d sides\n_Gif by [Tenor](https://tenor.com/)_", sides)).
			Image(gif)

		ctx.SuccessEmbed(embed.Build())
	},
}
