package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mritd/promptx"
)

func main() {
	cmd := exec.Command("kubectl", "config", "current-context")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	checkAndExit(err)
	currentContext := strings.TrimSpace(string(b))

	if len(os.Args) > 1 {
		cmd = exec.Command("kubectl", "config", "set-context", currentContext, "--namespace="+os.Args[1])
		cmd.Stdout = os.Stdout
		checkAndExit(cmd.Run())
		fmt.Printf("Kubernetes namespace switch to %s.\n", os.Args[1])
	} else {
		cmd = exec.Command("kubectl", "get", "ns", "-o", "template", "--template", "{{ range .items }}{{ .metadata.name }} {{ end }}")
		b, err = cmd.Output()
		checkAndExit(err)
		allNameSpace := strings.Fields(string(b))

		cfg := &promptx.SelectConfig{
			ActiveTpl:    "»  {{ . | cyan }}",
			InactiveTpl:  "  {{ . | white }}",
			SelectPrompt: "NameSpace",
			SelectedTpl:  "{{ \"» \" | green }}{{\"NameSpace:\" | cyan }} {{ . }}",
			DisPlaySize:  9,
			DetailsTpl:   ` `,
		}
		s := &promptx.Select{
			Items:  allNameSpace,
			Config: cfg,
		}

		selectNameSpace := allNameSpace[s.Run()]

		cmd = exec.Command("kubectl", "config", "set-context", currentContext, "--namespace="+selectNameSpace)
		cmd.Stdout = os.Stdout
		checkAndExit(cmd.Run())
		fmt.Printf("Kubernetes namespace switch to %s.\n", selectNameSpace)
	}
}

func checkErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func checkAndExit(err error) {
	if !checkErr(err) {
		os.Exit(1)
	}
}
