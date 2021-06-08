package main

import (
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
	"strconv"
)


func transaction_add() {
	item := item_{}
	i := 0
	for {
		select {
		case item = <-steam_transactions:
			i++
			if (item.name == "") {

			} else if (item._type == "Buy") {
				item.add_buy()
			} else if (item._type == "Sell") {
				item.add_sell()
			}
			if (i == steam_req_volume+dmarket_req_volume){
				quit <- 1
				i = 0
			}
		}
	}
}
func (item *item_)add_buy(){
	rows, err := db.Query("select id from transactions where id = $1", item.id)

	if err != nil {
		errorLog.Println(err)
		return
	}
	defer rows.Close()
	sz := 0
	for rows.Next() {
		time := ""
		err := rows.Scan(&time)
		if err != nil {
			errorLog.Println(err)
			return
		}
		sz++
	}
	if (sz == 0) {
		_, err := db.Exec("insert into transactions (id, name, price, market, par) values ($1, $2, $3, $4, $5)", item.id, item.name, item.price, item.market, item.id)
		if (err != nil) {
			errorLog.Println(err)
			return
		}
	}
}
func (item *item_)add_sell() bool{
	rows, err := db.Query("select status from sell_was where id = $1", item.id)
	if (err != nil){
		errorLog.Println(err)
	}
	sz := 0
	var c bool = true
	for rows.Next() {
		err := rows.Scan(&c)
		if err != nil {
			errorLog.Println(err)
			return true
		}
		sz++
	}

	if (sz == 0 || c == false){
		rows, err := db.Query("select id from transactions where (status = false AND name = $1)", item.name)
		if err != nil {
			errorLog.Println(err)
			return true
		}
		var id string = ""
		for rows.Next() {
			err = rows.Scan(&id)
			if err != nil {
				errorLog.Println(err)
				return true
			}
			break
		}
		if (id == ""){
			errorLog.Println("Нет шмотки с таким именем!: ", item.name)
			if (sz == 0) {
				_, err = db.Exec("insert into sell_was(id, price, market, status, name) values ($1,$2,$3, false, $4)", item.id, item.price, item.market, item.name)
			}
			if (err != nil){
				errorLog.Println(err)
				return true
			}
			return true
		}
		if (sz == 0){
			_, err = db.Exec("insert into sell_was(id, price, market, status, name) values ($1,$2,$3, false, $4)", item.id, item.price, item.market, item.name)
			if (err != nil){
				errorLog.Println(err)
				return true
			}
		}

		_, err = db.Exec("UPDATE transactions SET status = false, sell_price = $1, sell_market = $3 WHERE id = $2",item.price, id, item.market)
		if err != nil {
			errorLog.Println(err)
			spew.Dump(item)
			return true
		}
		_, err  = db.Exec("UPDATE sell_was SET status = TRUE WHERE id = $1", item.id)
		if (err != nil){
			errorLog.Println(err)
			return true
		}
	}
	return true
}
func link_items(){
	rows, err := db.Query("select id ,price , sell_price from transactions where status = false and sell_price != -1")
	if (err != nil){
		errorLog.Println(err)
		return
	}
	was_id := make(map[string]int)
	count := make(map[string]int)
	for rows.Next(){
		a := struct {
			id string
			price float64
			sell_price float64
		}{}
		rows.Scan(&a.id, &a.price, &a.sell_price)
		if (was_id[a.id] == 2){
			continue
		}
		need_rows, err := db.Query("select id, execel_row, par, price, sell_price, count from transactions where (price>=$1 and price<=$2) and (sell_price>=$3 and sell_price<=$4) and par = transactions.id and id != $5", a.price-0.01, a.price+0.01, a.sell_price-0.01, a.sell_price+0.01, a.id)
		if (err != nil){
			errorLog.Println(err)
			return
		}
		was := false
		for (need_rows.Next()){
			b:= struct {
				id string
				row int
				par string
				price float64
				sell_price float64
				count int
			}{}
			err = need_rows.Scan(&b.id, &b.row, &b.par, &b.price, &b.sell_price, &b.count)
			if (was_id[b.id] == 1){
				continue
			}
			was_id[b.id] =2;
			if (count[b.id] == 0){
				count[b.id] = 1;
			}
			count[b.id]++;
			if (err != nil){
				errorLog.Println(err)
				return
			}
			was = true
			_, err := db.Exec("update transactions set execel_row= $1, par = $2, price = $3, sell_price = $4, status = true where id  = $5", b.row, b.par, b.price, b.sell_price, a.id)
			if err != nil{
				errorLog.Println(err)
				return
			}
			_, err = db.Exec("update transactions set status  = false, count = $2 where id = $1", b.par, b.count+1)
			if err != nil{
				errorLog.Println(err)
				return
			}
		}
		if was == false{
			db.Exec("update transaction set status = false, par = $1 where id = $1", a.id)
		} else {
			was_id[a.id] = 1

		}
	}
}
func not_added_sell(){
	rows, err := db.Query("select id, price, market, name from sell_was where status = false")
	if err != nil {
		errorLog.Println(err)
		return
	}

	for rows.Next() {
		item := item_{}
		err := rows.Scan(&item.id, &item.price, &item.market, &item.name)
		if err != nil {
			errorLog.Println(err)
			return
		}
		item.add_sell()
	}
	return
}
func write_to_google(){
	rows, err := db.Query("select name,price, market, execel_row, sell_price, sell_market, count from transactions where status = false")
	if err != nil {
		errorLog.Println(err)
		return
	}
	update_items := []sql_item{}
	add_items := []sql_item{}
	for (rows.Next()){
		i := sql_item{}
		err := rows.Scan(&i.name, &i.price, &i.market, &i.execel_row, &i.sell_price, &i.sell_market, &i.count)
		if err != nil {
			errorLog.Println(err)
			return
		}
		if (i.execel_row == -1){
			add_items = append(add_items, i);
		} else {
			update_items = append(update_items, i);
		}
	}
	data := []*sheets.ValueRange{}
	for _, v := range update_items{
		val := [][]interface{}{[]interface{} {
			v.name, v.count, v.market, v.price, v.sell_market, v.sell_price}}
		i := &sheets.ValueRange{
			Range: "A"+strconv.Itoa(v.execel_row),
			Values: val,
		}
		data = append(data, i)
	}
	row := get_last_row()
	row++
	for _, v := range add_items{
		val := [][]interface{}{[]interface{} {
			v.name, v.count, v.market, v.price, v.sell_market, v.sell_price,
		}}
		i := &sheets.ValueRange{
			Range: "A"+strconv.Itoa(row),
			Values: val,
		}
		row++;
		data = append(data, i)
	}
	rb := &sheets.BatchUpdateValuesRequest{
		Data:                         data,
		ValueInputOption:             "RAW",
	}
	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Do()
	if err != nil {
		errorLog.Println(err)
		return
	}
	/*data := []*sheets.ValueRange{&sheets.ValueRange{
		Range:  "A1",
		Values: [][]interface{}{[]interface{}{"zxc"}},
	}, &sheets.ValueRange{
		Range: "C1",
		Values: [][]interface{}{[]interface{}{"asd"}},
	}}
	rb := &sheets.BatchUpdateValuesRequest{
		Data:                         data,
		ValueInputOption:             "RAW",
	}
	_, err := srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Do()
	if err != nil {
		errorLog.Println(err)
	}*/

}
func get_last_row()int{
	ctx := context.Background()

	range2 := "A:Z"

	valueRenderOption := "FORMATTED_VALUE"

	dateTimeRenderOption := "SERIAL_NUMBER"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, range2).ValueRenderOption(valueRenderOption).DateTimeRenderOption(dateTimeRenderOption).Context(ctx).Do()
	if err != nil {
		errorLog.Println(err)
	}

	return len(resp.Values)

}



func update_1(){
	key := []string{"Name1", "Name2", "Name4"}
	val := []string{"1", "2", "4"}
	for i, _ := range key{
		_, err := db.Exec("UPDATE table SET val = $1 WHERE name = $2",val[i], key[i])
		if err != nil {
			errorLog.Println(err)
			return
		}
	}
}











//select id, execel_row, par, price, sell_price, count from transactions where (price>=0.07 and price<=0.08) and (sell_price>=148799.99 and sell_price<=148800.01) and par = transactions.id