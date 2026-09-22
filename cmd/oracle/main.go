package main

import (
	"fmt"
	"os"
	"strings"

	rl "github.com/chzyer/readline"
	"github.com/gamebox/typesafe-ai-go"
)

func main() {
	key, ok := os.LookupEnv("TYPESAFE_API_KEY")
	if !ok {
		fmt.Printf("No key\n")
		os.Exit(1)
	}
	client := typesafe.NewClient(key)

	i, err := rl.NewEx(&rl.Config{
		Prompt:    ">",
		EOFPrompt: "quit",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create readline interface: %e\n", err)
	}
	for {
		collectRequestDetails(client, i)
		if !collectWantsAnother(i) {
			break
		}
	}
}
func collectWantsAnother(i *rl.Instance) bool {
	for {
		i.SetPrompt("Do you want to make another request? (y/n) > ")
		line, err := i.Readline()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read input: %e\n", err)
		}
		choice := strings.ToLower(strings.TrimSpace(line))
		switch choice {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Printf("This is an invalid option, please enter y or n\n")
		}
	}
}

func collectRequestDetails(client *typesafe.Client, i *rl.Instance) {
	questions := map[string]typesafe.Question{}
	num_qs := 0

	state := getState(i)
	for {
		q, ok := getQuestion(i, num_qs > 0)
		if !ok {
			break
		}
		questions[fmt.Sprintf("Q%d", num_qs)] = typesafe.Noul(q)
		num_qs++
	}

	response, err := client.Ask(typesafe.Request{
		State:     state,
		Questions: questions,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not make request: %e\n", err)
		os.Exit(1)
	}

	for key := range questions {
		q := questions[key]
		a, ok := response.DecodeNoul(key)
		if !ok {
			fmt.Printf("Could not decode answer for %s: %#v\n", key, response)
			os.Exit(1)
		}
		PrintNoulAnswer(q.(typesafe.NoulQuestion), a)
	}
}

func getState(i *rl.Instance) string {
	for {
		fmt.Print(strings.TrimLeft(`
What do you want to ask about?
1) Some text I enter here
2) A file
	`, " \t\n"))
		i.SetPrompt("Your choice > ")
		c, err := i.Readline()
		if err != nil {
			os.Exit(1)
		}
		choice := strings.TrimSpace(c)
		state := ""
		switch choice {
		case "1":
			i.SetPrompt("Enter your input > ")
			state, err = i.Readline()
			if err != nil {
				os.Exit(1)
			}
			return state
		case "2":
			i.SetPrompt("Enter the path to the file > ")
			path, err := i.Readline()
			if err != nil {
				os.Exit(1)
			}
			bytes, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not open file: %e\n", err)
				os.Exit(1)
			}
			return string(bytes)
		default:
			fmt.Fprintf(os.Stderr, "Invalid choice: [%s]\n", choice)
		}
	}
}

func getQuestion(i *rl.Instance, first bool) (q string, ok bool) {
	if first {
	loop:
		for {
			i.SetPrompt("Do you have another question?(Y/N) > ")
			raw_line, err := i.Readline()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not read input\n")
				os.Exit(1)
			}
			choice := strings.ToLower(strings.TrimSpace(raw_line))
			switch choice {
			case "n":
				return "", false
			case "y":
				break loop
			default:
				fmt.Printf("Invalid choice, please enter Y or N\n")
			}
		}
	}

	i.SetPrompt("Enter your question > ")
	q_line, err := i.Readline()
	if err != nil {
		os.Exit(1)
	}
	q = strings.TrimSpace(q_line)

	return q, true
}

func PrintNoulAnswer(q typesafe.NoulQuestion, a typesafe.NoulAnswer) {
	fmt.Printf("Q: %s\nA: %t (%d%% confidence)\n", q.Instructions.(string), a.Noul > 0.50, int(a.Noul*100))
}
