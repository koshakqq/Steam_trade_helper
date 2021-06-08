package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zellyn/kooky"
	_ "github.com/zellyn/kooky/allbrowsers" // register cookie store finders!
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"./dmarket"
)

func get_from_steam(){
	req, _ := http.NewRequest("GET", steam_get_url, nil)

	for _, v := range cookies_{
		req.AddCookie(&http.Cookie{Name: v.f, Value: v.s, Domain: v.t});
	}
	client := http.Client{}
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body);

	get_items_from_html(string(b))
}
func get_items_from_html(str string) int {
	names := get_real_name(&str)

	items = items[0:0]
	items = make([]item_, steam_req_volume)

	str = strings.Replace(str, "\\", "", -1);
	reader := strings.NewReader(str)
	str = ""
	doc, err := goquery.NewDocumentFromReader(reader)
	if err!=nil{
		fmt.Print(err)
	}

	doc.Find(".market_listing_item_name").Each(func(i int, s *goquery.Selection){
		items[i].name = s.Text();
		id, _ := s.Attr("id")
		items[i].id = parse_steam_id(id)
	})
	doc.Find(".market_listing_price").Each(func(i int, s *goquery.Selection){
		price_str := s.Text()
		r := strings.NewReplacer("\n", "", "r", "", "n", "", "t", "", ".", "", " pуб", "")
		price_str = r.Replace(price_str)
		price_str = strings.Replace(price_str, ",", ".", -1)
		items[i].price, _ = strconv.ParseFloat(price_str, 3);
	})
	doc.Find(".market_listing_gainorloss").Each(func(i int, s *goquery.Selection){
		if (i==0) {
			return
		}
		items[i-1]._type = parse_steam_type(s.Text())
	})
	for i := range items{
		items[i].market =  "Steam"
		if (names[items[i].id] != "") {
			items[i].name = names[items[i].id]
		}
		steam_transactions <- items[i]
	}
	return 1
}
func get_cookies_from_browser(){
	cookies := kooky.ReadCookies(kooky.Valid, kooky.DomainHasSuffix("steamcommunity.com"))
	for _, cookie := range cookies {
		if (cookie.Domain!=".steamcommunity.com") {
			cookies_= append(cookies_,  cookie_stirng{cookie.Name, cookie.Value, cookie.Domain})
		}
	}
}

func get_real_name(str *string)map[string]string{
	var res_ map[string]interface{}
	json.Unmarshal([]byte(*str), &res_)
	assets := res_["assets"].(map[string]interface{})

	items_name := make ([]struct {
		name string
		id string
	}, 0)
	b  := make(map[string]string)
	for _, v := range assets{
		game := v.(map[string]interface{})
		items_name_ := get_name_from_Json(&game).([]struct {
			name string
			id string
		})
		for _, item := range items_name_{
			b[item.id] = item.name
			items_name = append(items_name, struct{
				name string
				id string
			}{ name: item.name, id:item.id})
		}
	}
	var hovers string = res_["hovers"].(string)
	ids := make([]string_pair, 0)
	flag := 0
	cur_id := ""
	for _, v := range hovers{
		if (v == 39){
			if (flag == 0){
					cur_id =""
			}	else if (flag == 1){
					ids = append(ids, string_pair{f:parse_steam_id(cur_id)})
			}	else if (flag == 3){
					cur_id = ""
			}	else if (flag == 5){
					ids[len(ids)- 1].s = cur_id
			}	else if (flag == 11){
					flag = -1
			}
			flag++
		} else {
			if (flag == 1 || flag == 5){
				cur_id+=string(v)
			}
		}

	}
	res := make(map[string]string)
	for _, val := range ids {
		res[val.f] = b[val.s]
	}
	return res
}
func get_name_from_Json(game *map[string]interface{})interface{}{
	res := make ([]struct {
		name string
		id string
	}, 0)
	for key, _ := range (*game) {
		for _, contexid_val := range (*game)[key].(map[string]interface{}) {
			item := contexid_val.(map[string]interface{})
			res = append(res, struct {
				name string
				id   string
			}{name: item["market_name"].(string), id: item["id"].(string)})
		}
	}
	return res
}

func get_dm_items(){
	items := dm.Get_history()
	ob := items.Objects
	kol_vo := dmarket_req_volume
	for _, v := range ob{
		kol_vo--
		if (v.Type == "instant_sell" || v.Type == "sell"){
			v.Type = "Sell"
		} else if (v.Type == "purchase" || v.Type == "target_closed"){
			v.Type = "Buy"
		} else {
			item := item_{
				_type:  "Skip",
				market: "Dmarket",
			}
			steam_transactions <- item
			continue
		}
		item := item_{
			name:  v.Name ,
			_type:  v.Type,
			id: dmarket.Parse_id(&v.Id),
			price:  v.Changes[0].Money.Amount,
			market: "Dmarket",
		}
		steam_transactions <- item
	}
	for kol_vo!=0{
		item := item_{
			_type:  "Skip",
		}
		steam_transactions <- item
		kol_vo--;
	}
	return
}
