package acccheck

import (
	"fmt"
	"regexp"
)

// Output: print"\n        SUCCESS.... connected to $singleIp with username:'$singleUser' and password:'$singlePass'\n";
func Interpret(resultBytes []byte) []byte {
	ipRegex, _ := regexp.Compile(`\b(?:(?:2(?:[0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9])\.){3}(?:(?:2([0-4][0-9]|5[0-5])|[0-1]?[0-9]?[0-9]))\b`)
	usrPassRegex, _ := regexp.Compile(`'[^']*'`)

	ip := ipRegex.FindStringSubmatch(string(resultBytes))[0]
	userPass := usrPassRegex.FindAllStringSubmatch(string(resultBytes), 2)

	fmt.Println(ip)
	fmt.Println(userPass[0])
	fmt.Println(userPass[1])

	return []byte("aababa")
}

func Test() {
	out := ""

	for i := 1; i < 5; i++ {
		out += fmt.Sprintf("\n        SUCCESS.... connected to %s with username:'%s' and password:'%s'\n", "127.0.0.1", "phil", "password123")
	}

	fmt.Println(out)

	Interpret([]byte(out))
}
