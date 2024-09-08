package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"yato/config"
	"yato/screens"
	"yato/utils"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	authorizationChan = make(chan struct{})
	codeVerifier      string
	errorChan         = make(chan error)
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf(err.Error())
	}

	config := config.GetConfig()
	if config.MyAnimeList.AccessToken == "" {
		go StartOAuthFlow()
		select {
		case <-authorizationChan: // Authorization successful
		case err := <-errorChan:
			log.Fatalf("Unable to authenticate: %s", err)
		}
	}

	StartApp()
}

func StartApp() {
	p := tea.NewProgram(screens.Initialize(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}
}

func StartOAuthFlow() {
	var err error

	codeVerifier, err = utils.GetNewCodeVerifier()
	if err != nil {
		log.Fatalf("failed to generate code verifier: %s", err)
	}

	url := utils.GetOAuthURL(codeVerifier)

	server := &http.Server{Addr: ":42069"}
	http.HandleFunc("/authenticate", handleOAuthCallback)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	if err := openBrowser(url); err != nil {
		log.Printf("failed to open browser: %v. Visit %s in your browser to authenticate.", err, url)
	}

	select {
	case <-authorizationChan:
	case err := <-errorChan:
		log.Fatalf("Unable to authenticate: %s", err)
	}

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}
}

func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		errorChan <- fmt.Errorf("user cancelled authentication")
		http.Error(w, "missing code query parameter. user cancelled authentication", http.StatusBadRequest)
		return
	}

	malConfig, err := utils.ExchangeToken(code, codeVerifier)
	if err != nil {
		errorChan <- fmt.Errorf("failed to exchange token: %w", err)
		http.Error(w, "failed to exchange token", http.StatusInternalServerError)
		return
	}

	config.GetConfig().MyAnimeList = *malConfig
	if err := config.SaveConfig(); err != nil {
		errorChan <- fmt.Errorf("failed to save config: %w", err)
		http.Error(w, "failed to save config", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Authentication successful! You can now close this tab."))
	authorizationChan <- struct{}{}
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
