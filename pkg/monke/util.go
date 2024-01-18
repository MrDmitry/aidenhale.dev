package monke

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type Configuration struct {
	Gitdir string
}

var Config Configuration

func init() {
	workdir, err := os.Getwd()
	if err != nil {
		panic("failed to detect current work directory, aborting")
	}

	Config = Configuration{
		Gitdir: workdir + "/.git",
	}
}

func (c *Configuration) Validate() error {
	err := GitDirectoryValidator(c.Gitdir)
	if err != nil {
		return err
	}
	return nil
}

func sanitizePath(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] != '/' {
		s = "/" + s
	}
	trailingSlash := s[len(s)-1] == '/'
	s, err := filepath.Abs(s)
	if err != nil {
		return ""
	}
	s = strings.TrimRight(s, "/")

	if trailingSlash {
		return s + "/"
	}
	return s
}

func SanitizeUrl(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	u.Path = sanitizePath(u.Path)
	return u.String()
}

func gitExec(gitdir string, params ...string) (string, error) {
	args := append([]string{"--git-dir", gitdir}, params...)
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(out), "\n"), nil
}

func GitRevision(gitdir ...string) (string, error) {
	workdir := Config.Gitdir
	if len(gitdir) > 0 {
		workdir = gitdir[0]
	}
	out, err := gitExec(workdir, "rev-parse", "HEAD")
	if err != nil {
		log.Warnf("failed to detect git revision: %+v", err)
		return strconv.FormatInt(time.Now().Unix(), 10), err
	}
	return out, nil
}

func GitLastLogTimeISO(path string, gitdir ...string) (time.Time, error) {
	workdir := Config.Gitdir
	if len(gitdir) > 0 {
		workdir = gitdir[0]
	}
	out, err := gitExec(workdir, "log", "-1", "--pretty=%aI", "--", path)
	if err != nil {
		log.Warnf("failed to detect git log for %s: %+v", path, err)
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, out)
}
