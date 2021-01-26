package git

import (
	"os/exec"
	"strings"
)


// run get
func ExecShell(command string, arg ...string) (out string, err error) {
	var Stdout []byte
	cmd := exec.Command(command, arg...)
	Stdout, err = cmd.CombinedOutput()
	out = string(Stdout)
	return
}

// Repository
func Repo() string {
	res, err := ExecShell("/bin/sh", "-c", "git remote -v")
	if err != nil {
		return UNKNOWN
	}
	repo := res[strings.Index(res, ":")+1 : strings.Index(res, ".git")]
	if repo == "" {
		return UNKNOWN
	}
	return repo
}

// Current Branch
func Branch()string{
	res, err := ExecShell("/bin/sh", "-c", "git branch")
	if err != nil {
		return UNKNOWN
	}
	list := strings.Split(res, "\n")
	for _, v := range list {
		if strings.HasPrefix(v, "*") {
			res = v[strings.Index(v, "*")+2:]
			return res
		}
	}
	return UNKNOWN
}

// Current Commit
func Commit()string{
	commit, err := ExecShell("/bin/sh", "-c", "git rev-parse --short HEAD")
	if err != nil || len(commit) == 0{
		return UNKNOWN
	}

	return commit[:len(commit)-1]
}

// Author of a specified commit
func Author(commit string)string{
	arg := "git log --pretty=format:“%an” "+ commit +" -1"
	author, err := ExecShell("/bin/sh", "-c", arg)
	if err != nil || len(commit) == 0{
		return UNKNOWN
	}
	return author
}

// Get message of a commit
func Message(commit string)string{
	if len(commit) == 0{
		return UNKNOWN
	}
	arg := "git log --pretty=format:“%s” "+ commit +" -1"
	message, err := ExecShell("/bin/sh", "-c", arg)
	if err != nil || len(commit) == 0{
		return UNKNOWN
	}
	return message
}
