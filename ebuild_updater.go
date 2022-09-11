package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	gitLinks := map[string]string{
		"/usr/local/portage/himozzza_overlay/sys-kernel/zen-sources": "https://github.com/zen-kernel/zen-kernel/releases",
	}

	for n, i := range gitLinks {
		gitLink := i
		gitUpdate := parcingPackageName(gitLink)
		os.Chdir(n)
		files, _ := os.ReadDir(n)
		for _, i := range files {
			if strings.Compare(i.Name(), "Manifest") == 0 {
				continue
			} else if strings.Compare(i.Name(), "files") == 0 {
				continue
			}
			re := regexp.MustCompile("(.*)-[0-9]")
			preName := strings.TrimRight(re.FindString(i.Name()), "1234567890")

			newName := fmt.Sprintf("%s%s.ebuild", preName, gitUpdate)

			if strings.Compare(i.Name(), newName) == 0 {
				fmt.Println("Ebuild уже имеет последнюю версию.")
			} else {
				err := os.Rename(i.Name(), newName)
				if err != nil {
					fmt.Println("Не удалось назначить новое имя.")
				}
			}
			exec.Command("sudo", "ebuild", newName, "digest").Run()
		}
	}
}

func parcingPackageName(gitLink string) string {
	fmt.Printf("Получаем информацию о версии zen-sources на github...\n")
	resp, err := http.Get(gitLink)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile("[0-9.](.*).tar.gz")
	matched := re.FindString(string(body))
	matched = strings.SplitN(matched, "-", -1)[0]
	fmt.Println(matched)

	return matched
}
