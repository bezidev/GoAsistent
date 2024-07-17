package main

import (
	"fmt"
	"github.com/bezidev/GoAsistent"
	"os"
)

func main() {
	session, err := GoAsistent.Login(os.Getenv("EA_USER"), os.Getenv("EA_PASS"), true)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(session.GetSessionData())
	/*err = session.RefreshSession()
	if err != nil {
		fmt.Println(err)
		return
	}*/
	gradings, err := session.GetNotifications()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(gradings)
}
