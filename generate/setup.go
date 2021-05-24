package generate

import (
	"os/exec"
	"strings"
)

func GoList() (string, error) {
	out, err := exec.Command("go", "list").Output()
	return strings.TrimSpace(string(out)), err
}
