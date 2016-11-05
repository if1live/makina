package hitomi

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func ExecuteHaru(config Config, code int) (bool, string) {
	cmd := config.HaruExecutable
	args := []string{
		fmt.Sprintf("-id=%d", code),
		"-service=hitomi",
		"-cmd=download",
	}
	out, err := exec.Command(cmd, args...).CombinedOutput()

	// dump stdout/stderr
	stderrs := []string{}
	if err != nil {
		if _, ok := err.(*exec.Error); ok {
			stderrs = append(stderrs, err.Error())
		}
	}
	stdouts := strings.Split(string(out[:]), "\n")

	if config.ShowLog {
		for _, line := range stderrs {
			log.Println(line)
		}
		for _, line := range stdouts {
			log.Println(line)
		}
	}

	// 공백이 아닌 가장 마지막 출력 찾기
	// 거기에 파일명이 있을거다
	lastStdout := ""
	for i := len(stdouts) - 1; i >= 0; i-- {
		if len(stdouts[i]) > 0 {
			lastStdout = stdouts[i]
			break
		}
	}

	re := regexp.MustCompile(` (/.*\.zip)`)
	for _, m := range re.FindAllStringSubmatch(lastStdout, -1) {
		return true, m[1]
	}
	return false, ""
}
