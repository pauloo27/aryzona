package bot

import (
	"fmt"
	"strings"

	"github.com/Pauloo27/aryzona/internal/command"
	"github.com/Pauloo27/aryzona/internal/command/parameters"
	"github.com/Pauloo27/aryzona/internal/db/services"
	"github.com/Pauloo27/aryzona/internal/discord/model"
	"github.com/Pauloo27/aryzona/internal/i18n"
	"github.com/Pauloo27/logger"
)

var LanguageCommand = command.Command{
	Name: "language", Aliases: []string{"lang", "locale"},
	Parameters: []*command.CommandParameter{
		{
			Name: "language", Type: parameters.ParameterLowerCasedString,
			ValidValuesFunc: listValidLanguages,
			Required:        false,
		},
	},
	Handler: func(ctx *command.CommandContext) {
		if len(ctx.Args) == 0 {
			listLanguages(ctx)
		} else {
			selectLanguage(ctx)
		}
	},
}

func listLanguages(ctx *command.CommandContext) {
	t := ctx.T.(*i18n.CommandLanguage)

	var validLanguages strings.Builder

	for i, lang := range i18n.LanguagesName {
		if i != 0 {
			validLanguages.WriteString(", ")
		}
		validLanguages.WriteString(string(lang))
	}

	description := fmt.Sprintf(
		"%s\n\n%s", t.CurrentLanguage.Str(
			t.Name,
			t.DisplayName.Str(),
			t.Authors.Str(),
		),
		t.LanguageList.Str(validLanguages.String()),
	)

	embed := model.NewEmbed().
		WithTitle(t.Title.Str()).
		WithDescription(description)

	ctx.SuccessEmbed(embed)
}

func selectLanguage(ctx *command.CommandContext) {
	t := ctx.T.(*i18n.CommandLanguage)

	langName := ctx.Args[0].(string)

	var lang *i18n.Language

	for _, l := range i18n.LanguagesName {
		if strings.ToLower(string(l)) == langName {
			lang, _ = i18n.GetLanguage(l)
			break
		}
	}

	err := services.User.SetPreferredLang(ctx.AuthorID, lang.Name)
	if err != nil {
		logger.Error(err)
		ctx.Error(t.SomethingWentWrong.Str())
		return
	}

	newLang, _ := i18n.GetLanguage(lang.Name)
	ctx.Success(newLang.Commands.Language.LanguageChanged.Str(lang.DisplayName))
}

func listValidLanguages() []any {
	validLanguages := make([]any, len(i18n.LanguagesName))

	for i, lang := range i18n.LanguagesName {
		langStr := string(lang)
		validLanguages[i] = strings.ToLower(langStr)
	}
	return validLanguages
}
