package main

import (
	// "bufio"
	"fmt"
	// "os"
	// "strconv"
	// "strings"
)

func main() {


    // Input for number word

    // var yaa int
    // fmt.Println("Enter")
    // fmt.Scanln(&yaa)
    // fmt.Println(yaa)

    // name := "yash";
    // fmt.Println("My name: ", name);


    // Input for sentences

    // reader := bufio.NewReader(os.Stdin);
    // fmt.Println("Say your name: ");
    // yourName, _ := reader.ReadString('\n');
    // fmt.Println("Your name: ",yourName)



    // fmt.Println("Say your number: ");

    // input, _ := reader.ReadString('\n');

    // number, err := strconv.ParseFloat(strings.TrimSpace(input), 64)

    // if err != nil {
    //     fmt.Println(err);
    // } else {
    //     fmt.Println(number+4);
    // }

    var array = [3]int {3, 5, 9}

    var slice = []int{898, 899};
    slice = append(slice, 8, 9);

    fmt.Println(array, slice)

    hh := make([]int, 4)

    hh[0] = 9;
    hh[3] = 9;

    hh = append(hh, 8)

    fmt.Println(hh);

    var maps = make(map[string]int)
    mm := make(map[string]int)

    maps["jdfn"] = 8;
    maps["hfuf"] = 8;
    maps["fh"] = 8;

    delete(maps, "fh")
    
    fmt.Println(maps, mm);

    // range
    for key, value := range maps {
        fmt.Println(key, value);
    }

    if n := 3; n < 5 {
        fmt.Println(n);
    } 


    uu := User{"yash", 9,"yoo"};
    uu.greeting();
    fmt.Println(uu.Age);


};

type User struct {
    Name string
    Age int
    Bio string
}

// method

func (u *User) greeting() {
    u.Age = 22;
    fmt.Println("Hello", u.Name);
}