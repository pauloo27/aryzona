package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/lmittmann/tint"
	"github.com/pauloo27/aryzona/internal/command"
	"github.com/pauloo27/aryzona/internal/command/parameters"
	"github.com/pauloo27/aryzona/internal/data/services"
	"github.com/pauloo27/aryzona/internal/discord/model"
	"github.com/pauloo27/aryzona/internal/i18n"
)

var LanguageCommand = command.Command{
	Name: "language", Aliases: []string{"lang", "locale"},
	Parameters: []*command.Parameter{
		{
			Name: "language", Type: parameters.ParameterLowerCasedString,
			ValidValuesFunc: listValidLanguages,
			Required:        false,
		},
	},
	Handler: func(ctx *command.Context) command.Result {
		if len(ctx.Args) == 0 {
			return listLanguages(ctx)
		}

		return selectLanguage(ctx)
	},
}

func listLanguages(ctx *command.Context) command.Result {
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

	return ctx.SuccessEmbed(embed)
}

func selectLanguage(ctx *command.Context) command.Result {
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
		slog.Error("Cannot set preferred lang", tint.Err(err))
		return ctx.Error(t.SomethingWentWrong.Str())
	}

	newLang, _ := i18n.GetLanguage(lang.Name)
	return ctx.Success(newLang.Commands.Language.LanguageChanged.Str(lang.DisplayName))
}

func listValidLanguages() []any {
	validLanguages := make([]any, len(i18n.LanguagesName))

	for i, lang := range i18n.LanguagesName {
		langStr := string(lang)
		validLanguages[i] = strings.ToLower(langStr)
	}
	return validLanguages
}
