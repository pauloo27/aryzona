package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/Pauloo27/aryzona/command"
	"github.com/Pauloo27/aryzona/command/slash"
	"github.com/Pauloo27/aryzona/discord"
	"github.com/Pauloo27/aryzona/git"
	"github.com/Pauloo27/aryzona/utils"
	"github.com/Pauloo27/logger"
	"github.com/joho/godotenv"

	// import listeners
	_ "github.com/Pauloo27/aryzona/discord/listeners"

	// import scheduler
	_ "github.com/Pauloo27/aryzona/utils/scheduler"

	// import all command categories
	_ "github.com/Pauloo27/aryzona/command/categories/audio"
	_ "github.com/Pauloo27/aryzona/command/categories/sysmon"
	_ "github.com/Pauloo27/aryzona/command/categories/utils"
)

var commitHash, commitMessage string

func init() {
	logger.Info("Loading .env...")
	err := godotenv.Load()
	logger.HandleFatal(err, "Cannot load .env")
	logger.Success(".env loaded")

	git.CommitHash = commitHash
	git.CommitMessage = commitMessage
	git.RemoteRepo = os.Getenv("DC_BOT_REMOTE_REPO")
}

func init() {
	logger.AddLogListener(func(level logger.Level, params ...interface{}) {
		if discord.Session == nil || (level != logger.ERROR && level != logger.FATAL) {
			return
		}

		c, err := utils.OpenChatWithOwner(discord.Session)
		if err != nil {
			// to avoid loops, do not call the logger again
			fmt.Println("shit 0", err)
			return
		}

		if err != nil {
			fmt.Println("shit 1", err)
		}

		embed := utils.NewEmbedBuilder().
			FieldInline("Message", fmt.Sprintln(params...)).
			Description(utils.Fmt("```\n%s\n```", string(debug.Stack()))).
			Color(0xff5555).
			Title(utils.Fmt("Oops! [%s]", level.Name))

		_, err = discord.Session.ChannelMessageSendEmbed(c.ID, embed.Build())
		if err != nil {
			fmt.Println("shit 2", err)
			return
		}
	})
}

func main() {
	logger.Info("Connecting to Discord...")
	err := discord.Create(os.Getenv("DC_BOT_TOKEN"))
	if err != nil {
		logger.Fatal(err)
	}

	discord.RegisterListeners()
	err = discord.Connect()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Success("Connected to discord")

	command.Prefix = os.Getenv("DC_BOT_PREFIX")

	logger.Info("Registering slash commands handlers...")
	err = slash.RegisterCommands(false)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Success("Slash commands created!")

	stop := make(chan os.Signal, 1)
	//lint:ignore SA1016 i dont know, it just works lol
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-stop
	err = discord.Disconnect()
	if err != nil {
		logger.Error("Cannot disconnect... we are disconnecting anyway...", err)
	}
	logger.Success("Exiting...")
}
