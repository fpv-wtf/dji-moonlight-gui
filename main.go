package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/jchv/go-webview2"
)

const WindowWidth = 768
const WindowHeight = 768

type MoonlightManager struct {
	ConsoleOutputFunc func(string)
	RunningStateFunc  func(bool)
	RunningCmd        *exec.Cmd
	RunningCmdLock    sync.Mutex
}

func (m *MoonlightManager) GetGames() ([]string, error) {
	m.RunningCmdLock.Lock()
	defer m.RunningCmdLock.Unlock()

	if m.RunningCmd != nil {
		return nil, errors.New("already running a command")
	}

	cmd := m.getCommand()
	m.RunningCmd = cmd

	cmd.Args = append(cmd.Args, "list")

	m.RunningStateFunc(true)
	m.ConsoleOutputFunc(fmt.Sprintf("> %s", strings.Join(cmd.Args, " ")))
	out, _ := cmd.CombinedOutput()

	outString := strings.TrimSpace(string(out))
	outStrings := strings.Split(outString, "\n")

	m.ConsoleOutputFunc(string(out))

	gameStrings := outStrings[2:]

	if strings.Index(gameStrings[0], "1. ") != 0 {
		m.RunningCmd = nil
		m.RunningStateFunc(false)
		return nil, errors.New("failed to get games")
	}

	for i, gameString := range gameStrings {
		periodPos := strings.Index(gameString, ". ")
		gameStrings[i] = gameString[periodPos+2:]
	}

	m.RunningCmd = nil
	m.RunningStateFunc(false)
	return gameStrings, nil
}

func (m *MoonlightManager) Pair() error {
	m.RunningCmdLock.Lock()
	if m.RunningCmd != nil {
		return errors.New("already running a command")
	}
	m.RunningCmdLock.Unlock()

	go func() {
		m.RunningCmdLock.Lock()
		cmd := m.getCommand()
		m.RunningCmd = cmd

		cmd.Args = append(cmd.Args, "pair")

		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		combined := io.MultiReader(stdout, stderr)
		reader := bufio.NewReader(combined)

		m.RunningStateFunc(true)
		m.ConsoleOutputFunc(fmt.Sprintf("> %s", strings.Join(cmd.Args, " ")))
		cmd.Start()
		m.RunningCmdLock.Unlock()

		for {
			out, err := reader.ReadString('\n')
			if err != nil {
				log.Println(err)
				break
			}

			outString := strings.TrimSpace(out)
			m.ConsoleOutputFunc(outString)
		}

		m.RunningCmdLock.Lock()
		m.RunningCmd = nil
		m.RunningCmdLock.Unlock()
		m.RunningStateFunc(false)
	}()

	return nil
}

func (m *MoonlightManager) Unpair() error {
	m.RunningCmdLock.Lock()
	defer m.RunningCmdLock.Unlock()

	if m.RunningCmd != nil {
		return errors.New("already running a command")
	}

	cmd := m.getCommand()
	m.RunningCmd = cmd

	cmd.Args = append(cmd.Args, "unpair")

	m.RunningStateFunc(true)
	m.ConsoleOutputFunc(fmt.Sprintf("> %s", strings.Join(cmd.Args, " ")))
	out, _ := cmd.CombinedOutput()

	m.ConsoleOutputFunc(string(out))
	m.RunningCmd = nil
	m.RunningStateFunc(false)

	return nil
}

func (m *MoonlightManager) Quit() error {
	m.RunningCmdLock.Lock()
	defer m.RunningCmdLock.Unlock()

	if m.RunningCmd != nil {
		return errors.New("already running a command")
	}

	cmd := m.getCommand()
	m.RunningCmd = cmd

	cmd.Args = append(cmd.Args, "quit")
	m.RunningStateFunc(true)
	m.ConsoleOutputFunc(fmt.Sprintf("> %s", strings.Join(cmd.Args, " ")))
	out, _ := cmd.CombinedOutput()

	m.ConsoleOutputFunc(string(out))
	m.RunningCmd = nil
	m.RunningStateFunc(false)

	return nil
}

