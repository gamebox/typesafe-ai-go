package main

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gamebox/typesafe-ai-go"
)

//go:embed fireworks.vt
var winAnimation string

//go:embed barney.vt
var loseAnimation string

type Mark int

const (
	MarkNone = iota
	MarkX
	MarkO
)

func (m Mark) Char() rune {
	switch m {
	case MarkNone:
		return ' '
	case MarkX:
		return 'X'
	case MarkO:
		return 'O'
	}
	return ' '
}

func dotBlank(r rune) rune {
	if r == ' ' {
		return '.'
	}
	return r
}

type Board [3][3]Mark

func (b Board) String() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%c%c%c\n", dotBlank(b[0][0].Char()), dotBlank(b[0][1].Char()), dotBlank(b[0][2].Char()))
	fmt.Fprintf(&sb, "%c%c%c\n", dotBlank(b[1][0].Char()), dotBlank(b[1][1].Char()), dotBlank(b[1][2].Char()))
	fmt.Fprintf(&sb, "%c%c%c\n", dotBlank(b[2][0].Char()), dotBlank(b[2][1].Char()), dotBlank(b[2][2].Char()))
	return sb.String()
}
func (b Board) ValidMoves() [][2]int {
	moves := make([][2]int, 0, 9)
	for y := 0; y < 3; y = y + 1 {
		for x := 0; x < 3; x = x + 1 {
			if b[y][x] == MarkNone {
				moves = append(moves, [2]int{y, x})
			}
		}
	}
	return moves
}

type Game struct {
	b Board
	c *typesafe.Client
}

func NewGame(key string) Game {
	return Game{
		c: typesafe.NewClient(key),
	}
}

var playerMarks = [2]Mark{MarkX, MarkO}

func (g Game) HaveWinner() bool {
	for _, mark := range playerMarks {
		for y := 0; y < 3; y = y + 1 {
			if g.b[y][0] == mark && g.b[y][1] == mark && g.b[y][2] == mark {
				return true
			}
		}
		for x := 0; x < 3; x = x + 1 {
			if g.b[0][x] == mark && g.b[1][x] == mark && g.b[2][x] == mark {
				return true
			}
		}
		if g.b[0][0] == mark && g.b[1][1] == mark && g.b[2][2] == mark {
			return true
		}
		if g.b[0][2] == mark && g.b[1][1] == mark && g.b[2][0] == mark {
			return true
		}
	}
	return false
}

func (g Game) HaveDraw() bool {
	for y := 0; y < 3; y = y + 1 {
		for x := 0; x < 3; x = x + 1 {
			if g.b[y][x] == MarkNone {
				return false
			}
		}
	}
	return true
}

func (g *Game) Reset() {
	g.b = Board{}
}

func main() {
	key, ok := os.LookupEnv("TYPESAFE_API_KEY")
	if !ok {
		fmt.Fprintf(os.Stderr, "Must have TYPESAFE_API_KEY in environment")
	}

	restore := setupTerm()
	defer restore()

	closeLogs := setupLogs()
	defer closeLogs()

	g := NewGame(key)

	for {
		printBoard(g.b)
		err := readPlayerMove(&g)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			break
		}
		cursorHidden()
		if g.HaveWinner() {
			if !showWinScreen(&g) {
				break
			}
			continue
		}
		if g.HaveDraw() {
			if !showDrawScreen(&g) {
				break
			}
			continue
		}
		getJevMove(&g)
		if g.HaveWinner() {
			if !showLoseScreen(&g) {
				break
			}
			continue
		}
		if g.HaveDraw() {
			if !showDrawScreen(&g) {
				break
			}
			continue
		}
	}
}

