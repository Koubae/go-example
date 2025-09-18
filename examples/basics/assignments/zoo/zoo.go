package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type AnimalT string
type ActionT string

type Animal struct {
	food       string
	locomotion string
	noise      string
}

func (a Animal) Eat() {
	fmt.Println(a.food)
}
func (a Animal) Move() {
	fmt.Println(a.locomotion)
}
func (a Animal) Speak() {
	fmt.Println(a.noise)
}

var animals = [...]string{"cow", "bird", "snake"}
var actions = [...]string{"eat", "move", "speak"}

func main() {
	commandGetter := getCommand()

	cow := Animal{
		food:       "grass",
		locomotion: "walk",
		noise:      "moo",
	}
	bird := Animal{
		food:       "worms",
		locomotion: "fly",
		noise:      "peep",
	}
	snake := Animal{
		food:       "mice",
		locomotion: "slither",
		noise:      "hiss",
	}

	for {
		animalT, actionT, err := commandGetter()
		if err != nil {
			fmt.Println(err)
			continue
		}

		var animal Animal
		switch animalT {
		case "cow":
			animal = cow
		case "bird":
			animal = bird
		case "snake":
			animal = snake
		}

		switch actionT {
		case "eat":
			animal.Eat()
		case "move":
			animal.Move()
		case "speak":
			animal.Speak()
		}

	}
}

func getCommand() func() (AnimalT, ActionT, error) {
	reader := bufio.NewReader(os.Stdin)

	return func() (AnimalT, ActionT, error) {
		fmt.Print("> ")

		action, err := reader.ReadString('\n')
		if err != nil {
			return "", "", errors.New(fmt.Sprintf("ERROR: Error while reading next action, error: %v", err))
		}

		action = strings.TrimSpace(action)
		return extractAndValidateCommand(action)

	}

}

func extractAndValidateCommand(command string) (AnimalT, ActionT, error) {
	var animalT AnimalT
	var actionT ActionT
	actionSplit := strings.Fields(command)

	if len(actionSplit) != 2 {
		return "", "", errors.New(
			fmt.Sprintf(
				"ERROR: Invalid action, expected 2 argument got %d",
				len(actionSplit),
			),
		)
	}

	animalT, actionT = AnimalT(actionSplit[0]), ActionT(actionSplit[1])
	if err := validateAnimal(animalT); err != nil {
		return "", "", err
	}
	if err := validateAction(actionT); err != nil {
		return "", "", err
	}

	return animalT, actionT, nil
}

func validateAction(action ActionT) error {
	for _, v := range actions {
		if v == string(action) {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("ERROR: Invalid action '%s', valid are %v", action, actions[:]))
}

func validateAnimal(animal AnimalT) error {
	for _, v := range animals {
		if v == string(animal) {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("ERROR: Invalid animal '%s', valid are %v", animal, animals[:]))
}
