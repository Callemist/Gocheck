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

		hp, err := hashPassword(p)

		if err != nil {
			fmt.Println("Error hashing password: ", err)
			return
		}

		hp = strings.ToUpper(hp)
		hashSuffix := hp[5:]

		matchingHashes, err := getHashs(hp)

		if err != nil {
			fmt.Println("Error getting matching hashes: ", err)
			return
		}

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

func hashPassword(pass string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(pass))

	if err != nil {
		return "", err
	}

	bs := h.Sum(nil)
	return hex.EncodeToString(bs), nil
}

func getHashs(hp string) ([]string, error) {
	resp, err := http.Get("https://api.pwnedpasswords.com/range/" + hp[:5])

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(body), "\n"), nil
}
