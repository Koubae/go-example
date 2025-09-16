package assignments

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type AddressInfo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func Makejson() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf(">>> Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Printf(">>> Address: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	addressInfo := AddressInfo{
		Name:    name,
		Address: address,
	}
	marshalled, err := json.Marshal(addressInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(marshalled))
}
