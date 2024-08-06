package main

import (
	"fmt"
	"github.com/bezidev/GoAsistent"
	"os"
)

func main() {
	session, err := GoAsistent.WebLogin(os.Getenv("EA_USER"), os.Getenv("EA_PASS"), true, func(username string, sessionToken string) {

	})
	if err != nil {
		fmt.Println(err)
		return
	}
	//session.RefreshWebSession(os.Getenv("EA_PASS"))
	//fmt.Println(session.GetSessionData())
	/*err = session.RefreshSession()
	if err != nil {
		fmt.Println(err)
		return
	}*/
	gradings, err := session.GetAbsences()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(gradings)
}
