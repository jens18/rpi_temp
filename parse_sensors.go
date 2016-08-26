// Go bindings for libsensors.so from the lm-sensors project via cgo: https://github.com/md14454/gosensors
//
// Note: This program does not have a C library dependency and instead calls the 'sensors' command and 
// extracts the CPU temperature ('screen scraping').

package main

import "fmt"
import "regexp"
import "log"
import "bufio"
import "strings"
import "os/exec"

// Detect/discover the version of the sensors command and machine type (uname -m: x86_64 or armv6l or armv7l)
func init() {
	sensorsCmd := exec.Command("/usr/bin/sensors", "-v")
	// read output into a byte slice
	sensorsOut, err := sensorsCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sensors version: %s\n", sensorsOut)

	// set the string constants to match based on the version of the sensors command
}

func main() {
	// each argument is represented as a separate string, assign to Cmd sensorsCmd Cmd object
	sensorsCmd := exec.Command("/usr/bin/sensors", "-u", "-A")

	stdoutSensorsCmd, err := sensorsCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := sensorsCmd.Start(); err != nil {
		log.Fatal(err)
	}

     	scanner := bufio.NewScanner(stdoutSensorsCmd)

	for scanner.Scan() {
		str := scanner.Text()

		//
		matched, _ := regexp.Match("CPUTIN:", []byte(str))

		if matched {
			// fmt.Println("found CPUTIN")

			// read next line containing "temp2_input: 46.500" name value pair
			scanner.Scan()
			str = scanner.Text()

			// fmt.Println("next line = ", str)
			matched, _ := regexp.Match("temp2_input:", []byte(str))

			if matched {
				fmt.Println("found temp2_input", str)

				// Example: temp2_input: 46.500

				result := strings.Fields(str)
				// Display all elements.
				for i := range result {
					fmt.Println(result[i])
				}
				// the second element contains the temperature
				fmt.Println("temp = ", result[1])
				
				// format: 42.000
				temp := strings.Split(result[1], ".")

				leftSide := temp[0]
//				rightSide := temp[1]

				// slice expression to extract the first character after the '.'
				// format: 42.0
				fmt.Printf("%s.%s\n", leftSide, temp[1][0:1])
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
    }
}
	


