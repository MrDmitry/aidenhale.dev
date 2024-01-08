package pages

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type HeadSnippet struct {
	GitRev string
}

func getRevision() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		log.Warnf("failed to detect git revision: %+v", err)
		return strconv.FormatInt(time.Now().Unix(), 10)
	}
	return strings.Trim(string(out), "\n")
}

func NewHeadSnippet() HeadSnippet {
	return HeadSnippet{
		GitRev: getRevision(),
	}
}
