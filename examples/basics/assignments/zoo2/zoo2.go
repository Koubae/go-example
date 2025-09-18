package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Animal interface {
	Eat()
	Move()
	Speak()
}

type Cow struct {
	food       string
	locomotion string
	noise      string
}

func (a Cow) Eat() {
	fmt.Println(a.food)
}

func (a Cow) Move() {
	fmt.Println(a.locomotion)
}

func (a Cow) Speak() {
	fmt.Println(a.noise)
}

type Bird struct {
	food       string
	locomotion string
	noise      string
}

func (a Bird) Eat() {
	fmt.Println(a.food)
}

func (a Bird) Move() {
	fmt.Println(a.locomotion)
}

func (a Bird) Speak() {
	fmt.Println(a.noise)
}

type Snake struct {
	food       string
	locomotion string
	noise      string
}

func (a Snake) Eat() {
	fmt.Println(a.food)
}

func (a Snake) Move() {
	fmt.Println(a.locomotion)
}

func (a Snake) Speak() {
	fmt.Println(a.noise)
}

const CommandNewAnimalCode = "newanimal"
const CommandQueryCode = "query"

func main() {
	getCommand := getCommandBuilder()

	zooDb := make(map[string]Animal)

	for {
		action, command, err := getCommand()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch action {
		case CommandNewAnimalCode:
			err = processNewAnimalCommand(zooDb, command)
		case CommandQueryCode:
			err = processQueryCommand(zooDb, command)
		default:
			err = errors.New(
				fmt.Sprintf(
					"ERROR: Invalid action '%s', supported actions are %v",
					action,
					[...]string{CommandNewAnimalCode, CommandQueryCode},
				),
			)
		}

		if err != nil {
			fmt.Println(err)
		}

	}

}

func getCommandBuilder() func() (string, []string, error) {
	reader := bufio.NewReader(os.Stdin)
	return func() (string, []string, error) {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", nil, errors.New(fmt.Sprintf("ERROR: Error while reading next action, error: %v", err))
		}
		input = strings.TrimSpace(input)

		fields := strings.Fields(input)

		action := fields[0]
		command := fields[1:]
		return action, command, nil

	}
}

func processNewAnimalCommand(zooDb map[string]Animal, command []string) error {
	if len(command) != 2 {
		return errors.New(
			fmt.Sprintf(
				"ERROR: Invalid command for new animal, expected 2 argument got %d",
				len(command),
			),
		)
	}

	name := command[0]
	animalT := command[1]

	var animal Animal
	switch animalT {
	case "cow":
		animal = &Cow{
			food:       "grass",
			locomotion: "walk",
			noise:      "moo",
		}
	case "bird":
		animal = &Bird{
			food:       "worms",
			locomotion: "fly",
			noise:      "peep",
		}
	case "snake":
		animal = &Snake{
			food:       "mice",
			locomotion: "slither",
			noise:      "hiss",
		}
	default:
		return errors.New(fmt.Sprintf("ERROR: Unknown animal type '%s'", animalT))
	}

	zooDb[name] = animal

	fmt.Println("Created it!")
	return nil
}

func processQueryCommand(zooDb map[string]Animal, command []string) error {
	if len(command) != 2 {
		return errors.New(
			fmt.Sprintf(
				"ERROR: Invalid command for query command, expected 2 argument got %d",
				len(command),
			),
		)
	}

	name := command[0]
	action := command[1]

	animal, ok := zooDb[name]
	if !ok {
		return errors.New(fmt.Sprintf("ERROR: Animal '%s' not found", name))
	} else if animal == nil {
		return errors.New(fmt.Sprintf("ERROR: Animal '%s' is nil", name))
	}

	switch action {
	case "eat":
		animal.Eat()
	case "move":
		animal.Move()
	case "speak":
		animal.Speak()
	}

	return nil
}
