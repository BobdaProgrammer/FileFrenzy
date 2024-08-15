package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
)

var (
	startPath string
	ChosenOne string
	grid      [][]string
	Player    PlayerCh
	selected  string
	isCorrect bool
)

type PlayerCh struct {
	posX, posY int
}

func indexOf(slice []string, str string) int {
	for i, s := range slice {
		if s == str {
			return i
		}
	}
	return -1 // Return -1 if the string is not found
}

func genStartingPath() {
	curr, err := os.Getwd()
	if err == nil {
		Folders := strings.Count(curr, "\\") - 1
		Place := int(rand.Float64() * float64(Folders))
		for i := 0; i < Place; i++ {
			curr = filepath.Dir(curr)
		}
		fmt.Println(curr)
		startPath = curr
	} else {
		log.Fatal("Couldn't get your directory :(")
	}
}

// Check if a directory is empty
func isDirEmpty(name string) bool {
	entries, err := os.ReadDir(name)
	if err != nil {
		return true
	}
	return len(entries) == 0
}

// GetRandFile tries to select a random non-empty directory or file
func GetRandFile(files []os.DirEntry, basePath string) string {
	for len(files) > 0 {
		pos := rand.Intn(len(files))
		selected := files[pos]
		// If it's a directory and not empty, return it
		if selected.IsDir() && !isDirEmpty(filepath.Join(basePath, selected.Name())) {
			return selected.Name()
		}
		// Remove the selected item and continue
		files = append(files[:pos], files[pos+1:]...)
	}
	// Return an empty string if no non-empty directory is found
	return ""
}

// GetEndFile selects a random file from the list
func GetEndFile(files []os.DirEntry) string {
	for len(files) > 0 {
		pos := rand.Intn(len(files))
		selected := files[pos]
		// If it's a file, return it
		if !selected.IsDir() {
			return selected.Name()
		}
		// Remove the selected item and continue
		files = append(files[:pos], files[pos+1:]...)
	}
	// Return an empty string if no file is found
	return ""
}

func Render(s tcell.Screen, files []os.DirEntry) {
	screenWidth, screenHeight := s.Size()
	text := "Target: " + selected
	isCorrect = false
	correctCol := tcell.StyleDefault.Foreground(tcell.ColorLightGreen.TrueColor())
	incorrectCol := tcell.StyleDefault.Foreground(tcell.ColorRed.TrueColor())
	WhichCol := incorrectCol
	end := ""
	for i, ch := range text {
		s.SetContent(i, 0, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorRed.TrueColor()))
	}
	for x := 0; x < screenWidth; x++ {
		s.SetContent(x, 1, '#', nil, tcell.StyleDefault.Foreground(tcell.ColorLightGreen.TrueColor()))
	}
	for y := range grid {
		s.SetContent(0, y+2, '#', nil, tcell.StyleDefault.Foreground(tcell.ColorLightGreen.TrueColor()))
		for x, name := range grid[y] {
			if name != "" {
				cont := '📁'
				if (Player.posX == x || Player.posX == x+1) && Player.posY == y {
					name := grid[Player.posY][Player.posX]
					if Player.posX == x+1 {
						name = grid[Player.posY][Player.posX-1]
					}
					end = "File: " + name
					if name == selected {
						WhichCol = correctCol
						isCorrect = true
					}
					cont = ' '
				}
				s.SetContent(x+1, y+2, cont, nil, tcell.StyleDefault)
			} else {
				s.SetContent(x+1, y+2, ' ', nil, tcell.StyleDefault)
			}
			if x == Player.posX && y == Player.posY {
				s.SetContent(Player.posX+1, Player.posY+2, ' ', nil, tcell.StyleDefault.Background(tcell.ColorRed.TrueColor()))
			}
		}
		s.SetContent(screenWidth-1, y+2, '#', nil, tcell.StyleDefault.Foreground(tcell.ColorLightGreen.TrueColor()))
	}
	for x := 0; x < screenWidth; x++ {
		s.SetContent(x, screenHeight-2, '#', nil, tcell.StyleDefault.Foreground(tcell.ColorLightGreen.TrueColor()))
	}
	if end == "" {
		for x := 0; x < screenWidth; x++ {
			s.SetContent(x, screenHeight-1, ' ', nil, tcell.StyleDefault)
		}
	} else {
		for i, ch := range end {
			s.SetContent(i, screenHeight-1, ch, nil, WhichCol)
		}
	}
}

func genGrid(files []os.DirEntry, s tcell.Screen) {
	screenWidth, screenHeight := s.Size()
	selected = GetRandFile(files, startPath)
	if selected == "" {
		selected = GetEndFile(files)
		if selected == "" {
			s.Clear()
			WinScreen()
			os.Exit(0)
		}
	}
	width := screenWidth - 3
	height := screenHeight - 4
	grid = make([][]string, height)
	for i := 0; i < len(grid); i++ {
		grid[i] = make([]string, width)
	}
	for _, file := range files {
		for {
			posX := int(rand.Float64() * float64(width))
			posY := int(rand.Float64() * float64(height))
			if grid[posY][posX] == "" {
				grid[posY][posX] = file.Name()
				break
			}
		}
	}

	for {
		posX := int(rand.Float64() * float64(width))
		posY := int(rand.Float64() * float64(height))
		if grid[posY][posX] == "" {
			Player.posY = posY
			Player.posX = posX
			break
		}
	}
}

func WinScreen() {
	fmt.Printf(`	
__   __           __        __          _
\ \ / /__  _   _  \ \      / /__  _ __ | |
 \ V / _ \| | | |  \ \ /\ / / _ \| '_ \| |
  | | (_) | |_| |   \ V  V / (_) | | | |_|
  |_|\___/ \__,_|    \_/\_/ \___/|_| |_(_)
`)
}

func main() {
	genStartingPath()
	files, err := os.ReadDir(startPath)
	if err != nil {
		log.Fatal("Couldn't read directory")
	}
	s, err := tcell.NewScreen()
	if err = s.Init(); err != nil {
		log.Fatal("Couldn't initialize the screen :(")
	}
	genGrid(files, s)
	s.Clear()
	if err == nil {
		for {
			s.ShowCursor(-1, -1)
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlQ:
					os.Exit(0)
				case tcell.KeyUp:
					if Player.posY != 0 {
						Player.posY--
					}
				case tcell.KeyDown:
					if Player.posY != len(grid)-1 {
						Player.posY++
					}
				case tcell.KeyLeft:
					if Player.posX != 0 {
						Player.posX--
					}
				case tcell.KeyRight:
					if Player.posX != len(grid[0])-1 {
						Player.posX++
					}
				case tcell.KeyEnter:
					if isCorrect {
						var fileNames []string
						for _, file := range files {
							fileNames = append(fileNames, file.Name())
						}
						if files[indexOf(fileNames, selected)].IsDir() {
							startPath = filepath.Join(startPath, selected)
							files, err := os.ReadDir(startPath)
							if err != nil {
								s.Clear()
								log.Fatal("Couldn't read directory")
								os.Exit(0)
							}
							genGrid(files, s)
							isCorrect = false
						} else {
							s.Clear()
							WinScreen()
							os.Exit(0)
						}
					}
				}
			case *tcell.EventResize:
				s.Sync()
			}

			Render(s, files)
			s.ShowCursor(-1, -1)
			s.Sync()
		}
	} else {
		log.Fatal("Couldn't initialize the screen :(")
	}
}
