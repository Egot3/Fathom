package config

var pathToQuizzes string = "/home/ETS/"

func TurnToAbs(name string) string {
	return pathToQuizzes + name + ".md"
}
