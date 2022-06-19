package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/gliderlabs/ssh"
	terminal "github.com/quackduck/term"
)

type Player struct {
	name  string
	left  int
	right int

	// takes a prompt, returns input
	input func(string) string
}

// AddToHand adds chopsticks to this player's hand. If isLeft == true, it adds to the left hand. Otherwise, it adds it to the right hand.
func (this *Player) AddToHand(n int, isLeft bool) {
	if isLeft {
		this.left = (this.left + n) % 5
	} else {
		this.right = (this.right + n) % 5
	}
}
func (this *Player) DetectLose() bool {
	return this.left == 0 && this.right == 0
}

func DisplayState(p1 *Player, p2 *Player) {
	state := fmt.Sprintln(p1.name+"'s hand:", p1.left, p1.right) +
		fmt.Sprintln(p2.name+"'s hand:", p2.left, p2.right) +
		fmt.Sprintln("----------------------\n")

	res := string(markdown.Render("![lol](./finger.png)", 80, 6))
	fmt.Println(res)

	fmt.Println(state)
	termPrintln(state)
}

func GetLeftRight(input string) bool {
	if input == "left" {
		return true
	} else if input == "right" {
		return false
	} else {
		fmt.Println("I'm just going with right")
		return false
	}
}

var (
	scanner = bufio.NewScanner(os.Stdin)
	term    *terminal.Terminal
)

func input(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func termPrintln(s string) {
	term.Write([]byte(s + "\n"))
}

func termInput(prompt string) string {
	term.SetPrompt(prompt)
	line, _ := term.ReadLine()
	return line
}

func main() {
	p1 := &Player{input("Enter your name: "), 1, 1, input}
	p2 := &Player{"p2", 1, 1, termInput} // name gets changed later

	fmt.Println("Da Chopsticks Game Starts:")

	gameReadyChan := make(chan bool)

	ssh.Handle(func(s ssh.Session) {
		term = terminal.NewTerminal(s, "> ")
		pty, winChan, _ := s.Pty()
		w := pty.Window
		_ = term.SetSize(w.Width, w.Height)

		go func() {
			for w = range winChan {
				_ = term.SetSize(w.Width, w.Height)
			}
		}()

		term.SetPrompt("Enter name: ")
		line, err := term.ReadLine()
		if err != nil {
			fmt.Println(err)
		}
		p2.name = line
		gameReadyChan <- true
		for {
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "2155"
	}

	go func() {
		err := ssh.ListenAndServe(fmt.Sprintf(":%s", port), nil, ssh.HostKeyFile(os.Getenv("HOME")+"/.ssh/id_rsa"))
		if err != nil {
			fmt.Println(err)
		}
	}()

	<-gameReadyChan

	DisplayState(p1, p2)

	curr := p1
	other := p2

	for !(other.DetectLose() || curr.DetectLose()) {
		fromLeft := GetLeftRight(curr.input(curr.name + "'s turn\nFrom which hand? (left, right): "))
		toLeft := GetLeftRight(curr.input("To which hand? (left, right): "))
		if fromLeft {
			other.AddToHand(curr.left, toLeft)
		} else {
			other.AddToHand(curr.right, toLeft)
		}
		DisplayState(p1, p2)
		curr, other = other, curr
	}
}
