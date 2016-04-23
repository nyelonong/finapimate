package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nyelonong/finapimate/tx"
	"github.com/nyelonong/finapimate/user"
	"github.com/nyelonong/finapimate/utils"
)

var UserModule *user.UserModule
var TxModule *tx.TxModule

func init() {
	config, err := utils.NewConfig("files/config.ini")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sqlx.Connect("postgres", config.Database.Finmate)
	if err != nil {
		log.Fatalln(err)
	}

	log.SetFlags(log.Lshortfile)

	// parse all html files.
	templates := template.Must(template.New("").ParseGlob("files/*.html"))

	UserModule = user.NewUserModule(db)
	TxModule = tx.NewTxModule(db, UserModule, templates)

	// Construct and implicitly start scheduler.
	if config.Scheduler.Run {
		tx.NewScheduler("@every 5s", TxModule)
	}
}

func main() {
	fmt.Println("FINMATE STARTED")

	// User
	http.HandleFunc("/v1/user/register", UserModule.RegisterHandler)
	http.HandleFunc("/v1/user/login", UserModule.LoginHandler)
	http.HandleFunc("/v1/user/friend/search", UserModule.SearchFriendHandler)
	http.HandleFunc("/v1/user/friend/add", UserModule.AddFriendshandler)
	http.HandleFunc("/v1/user/friend/request", UserModule.FriendRequesthandler)
	http.HandleFunc("/v1/user/friend/approve", UserModule.ApproveFriendshandler)
	http.HandleFunc("/v1/user/friend/list", UserModule.ListFriendHandler)

	// Tx
	http.HandleFunc("/v1/tx/request", TxModule.RequestBorrowHandler)
	http.HandleFunc("/v1/tx/approve", TxModule.ApproveBorrowHandler)
	http.HandleFunc("/v1/tx/decline", TxModule.DeclineBorrowHandler)
	http.HandleFunc("/v1/tx/payment", TxModule.PaymentBorrowHandler)
	http.HandleFunc("/v1/tx/topup", TxModule.TopUpHandler)

	// Notif
	http.HandleFunc("/v1/tx/notif", TxModule.NotifBorrowHandler)

	// History
	http.HandleFunc("/v1/tx/list/borrow", TxModule.BorrowListHandler)
	http.HandleFunc("/v1/tx/list/lend", TxModule.LendListHandler)

	// webview.
	http.HandleFunc("/v1/tx/notif/webview", TxModule.NotifBorrowWebviewHandler)

	// Testing.
	http.HandleFunc("/token", user.TestToken)

	log.Fatal(http.ListenAndServe(":8005", nil))
}
