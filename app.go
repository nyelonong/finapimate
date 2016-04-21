package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nyelonong/finapimate/user"
	"github.com/nyelonong/finapimate/utils"
)

var UserModule *user.UserModule

func init() {
	config, err := utils.NewConfig("files/config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sqlx.Connect("postgres", config.Database.Finmate)
	if err != nil {
		log.Fatalln(err)
	}

	UserModule = user.NewUserModule(db)
}

func main() {
	fmt.Println("FINMATE STARTED")

	http.HandleFunc("/v1/user/register", UserModule.RegisterHandler)
	http.HandleFunc("/v1/user/login", UserModule.LoginHandler)
	http.HandleFunc("/v1/user/friend/search", UserModule.SearchFriendHandler)
	http.HandleFunc("/v1/user/friend/add", UserModule.AddFriendshandler)
	http.HandleFunc("/v1/user/friend/request", UserModule.FriendRequesthandler)
	http.HandleFunc("/v1/user/friend/approve", UserModule.ApproveFriendshandler)

	log.Fatal(http.ListenAndServe(":8005", nil))
}
