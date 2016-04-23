package tx

import (
	"fmt"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nyelonong/finapimate/utils/jsonapi"
)

func (tm *TxModule) RequestBorrowHandler(res http.ResponseWriter, req *http.Request) {
	var txData []Transaction

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	if err := tm.RequestBorrow(txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Failed to borrow.")
		return
	}

	jsonapi.SuccessWriter(res, txData)
}

// input : tx_id
// status 2
// amount
func (tm *TxModule) ApproveBorrowHandler(res http.ResponseWriter, req *http.Request) {
	var txData []Transaction

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	if err := tm.ChangeStatusTx(txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Approve failed.")
		return
	}

	jsonapi.SuccessWriter(res, txData)
}

// input : tx_id
// status : 3
func (tm *TxModule) DeclineBorrowHandler(res http.ResponseWriter, req *http.Request) {
	var txData []Transaction

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	if err := tm.ChangeStatusTx(txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Decline failed.")
		return
	}

	jsonapi.SuccessWriter(res, txData)
}

// input : tx_id
// status 4
// amount
func (tm *TxModule) PaymentBorrowHandler(res http.ResponseWriter, req *http.Request) {
	var txData []Transaction

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	if err := tm.ChangeStatusTx(txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Approve failed.")
		return
	}

	jsonapi.SuccessWriter(res, txData)
}

// input : borrower_id
func (tm *TxModule) BorrowListHandler(res http.ResponseWriter, req *http.Request) {
	var txData Transaction

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	datas, err := tm.ListBorrow(txData)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Not found.")
		return
	}

	jsonapi.SuccessWriter(res, datas)
}

func (tm *TxModule) BorrowListWebviewHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid, err := strconv.Atoi(r.Form.Get("user_id")) // this is borrower id.
	if err != nil {
		log.Println(err)
		return
	}

	txData := Transaction{
		BorrowerID: int64(uid),
	}

	datas, err := tm.ListBorrow(txData)
	if err != nil {
		log.Println(err)
		jsonapi.ErrorsWriter(w, 400, "Error Server.")
		return
	}

	data := struct {
		List []Transaction
	}{
		List: datas,
	}

	err = tm.renderTemplate(w, r, "listLender", data)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// input : lender_id
func (tm *TxModule) LendListHandler(res http.ResponseWriter, req *http.Request) {
	var txData Transaction

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	datas, err := tm.ListLend(txData)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Not found.")
		return
	}

	jsonapi.SuccessWriter(res, datas)
}

// input : lender_id
func (tm *TxModule) NotifBorrowHandler(res http.ResponseWriter, req *http.Request) {
	var txData Transaction

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &txData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	datas, err := tm.NotifBorrow(txData)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Not found.")
		return
	}

	jsonapi.SuccessWriter(res, datas)
}

func (tm *TxModule) TopUpHandler(res http.ResponseWriter, req *http.Request) {
	var data TopUp

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &data); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	if err := data.UserTopUp(tm); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Failed to Login.")
		return
	}

	jsonapi.SuccessWriter(res, data)
}

func (tm *TxModule) NotifBorrowWebviewHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	uid, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		fmt.Println(err)
		return
	}

	txData := Transaction{
		LenderID: int64(uid),
	}

	datas, err := tm.NotifBorrow(txData)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(w, 400, "Not found.")
		return
	}

	data := struct {
		List []Transaction
	}{
		List: datas,
	}

	// jsonapi.SuccessWriter(w, datas)
	err = tm.renderTemplate(w, r, "listBorrower", data)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func (tm *TxModule) renderTemplate(w http.ResponseWriter, r *http.Request, tname string, data interface{}) error {

	template := tm.templates.Lookup(tname)
	if template != nil {
		err := template.Execute(w, data)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}