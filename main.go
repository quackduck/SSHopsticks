package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
	markdown "github.com/quackduck/go-term-markdown"
	terminal "github.com/quackduck/term"
)

var (
	scanner = bufio.NewScanner(os.Stdin)

	term *terminal.Terminal
	w    ssh.Window

	downFinger = string(markdown.Render("![lol](https://cloud-lvtf5ds2i-hack-club-bot.vercel.app/0image.png)", 10, 0))
	finger     = string(markdown.Render("![lol](https://cloud-6zj0ryec6-hack-club-bot.vercel.app/0finger.png)", 10, 0))

	gameReadyChan = make(chan bool)

	p1 = &Player{input("Enter your name: "), 1, 1, input, stdPrintln}
	p2 = &Player{"p2", 1, 1, termInput, termPrintln}
)

type Player struct {
	name  string
	left  int
	right int

	// takes a prompt, returns input
	input  func(string) string
	output func(string)
}

func (this *Player) AddToHand(n int, isLeft bool) {
	if isLeft {
		this.left = (this.left + n) % 5
	} else {
		this.right = (this.right + n) % 5
	}
}
func (this *Player) DetectLoss() bool {
	return this.left == 0 && this.right == 0
}

func DisplayState(curr *Player, other *Player) {
	curr.output(stateAs(curr, other))
	other.output(stateAs(other, curr))
}

func stateAs(curr *Player, other *Player) string {
	return other.name + "'s hand" +
		showFingers(other.left, other.right, false) +
		"\n" + curr.name + "'s hand" +
		showFingers(curr.left, curr.right, true)
}

func showFingers(num int, num2 int, up bool) string {
	fSplit := strings.Split(finger, "\n")
	if !up {
		fSplit = strings.Split(downFinger, "\n")
	}
	ret := ""

	for i := 0; i < len(fSplit); i++ {
		if i == len(fSplit)-10 { // magic number, do not change.
			break
		}
		for j := 0; j < num; j++ {
			ret += fSplit[i]
		}
		ret += "       "
		for j := 0; j < num2; j++ {
			ret += fSplit[i]
		}
		ret += "\n"
	}

	return ret
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

func input(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func stdPrintln(s string) {
	fmt.Println(s)
}

func termPrintln(s string) {
	term.Write([]byte(s + "\n"))
}

func termInput(prompt string) string {
	term.SetPrompt(prompt)
	line, _ := term.ReadLine()
	return line
}

func SshHandler(s ssh.Session) {
	term = terminal.NewTerminal(s, "> ")
	pty, winChan, _ := s.Pty()
	w = pty.Window
	_ = term.SetSize(w.Width, w.Height)

	go func() {
		for w = range winChan {
			_ = term.SetSize(w.Width, w.Height)
		}
	}()

	p2.name = termInput("Enter your name: ")
	gameReadyChan <- true
	for {
	}
}

func DoTurn(curr *Player, other *Player) error {
	whatKind := curr.input("From which hand (or split)? (left, right, split): ")

	if whatKind == "split" {
		fromLeft := GetLeftRight(curr.input("From which hand? (left, right): "))
		lols, err := strconv.Atoi(curr.input("How many? "))

		if err != nil {
			return err
		}

		if lols < 0 {
			return fmt.Errorf("naughtry boy")
		}

		if fromLeft {
			if lols+curr.right >= 5 {
				return fmt.Errorf("naughtry boy")
			}

			curr.right += lols
			curr.left -= lols
		} else {
			if lols+curr.left >= 5 {
				return fmt.Errorf("naughtry boy")
			}

			curr.left += lols
			curr.right -= lols
		}

		return nil
	}

	fromLeft := GetLeftRight(whatKind)
	toLeft := GetLeftRight(curr.input("To which hand? (left, right): "))

	if fromLeft {
		if other.left != 0 {
			other.AddToHand(curr.left, toLeft)
		} else {
			curr.output("You can't do that")
			return fmt.Errorf("no")
		}
	} else {
		if other.right != 0 {
			other.AddToHand(curr.right, toLeft)
		} else {
			curr.output("STOP!!!!!!")
			return fmt.Errorf("no")
		}
	}

	return nil
}

func GameLoop() {
	curr := p1
	other := p2

	for !(other.DetectLoss() || curr.DetectLoss()) {
		curr.output("your turn")
		other.output(curr.name + "'s turn...")

		err := DoTurn(curr, other)
		if err != nil {
			fmt.Println(err)
			continue
		}

		DisplayState(curr, other)
		curr, other = other, curr
	}
}

func ServeSsh(port string) {
	hostKey := ssh.HostKeyFile(os.Getenv("HOME") + "/.ssh/id_rsa")

	err := ssh.ListenAndServe(fmt.Sprintf(":%s", port), nil, hostKey)

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("Da Chopsticks Game Starts:")
	ssh.Handle(SshHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "2155"
	}

	fmt.Println("Please get the other player to ssh to port: " + port)

	go ServeSsh(port)
	<-gameReadyChan

	DisplayState(p1, p2)
	GameLoop()
}
