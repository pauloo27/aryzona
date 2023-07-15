package command

func RunValidation(ctx *Context, validation *Validation) (bool, string) {
	for _, depends := range validation.DependsOn {
		ok, msg := RunValidation(ctx, depends)
		if !ok {
			return ok, msg
		}
	}
	return validation.Checker(ctx)
}