func (m *MoonlightManager) ForceStop() {
	m.RunningCmdLock.Lock()
	defer m.RunningCmdLock.Unlock()

	if m.RunningCmd != nil {
		m.RunningCmd.Process.Kill()
		m.RunningCmd = nil
		m.RunningStateFunc(false)
	}
}

type StreamGameResolutions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type StreamGameParams struct {
	Bitrate    int                   `json:"bitrate"`
	Fps        int                   `json:"fps"`
	Game       string                `json:"game"`
	Mode       string                `json:"mode"`
	Resolution StreamGameResolutions `json:"resolution"`
}

func (m *MoonlightManager) StreamGame(params StreamGameParams) error {
	m.RunningCmdLock.Lock()
	if m.RunningCmd != nil {
		return errors.New("already running a command")
	}
	m.RunningCmdLock.Unlock()

	go func() {
		m.RunningCmdLock.Lock()
		cmd := m.getCommand()
		m.RunningCmd = cmd

		cmd.Args = append(
			cmd.Args,
			"stream",
			"-app",
			params.Game,
			"-platform",
			params.Mode,
			"-bitrate",
			fmt.Sprintf("%d", params.Bitrate),
			"-fps",
			fmt.Sprintf("%d", params.Fps),
			"-w",
			fmt.Sprintf("%d", params.Resolution.Width),
			"-h",
			fmt.Sprintf("%d", params.Resolution.Height),
			"-localaudio",
			"-viewonly",
		)

		fmt.Println(cmd.Args)

		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		combined := io.MultiReader(stdout, stderr)
		reader := bufio.NewReader(combined)

		m.RunningStateFunc(true)
		cmd.Start()
		m.RunningCmdLock.Unlock()

		for {
			println("reading")
			out, err := reader.ReadString('\n')
			if err != nil {
				log.Println(err)
				break
			}

			m.ConsoleOutputFunc(out)
		}

		m.RunningCmdLock.Lock()
		m.RunningCmd = nil
		m.RunningCmdLock.Unlock()
		m.RunningStateFunc(false)
	}()

	return nil
}

func (m *MoonlightManager) getCommand() *exec.Cmd {
	cwd, _ := os.Getwd()

	moonlightDir := cwd + "\\moonlight"
	moonlightExe := moonlightDir + "\\moonlight.exe"

	cmd := exec.Command(moonlightExe)
	cmd.Dir = moonlightDir

	return cmd
}

func main() {
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Center: true,
			Height: WindowHeight,
			IconId: 0,
			Title:  "DJI Moonlight (v1.0.2)",
			Width:  WindowWidth,
		},
	})

	if w == nil {
		log.Fatalln("Failed to load webview.")
	}

	m := &MoonlightManager{}
	m.ConsoleOutputFunc = func(s string) {
		s = strings.TrimSpace(s)
		log.Println(s)

		js, _ := json.Marshal(s)
		jsString := string(js)

		w.Dispatch(func() {
			w.Eval(`window.e_receiveConsoleOutput(` + jsString + `);`)
		})
	}
	m.RunningStateFunc = func(running bool) {
		w.Dispatch(func() {
			w.Eval(`window.e_receiveRunningState(` + fmt.Sprintf("%v", running) + `);`)
		})
	}

	// Get current directory
	dir, _ := os.Getwd()

	w.Bind("b_getGames", m.GetGames)
	w.Bind("b_pair", m.Pair)
	w.Bind("b_unpair", m.Unpair)
	w.Bind("b_forceStop", m.ForceStop)
	w.Bind("b_streamGame", m.StreamGame)
	w.Bind("b_quit", m.Quit)

	w.SetSize(WindowWidth, WindowHeight, webview2.HintFixed)
	w.Navigate("file://" + dir + "/assets/index.html")

	w.Run()
}
