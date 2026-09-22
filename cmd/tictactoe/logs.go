package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type NullWriter struct{}

func (w NullWriter) Write(bytes []byte) (n int, err error) {
	return len(bytes), nil
}

func setupLogs() func() {
	if debug, ok := os.LookupEnv("DEBUG"); ok && debug != "" {
		logFile, err := os.OpenFile("game.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create log file")
			os.Exit(1)
		}
		slog.SetDefault(slog.New(slog.NewTextHandler(logFile, nil)))
		return func() { logFile.Close() }
	}
	var w io.Writer
	w = NullWriter{}
	slog.SetDefault(slog.New(slog.NewTextHandler(w, nil)))
	return func() {}
}
