package main

import (
	_ "github.com/lib/pq"
	"math/rand"
	"time"

	//"github.com/davecgh/go-spew/spew"
	_ "github.com/zellyn/kooky"
	_ "github.com/zellyn/kooky/allbrowsers" // register cookie store finders!
)


func main(){
	tabbles_conn()
	rand.Seed(time.Now().UnixNano())
	_, err := db.Exec("DELETE FROM transactions")
	if err != nil {
		errorLog.Println(err)
		return
	}
	_, err = db.Exec("DELETE FROM sell_was")
	if err != nil {
		errorLog.Println(err)
		return
	}
	get_cookies_from_browser()
	add_30_cases()

	go transaction_add()
	for i :=0;i<1;i++{
		get_from_steam()
		get_dm_items()
		<-quit
		not_added_sell()
		link_items()
		write_to_google()
	}
}



func add_30_cases(){
	db.Exec("INSERT INTO public.sell_was(id, price, market, name, status) VALUES ('1488', 148800,'snus' , 'Golden Dread Requisition', false)")
	db.Exec("INSERT INTO public.sell_was(id, price, market, name, status) VALUES ('228', 148800,'snus' , 'Golden Dread Requisition', false)")
}