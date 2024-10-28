package main

func (app *application) promptConfirm(message string, defaultYes bool) bool {
	options := "y/N"
	if defaultYes {
		options = "Y/n"
	}
	app.shell.Printf("%s [%s] ", message, options)

	ans := app.shell.ReadLine()
	switch ans {
	case "":
		return defaultYes
	case "y", "yes":
		return true
	default:
		return false
	}
}
