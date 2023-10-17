package pages

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type HeadSnippet struct {
	PageTitle string
	GitRev    string
}

func NewHeadSnippet(t string) HeadSnippet {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	rev := strconv.FormatInt(time.Now().Unix(), 10)
	out, err := cmd.Output()
	if err != nil {
		log.Warnf("Failed to detect git revision: %+v", err)
	} else {
		rev = strings.Trim(string(out), "\n")
	}
	return HeadSnippet{
		PageTitle: t,
		GitRev:    rev,
	}
}
