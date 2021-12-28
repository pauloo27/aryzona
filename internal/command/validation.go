package command

func RunValidation(ctx *CommandContext, validation *CommandValidation) (bool, string) {
	for _, depends := range validation.DependsOn {
		ok, msg := RunValidation(ctx, depends)
		if !ok {
			return ok, msg
		}
	}
	return validation.Checker(ctx)
}
