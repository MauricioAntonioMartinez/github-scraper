package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Users []User

type User struct {
	Id int "json:id"
	Login string "json:login"
}


func main(){
 resp ,err := http.Get("https://api.github.com/users/torvalds/followers?per_page=100&page=1")
	if err != nil {
		log.Fatal(err)
	}

	 defer resp.Body.Close()
	fmt.Println("This is a test 123")
	fmt.Println("Another print, because they ask")
	//  body,_ := ioutil.ReadAll(resp.Body)

	 users := &Users{}
	 decoder := json.NewDecoder(resp.Body)
	 err = decoder.Decode(users)
	 
	if err != nil {
		log.Fatal(err)
	}

	for _,user := range *users {
		urs := user.fetchUsers(1)
		fmt.Println(urs)
	}


}


func (u User) fetchUsers(page int64) Users {
	urlFollowers := fmt.Sprintf("https://api.github.com/users/%s/followers?per_page=100&page=%d",u.Login,page)
	// url_following := fmt.Sprintf("https://api.github.com/users/%d/following?per_page=100&page=1",u.Login)

	resp ,err := http.Get(urlFollowers)
	if err != nil {
		log.Fatal(err)
	}
	 defer resp.Body.Close()

	 users := &Users{}
	 decoder := json.NewDecoder(resp.Body)
	 err = decoder.Decode(users)
	 
	if err != nil {
		log.Fatal(err)
	}

	return *users

}
