package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gliderlabs/ssh"
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
	fmt.Println("Player 1's hand:", p1.left, p1.right)
	fmt.Println("Player 2's hand:", p2.left, p2.right)
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

func main() {
	p1 := &Player{"Player 1", 1, 1}
	p2 := &Player{"Player 2", 1, 1}

	fmt.Println("Da Chopsticks Game Starts:")

	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "Hello world\n")
	})

	go func() {
		err := ssh.ListenAndServe(":2222", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	DisplayState(p1, p2)

	curr := p1
	other := p2

	for !(other.DetectLose() || curr.DetectLose()) {

		fmt.Print(curr.name + "'s turn\nFrom which hand? (left, right): ")
		scanner.Scan()
		fromLeft := GetLeftRight(strings.TrimSpace(scanner.Text()))
		fmt.Print("To which hand? (left, right): ")
		scanner.Scan()
		toLeft := GetLeftRight(strings.TrimSpace(scanner.Text()))

		if fromLeft {
			other.AddToHand(curr.left, toLeft)
		} else {
			other.AddToHand(curr.right, toLeft)
		}
		DisplayState(p1, p2)
		curr, other = other, curr
	}
}
