# Terminal Algorithm Visualizer

This is a Terminal-based Algorithm Visualizer (TUI) built with Go. It allows
you to select from various algorithms and watch their execution unfold
step-by-step, right in your terminal. This project is developed entirely from
scratch without the use of any TUI frameworks, utilizing ANSI escape codes for
terminal control.

## Table of Contents
- [Features](#features)
- [Available Algorithms](#available-algorithms)
- [Installation](#installation)
- [App Showcase](#app-showcase)
- [Shortcomings](#shortcomings)

## Features

- Simple and easy-to-use terminal interface
- Select algorithms easily from a list
- Real-time visualization of algorithm execution
- Control execution flow with Start, Stop, and Pause options
- Step forward and backward through the algorithm visualization
- Each visualization step includes a description of what is happening
- Adjust execution speed (speed up or slow down) for better understanding

## Available Algorithms

### Sorting Algorithms
- Bubblesort
- Insertionsort
- Selectionsort
- Quicksort
- Mergesort
- Heapsort

### Search Algorithms
- In progress


## Installation

Follow these step-by-step instructions to set up the **Terminal Algorithm
Visualizer** project locally. No prerequisites are required to run the compiled
binaries. However, if you wish to build the binary yourself, you will need to
have **Go** installed. If you haven't installed it yet, please follow the
[official Go installation guide](https://go.dev/doc/install).

```bash
git clone https://github.com/aaronbittel/terminal-algorithm-visualizer
cd terminal-algorithm-visualizer
```

You can run the compiled executable directly from the /bin folder. Use the command that corresponds to your operating system:

For Linux
```bash
./bin/algo
```

For Windows
```bash
./bin/algo.exe
```
If you would like to build the project yourself, you can do so by using the
following command. This step requires Go to be installed:

For Linux
```bash
make build
```

For Windows
```bash
make build-win
```

To see a list of available commands and options for your application, run:
```bash
make help
```

## App Showcase

![tav](https://github.com/user-attachments/assets/f0fe458d-7959-4bbb-b5a3-8c1e0fe7c7dc)

## Shortcomings

- Currently, the positions of widgets are hard-coded, which requires a minimum
terminal size to display everything correctly.
- Heapsort not quite finished yet
