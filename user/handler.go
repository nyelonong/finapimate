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

func (um *UserModule) SearchFriendHandler(res http.ResponseWriter, req *http.Request) {
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

	datas, err := um.SearchFriend(data)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Not found.")
		return
	}

	jsonapi.SuccessWriter(res, datas)
}

func (um *UserModule) AddFriendshandler(res http.ResponseWriter, req *http.Request) {
	var data []UserRelation

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

	if err := um.AddFriends(data); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Failed add friends.")
		return
	}

	jsonapi.SuccessWriter(res, data)
}

func (um *UserModule) ApproveFriendshandler(res http.ResponseWriter, req *http.Request) {
	var data []UserRelation

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

	if err := um.ApproveFriends(data); err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Failed approve friends.")
		return
	}

	jsonapi.SuccessWriter(res, data)
}

func (um *UserModule) FriendRequesthandler(res http.ResponseWriter, req *http.Request) {
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

	datas, err := um.FriendRequest(data)
	if err != nil {
		fmt.Println(err)
		jsonapi.ErrorsWriter(res, 400, "Failed getting friend request.")
		return
	}

	jsonapi.SuccessWriter(res, datas)
}
