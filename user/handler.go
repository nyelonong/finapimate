package user

import (
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nyelonong/finapimate/utils/jsonapi"
)

func (um *UserModule) RegisterHandler(res http.ResponseWriter, req *http.Request) {
	var registerData User

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid Body Request.")
		return
	}

	if err = json.Unmarshal([]byte(body), &registerData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Invalid JSON Request.")
		return
	}

	if err := um.UserRegister(registerData); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Failed to register.")
		return
	}

	jsonapi.SuccessWriter(res, registerData)
}

func (um *UserModule) LoginHandler(res http.ResponseWriter, req *http.Request) {
	var data User

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

	if err := data.UserLogin(um); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Failed to Login.")
		return
	}

	jsonapi.SuccessWriter(res, data)
}
