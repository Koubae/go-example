package main

import (
	"fmt"
)

type Animal interface {
	Eat()
	Move()
	Speak()
}

type Cow struct{}

func (v Cow) Eat() {
	fmt.Println("grass")
}

func (v Cow) Move() {
	fmt.Println("walk")
}

func (v Cow) Speak() {
	fmt.Println("moo")
}

type Bird struct{}

func (v Bird) Eat() {
	fmt.Println("worms")
}

func (v Bird) Move() {
	fmt.Println("fly")
}

func (v Bird) Speak() {
	fmt.Println("peep")
}

type Snake struct{}

func (v Snake) Eat() {
	fmt.Println("mice")
}

func (v Snake) Move() {
	fmt.Println("slither")
}

func (v Snake) Speak() {
	fmt.Println("hsss")
}

func main() {
	animals := make(map[string]Animal)
	var command, name, param string
	for true {
		fmt.Print("> ")
		fmt.Scanf("%s %s %s", &command, &name, &param)
		switch command {
		case "newanimal":
			switch param {
			case "cow":
				animals[name] = Cow{}
				fmt.Println("Created it!")
			case "bird":
				animals[name] = Bird{}
				fmt.Println("Created it!")
			case "snake":
				animals[name] = Snake{}
				fmt.Println("Created it!")
			default:
				fmt.Println("Invalid animal type")
			}
		case "query":
			if val, ok := animals[name]; ok {
				switch param {
				case "eat":
					val.Eat()
				case "move":
					val.Move()
				case "speak":
					val.Speak()
				default:
					fmt.Println("Invalid query")
				}
			} else {
				fmt.Println("No animal with this name")
			}
		default:
			fmt.Println("Invalid command type")
		}
	}
}
