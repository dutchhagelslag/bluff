package server

import (
    "fmt"
    "os"
    "time"
)

func spin(c chan string){
    current := 0
    speed := time.Duration(300) // ms

    pattern := "~o~"
    back := "-----------"

    for {
        select{
            case <-c: // handle and print error in future for debug
                return

            default:
                // from the right
                if(current < len(pattern)){
                    fmt.Printf("%s%s", pattern[len(pattern)-current:], back[current:])
                // over the left
                } else if(current == len(back)){
                    for i := 0; i < len(pattern); i++{
                        dif := len(pattern)-i
                        fmt.Printf("%s%s", back[:current - dif], pattern[:dif])
                        time.Sleep(speed * time.Millisecond)
                        fmt.Fprintf(os.Stdout, "\r")
                    }
                    current = 0
                    continue

                // moving in the mid
                } else {
                     fmt.Printf("%s%s%s", back[:current-len(pattern)], pattern, back[current:])
                }

                current = current % len(back)
                current++

                fmt.Fprintf(os.Stdout, "\r")
                time.Sleep(speed * time.Millisecond)
        }
    }
}
