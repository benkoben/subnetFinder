package main

import (
    "fmt"
    "strings"
    "strconv"
)

func binarysum(bits ...int) int {
    var summary int
    for _, i := range bits {
        summary += i
    }
    return summary
}

func main(){
	var input = []int{24,23,22,19,8,1}
    
	for _, in := range input {
        hostBits := 32 - in
        binary := fmt.Sprintf("%v", strings.Repeat("1", hostBits))
        fmt.Println("24 mask in binary: ", binary)
        decimal, _ := strconv.ParseInt(binary,2,64)
        fmt.Println("24 mask in decimal: ", decimal + 1)
    }
}
