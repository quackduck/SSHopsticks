package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gliderlabs/ssh"
	terminal "github.com/quackduck/term"
)

type Player struct {
	name  string
	left  int
	right int
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
	fmt.Println(p1.name+"'s hand:", p1.left, p1.right)
	fmt.Println(p2.name+"'s hand:", p2.left, p2.right)
	fmt.Println("----------------------\n")
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
)

func input() string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func main() {
	fmt.Print("Enter your name: ")
	p1 := &Player{input(), 1, 1}
	p2 := &Player{"Player 2", 1, 1}

	fmt.Println("Da Chopsticks Game Starts:")

	ssh.Handle(func(s ssh.Session) {
		term := terminal.NewTerminal(s, "> ")
		pty, winChan, _ := s.Pty()
		w := pty.Window
		_ = term.SetSize(w.Width, w.Height)
		go func() {
			for w = range winChan {
				_ = term.SetSize(w.Width, w.Height)
			}
		}()
		term.Write([]byte("hello world! enter smth\n"))
		line, err := term.ReadLine()
		if err != nil {
			fmt.Println(err)
		}
		term.Write([]byte("You said " + line))
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

	DisplayState(p1, p2)

	curr := p1
	other := p2

	for !(other.DetectLose() || curr.DetectLose()) {

		fmt.Print(curr.name + "'s turn\nFrom which hand? (left, right): ")
		fromLeft := GetLeftRight(input())
		fmt.Print("To which hand? (left, right): ")
		toLeft := GetLeftRight(input())

		if fromLeft {
			other.AddToHand(curr.left, toLeft)
		} else {
			other.AddToHand(curr.right, toLeft)
		}
		DisplayState(p1, p2)
		curr, other = other, curr
	}
}
