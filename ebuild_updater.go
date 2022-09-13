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
		"/usr/local/portage/himozzza_overlay/sys-kernel/zen-sources": "https://github.com/zen-kernel/zen-kernel/tags",
	}

	for n, i := range gitLinks {
		re := regexp.MustCompile(`[\w\d-_]+$`)
		gitName := re.FindString(n) // Имя запрашиваемого репозитория.

		gitVersion := parcingPackageName(i, gitName)
		os.Chdir(n)
		files, _ := os.ReadDir(n)
		for _, i := range files {
			if strings.Compare(i.Name(), "Manifest") == 0 {
				continue
			} else if strings.Compare(i.Name(), "files") == 0 {
				continue
			}

			newName := fmt.Sprintf("%s-%s.ebuild", gitName, gitVersion)

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

func parcingPackageName(gitLink, gitName string) string {
	fmt.Printf("Получаем информацию о версии %s на GitHub...\n", gitName)
	resp, err := http.Get(gitLink)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile(`/[\w\d-.+]*.tar.gz`)
	matched := re.FindString(string(body)) // Пакета tar.gz с последней версией.
	re, err = regexp.Compile(`[\d]+\.[\d.]+`)
	if err != nil {
		fmt.Println("Не удалость получить версию пакета.", err)
		os.Exit(0)
	}

	matched = strings.TrimRight(re.FindString(matched), ".") // Форматирование строки версией tar.gz для приведения к виду "5.19.8".
	fmt.Println(matched)

	return matched
}
