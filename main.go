package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Print("Enter Password: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		p := scanner.Text()

		hp := hashPassword(p)

		hp = strings.ToUpper(hp)
		hashSuffix := hp[5:]

		matchingHashes := getHashs(hp)

		found := false

		for _, item := range matchingHashes {
			i := strings.Split(item, ":")
			hash := i[0]
			leaked := i[1][:len(i[1])-1]

			if hashSuffix == hash {
				fmt.Printf("Found: %s times\n", leaked)
				fmt.Println("Hash:", hp)
				found = true
				break
			}
		}

		if !found {
			fmt.Println("Password have not been")
		}

		fmt.Println("")

	}
}

func hashPassword(pass string) string {
	h := sha1.New()
	h.Write([]byte(pass))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func getHashs(hp string) []string {
	resp, err := http.Get("https://api.pwnedpasswords.com/range/" + hp[:5])

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	return strings.Split(string(body), "\n")
}
