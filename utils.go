package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
)

func genHelp(title string, items map[string]string) string {
	maxlen := 0
	keys := make([]string, 0, len(items))
	for key, _ := range items {
		keys = append(keys, key)
		if maxlen < len(key) {
			maxlen = len(key)
		}
	}
	help := title + "\n"
	sort.Strings(keys)
	for _, key := range keys {
		help += fmt.Sprintf("%-"+strconv.Itoa(maxlen+8)+"s%s\n", key, items[key])
	}
	return help
}

func TrapSignal(f func(sig os.Signal), sigs ...os.Signal) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, sigs...)
	go func() {
		for sig := range sigCh {
			f(sig)
		}
	}()
}

func SelfDir() string {
	return filepath.Dir(os.Args[0])
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
