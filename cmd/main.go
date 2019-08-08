package main

import (
    "fmt"
    "os"
    "github.com/dragos240/gopng"
)

func main(){
    if len(os.Args) < 2 || len(os.Args) > 2{
        fmt.Println("Please specify the filename as the only argument.")
        return
    }

    fname := os.Args[1]
    if len(fname) < 5 || fname[len(fname)-4:] != ".png"{
        fmt.Println("Invalid filename!")
        return
    }

    ret := gopng.Parse(os.Args[1])
    if ret == 1{
        fmt.Println("Not a valid PNG file!")
    } else if ret == 2{
        fmt.Println("Parsing error!")
    }
}

