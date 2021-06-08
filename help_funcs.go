package main

import (
	"database/sql"
	"log"
	"math/rand"
	"strings"
)

func parse_steam_type(s string)(string){
	r := strings.NewReplacer("\n", "", "r", "", "n", "", "t", "")
	s = r.Replace(s)
	if s==""	{
		return "List"
	}
	if s=="-"{
		return "Sell"
	}
	if s=="+"{
		return "Buy"
	}
	return "Err"
}

func parse_steam_id(s string)(string){
	var s1 string = ""
	var flag bool = false
	for _, v := range s{
		if v >= 48 && v <= 57{
			s1+=string(v)
			flag = true
		} else {
			if (flag){
				break
			}
		}
	}
	return s1
}

func tabbles_conn(){
	var err error
	db, err = sql.Open("postgres", connStr)
	if (err != nil){
		log.Fatalf("DB: ", err)
	}
	google_sheet_conn()
}

func dmarket_keys() dmarket.Keys{
	k := dmarket.Keys{
		Private: "",
		Public:  "",
	}
	return k
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}