package utils

import (
	"crypto/rand"
	"math/big"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/discord"
	"github.com/Pauloo27/aryzona/internal/utils"
)

var gif = "https://c.tenor.com/KzfDyYrsLngAAAAM/dice-roll.gif"

var RollCommand = command.Command{
	Name: "roll", Aliases: []string{"rolar", "dado", "dados", "dice", "rool", "d"},
	Description: "Roll a dice",
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
			return
		}
		luckyNumber := bigLuckyNumber.Int64() + 1

		embed := discord.NewEmbed().
			WithTitle(utils.Fmt(":game_die: You got ||  %d  || (click in the black box to reveal)", luckyNumber)).
			WithDescription(utils.Fmt("You rolled a %d sides\n_Gif by [Tenor](https://tenor.com/)_", sides)).
			WithImage(gif)

		ctx.SuccessEmbed(embed)
	},
}
