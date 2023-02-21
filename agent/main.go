//https://gist.github.com/thesubtlety/be6e7ec9c19083473bed4cae11c8160d
//https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getcursorpos

package main

import (
    "log"
    "time"
)

func main() {
    log.Println("Starting agent service")
    //infinite loop for service
    for {
        log.Println("Agent service logic")
        time.Sleep(time.Millisecond * 5000)
    }
}
