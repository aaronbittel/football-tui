package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"tui/internal/component"
	utils "tui/internal/term-utils"

	"golang.org/x/term"
)

type Task struct {
	name    string
	steps   int
	running bool
}

func NewTask(name string, steps int) *Task {
	return &Task{
		name:  name,
		steps: steps,
	}
}

func main() {
	fd, oldState, err := utils.Start()
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	defer utils.TearDown()

	list := component.NewList("Einer", "Zehner", "Hunderter").At(10, 10)
	component.Print(list)

	var task *Task

	controlCh := make(chan string)
	defer close(controlCh)

	reader := bufio.NewReader(os.Stdin)
	running := true

	for running {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', utils.CtrlC:
			running = false
		case 'j':
			list.Next()
			task = stopTask(controlCh, task)
			utils.Debug("task", task)
		case 'k':
			list.Prev()
			task = stopTask(controlCh, task)
		case utils.Enter:
			var steps int
			switch list.Selected {
			case 0:
				steps = 1
			case 1:
				steps = 10
			case 2:
				steps = 100
			}
			task = NewTask(list.SelectedValue(), steps)
			utils.Debug(fmt.Sprintf("%s task is running", task.name))
			go Algo(controlCh, task)
		case 'x':
			task = stopTask(controlCh, task)
		case utils.Space:
			if task == nil {
				break
			}
			if task.running {
				controlCh <- "PAUSED"
				utils.Debug(fmt.Sprintf("%s task is paused", task.name))
				task.running = false
			} else {
				controlCh <- "RUNNING"
				utils.Debug(fmt.Sprintf("%s task is running", task.name))
				task.running = true
			}
		case 'n':
			if task == nil {
				break
			}
			controlCh <- "PAUSED"
			controlCh <- "NEXT"
			task.running = false
		case 'p':
			if task == nil {
				break
			}
			controlCh <- "PAUSED"
			controlCh <- "PREV"
			task.running = false
		}
	}
}

func Algo(controlCh <-chan string, task *Task) {
	i := 0
	state := "RUNNING"

outer:
	for {
		select {
		case state = <-controlCh:
			switch state {
			case "STOP":
				break outer
			case "NEXT":
				utils.ClearLine(10, 30)
				i += task.steps
				fmt.Print(i)
				state = "PAUSED"
			case "PREV":
				i -= task.steps
				utils.ClearLine(10, 30)
				fmt.Print(i)
				state = "PAUSED"
			}
		default:
			if state == "PAUSED" {
				break
			}
			utils.MoveCursor(10, 30)
			i += task.steps
			fmt.Print(i)
			time.Sleep(time.Millisecond * 500)

		}
	}

	utils.Debug("FINISHED")
}

func stopTask(controlCh chan<- string, task *Task) *Task {
	if task == nil {
		return task
	}
	controlCh <- "STOP"
	utils.Debug(fmt.Sprintf("%s task is stopped", task.name))
	utils.ClearLine(10, 30)
	return nil
}
