package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func parseTemplate(ctx *context, t string,
	b *strings.Builder,
	fs io.Writer,
) error {
	w := io.MultiWriter(b, fs)
	tmpl, err := template.New("").Parse(t)
	if err != nil {
		return fmt.Errorf("parse template faild %w", err)
	}
	if err := tmpl.Execute(w, ctx); err != nil {
		return fmt.Errorf("execute template faild %w", err)
	}
	return nil
}

func mkdir(s string) error {
	if err := os.MkdirAll(s, 0775); err != nil {
		return fmt.Errorf("cannot create directory, %v", err)
	}
	return nil
}

func toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	c := rune(s[0])

	if !unicode.IsLetter(c) {
		return ""
	}

	clean := ""
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			clean += string(r)
		} else if r == '_' || r == ' ' || r == '-' {
			clean += string(' ')
		}
	}

	words := strings.Fields(clean)
	title := cases.Title(language.English)
	for i, word := range words {
		words[i] = title.String(word)
	}
	return strings.Join(words, "")
}

func findGoMod(s string, maxDepth int) (string, string, error) {
	root := filepath.VolumeName(s)
	current := s
	for depth := 0; depth < maxDepth; depth++ {
		goModPath := filepath.Join(current, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			module, err := readGoModModule(goModPath)
			return current, module, err
		}
		parent := filepath.Dir(current)
		if parent == root {
			break
		}
		current = parent
	}
	return "", "", errors.New("cannot find go.mod file, you must initial a golang project first")
}

func readGoModModule(f string) (string, error) {
	fs, err := os.Open(f)
	if err != nil {
		return "", err
	}
	defer fs.Close()
	var line string
	scanner := bufio.NewScanner(fs)
	if scanner.Scan() {
		line = scanner.Text()
	}
	if scanner.Err() != nil {
		return "", err
	}
	if !strings.HasPrefix(line, "module") {
		return "", fmt.Errorf("invalid go.mod format at %s", f)
	}
	return strings.TrimPrefix(line, "module "), nil
}

func calculateModule(module string, name string, root string, dest string) (string, error) {
	if len(dest) == 0 {
		return module + "/" + name, nil
	}
	root = filepath.Clean(root)
	dest = filepath.Clean(dest)

	sep := string(filepath.Separator)
	if !strings.HasSuffix(root, sep) {
		root += sep
	}
	if !strings.HasSuffix(dest, sep) {
		dest += sep
	}
	index := strings.Index(dest, root)
	if index < 0 {
		return "", fmt.Errorf("you must specify outPath in golang project")
	}

	tails := strings.Trim(dest[index+len(root):], sep)
	tails = strings.ReplaceAll(tails, sep, "/")
	return module + "/" + tails, nil
}

func checkPath(path string) (string, error) {
	if len(path) == 0 {
		return os.Getwd()
	}
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return abs, err
	}

	fi, err := os.Stat(abs)

	if os.IsNotExist(err) {
		return abs, mkdir(abs)
	}

	if !fi.IsDir() {
		return "", fmt.Errorf("%s is not a directory", abs)
	}
	return abs, nil
}
