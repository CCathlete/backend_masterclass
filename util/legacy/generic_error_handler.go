package util

// func HandleError(
// 	f func(...any) ([]any, error),
// 	callerName string,
// 	args ...any,
// ) any {

// 	vals, err := f(args...)

// 	// Error handling logic.
// 	switch callerName {
// 	case "main":
// 		if err != nil {
// 			panic(err)
// 		}
// 	case "createAccount":
// 		ctx := args[0].(*gin.Context)
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

// 	}

// 	return vals
// }
