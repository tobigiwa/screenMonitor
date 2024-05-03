package frontend

import "math/rand"

var (
	art1 string = `
┓ •┏┓     ┳┳┓
┃ ┓┗┓┏┏┓┏┓┃┃┃┏┓┏┓
┗┛┗┗┛┗┛ ┗ ┛ ┗┗┛┛┗`

	art2 string = `
╦ ┬ ╔═╗┌─┐┬─┐┌─┐╔╦╗┌─┐┌┐┌
║ │ ╚═╗│  ├┬┘├┤ ║║║│ ││││
╩═╝┴╚═╝└─┘┴└─└─┘╩ ╩└─┘┘└┘`

	art3 string = `
+-+-+-+-+-+-+-+-+-+
|L|i|S|c|r|e|M|o|n|
+-+-+-+-+-+-+-+-+-+`
)

var arrayOfArt = []string{art1, art2, art3}

func getRandomArt() string {
	randonIndex := rand.Intn(len(arrayOfArt))
	return arrayOfArt[randonIndex]
}
