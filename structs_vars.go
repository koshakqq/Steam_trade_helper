package main

import (
	"database/sql"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"strconv"
	"./dmarket"
)

type item_ struct {
	name string
	_type string
	id string
	price float64
	market string
}
type sql_item struct {
	id string
	name string
	price string
	market string
	sell_price string
	execel_row int
	sell_market string
	count int
}

type cookie_stirng struct {
	f string
	s string
	t string
}

type string_pair struct{
	f string
	s string
}

type shets_json struct {
	Range string `json:"Renge"`
}
var items []item_

var cookies_ []cookie_stirng

var steam_req_volume int = 500
var steam_get_url string = "https://steamcommunity.com/market/myhistory/render/?query=&start=0&count="+strconv.Itoa(steam_req_volume)
var steam_transactions = make(chan item_, steam_req_volume*5)

var connStr string = "user=postgres password=228 dbname=trade sslmode=disable"
var db *sql.DB

var srv *sheets.Service
var spreadsheetId string = "1xTfwMDRYwEVuKBLe1_m7eo-UdvDWadskVifWRuT_gWc"

var quit = make(chan int)
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")


var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

var dmarket_req_volume = dmarket.Dmarket_req_volume
var dm = dmarket_keys()
//https://api.dmarket.com/account/v1/user?X-Api-Key:0123f46b02f429cd493b585c5dd0ef3897f79a1d93a8f190809feb0cedc0232c232bf86a627a343cd34945fddf55ba60ef077db6cc708624f016a34b85370154&X-Sign-Date:1&X-Request-Sign:228