func getJevMove(g *Game) error {
	instructions := strings.Builder{}
	fmt.Fprintf(&instructions, "This is a game of tic-tac-toe.  You are O. Your opponent is X. A '.' indicates an open spot on the board.\n")
	fmt.Fprintf(&instructions, "%s", g.b.String())
	fmt.Fprint(&instructions, "The rules of tic-tac-toe:")
	fmt.Fprint(&instructions, "1. There are two players that take turns placing marks.")
	fmt.Fprint(&instructions, "2. One player places X marks, the other places O marks.")
	fmt.Fprint(&instructions, "3. The object of the game is to get three of their marks in a row.  Vertical, Horizontally, or Diagonally")
	slog.Info("Making jev request", "instructions", instructions.String())
	r := typesafe.NewRequest(instructions.String())
	next_q := typesafe.Choice("What is your best move (line, column)?")
	moves := g.b.ValidMoves()
	for i, m := range moves {
		next_q.AddCriteria(fmt.Sprintf("Choice %d", i), fmt.Sprintf("(%d, %d)", m[0], m[1]))
	}
	r.AddQuestion("next", next_q)
	rules_q := typesafe.Noul("Do you know the rules of tic-tac-toe?")
	r.AddQuestion("rules", rules_q)
	resp, err := g.c.Ask(r)
	if err != nil {
		return err
	}
	rules_answer, ok := resp.DecodeNoul("rules")
	if !ok {
		return fmt.Errorf("Could not decode rules answer")
	}
	slog.Info("Got rules answer", "answer", rules_answer.Noul)
	next_answer, ok := resp.DecodeChoice("next")
	if !ok {
		return fmt.Errorf("Could not decode answer")
	}
	var choice int
	_, err = fmt.Fscanf(strings.NewReader(next_answer.Choice), "Choice %d", &choice)
	if err != nil {
		return err
	}
	move := moves[choice]
	g.b[move[0]][move[1]] = MarkO
	return nil
}

func playAgainPrompt() bool {
	fmt.Println("Play again? y for yes, n for no")
	for {
		_, err := os.Stdin.Read(b[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read input: %e", err)
			os.Exit(1)
		}

		m := string(b[:])
		b = [1]byte{}
		switch m {
		case "y":
			return true
		case "n":
			return false
		case "q":
			return false
		default:
			// Ignore other characters
		}
	}
}

func showWinScreen(g *Game) bool {
	cursorHidden()
	playVtAnim(winAnimation, 30*time.Millisecond)
	eraseScreen()
	moveCursorHome()
	fmt.Println("You win! 🥇")
	moveCursorBONL()
	if !playAgainPrompt() {
		return false
	}
	g.Reset()
	return true
}
func showLoseScreen(g *Game) bool {
	cursorHidden()
	playVtAnim(loseAnimation, 5*time.Millisecond)
	eraseScreen()
	moveCursorHome()
	fmt.Println("You lose... 😭")
	moveCursorBONL()
	if !playAgainPrompt() {
		return false
	}
	g.Reset()
	return true
}

func showDrawScreen(g *Game) bool {
	cursorHidden()
	eraseScreen()
	moveCursorHome()
	fmt.Println("It's a draw! 🤝")
	moveCursorBONL()
	if !playAgainPrompt() {
		return false
	}
	g.Reset()
	return true
}

var b [1]byte

func readPlayerMove(g *Game) error {
	posX, posY := 0, 0
	cursorVisible()
	moveCursorTo((posY+1)*2, (posX+1)*2)
loop:
	for {
		_, err := os.Stdin.Read(b[:])
		if err != nil {
			slog.Error("Could not read input", "error", err)
			os.Exit(1)
		}

		m := string(b[:])
		b = [1]byte{}

		switch m {
		case "w":
			posY = max(posY-1, 0)
		case "a":
			posX = max(posX-1, 0)
		case "s":
			posY = min(posY+1, 2)
		case "d":
			posX = min(posX+1, 2)
		case "q":
			return fmt.Errorf("Good bye!")
		case "\r":
			if m := g.b[posY][posX]; m == MarkNone {
				break loop
			}
			continue
		default:
			continue
		}
		moveCursorTo((posY+1)*2, (posX+1)*2)
	}

	g.b[posY][posX] = MarkX
	printBoard(g.b)
	return nil
}

func printBoard(b Board) {
	eraseScreen()
	moveCursorHome()

	fmt.Print("┏━┳━┳━┓")
	moveCursorBONL()
	fmt.Printf("┃%c┃%c┃%c┃", b[0][0].Char(), b[0][1].Char(), b[0][2].Char())
	moveCursorBONL()
	fmt.Print("┣━╋━╋━┫")
	moveCursorBONL()
	fmt.Printf("┃%c┃%c┃%c┃", b[1][0].Char(), b[1][1].Char(), b[1][2].Char())
	moveCursorBONL()
	fmt.Print("┣━╋━╋━┫")
	moveCursorBONL()
	fmt.Printf("┃%c┃%c┃%c┃", b[2][0].Char(), b[2][1].Char(), b[2][2].Char())
	moveCursorBONL()
	fmt.Print("┗━┻━┻━┛")
	moveCursorBONL()
	fmt.Printf("W/S/A/D Move Up/Down/Left/Right")
	moveCursorBONL()
	fmt.Print("Q Quit")
}

func playVtAnim(data string, sleepTime time.Duration) {
	eraseScreen()
	moveCursorHome()
	for line := range strings.Lines(data) {
		fmt.Print(line)
		time.Sleep(sleepTime)
	}
}
