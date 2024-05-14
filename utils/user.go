package utils

import (
	"fmt"
	"os"
	"os/user"
)

func main() {
	fmt.Println(user.Current()) //return current username
	fmt.Println(os.Hostname())  //return the hostname(domain)
}
