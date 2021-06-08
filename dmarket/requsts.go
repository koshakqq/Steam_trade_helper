package dmarket

import (
	json2 "encoding/json"
	"math"
	"strconv"
)

func (k Keys)Get_history() History{
	resp_ := k.request("", "/exchange/v1/history?offset=0&limit="+strconv.Itoa(Dmarket_req_volume), "GET")
	res := History{}
	json2.Unmarshal(resp_, &res)
	return res
}

func (k Keys)Usd_to_Rub(f float64) float64{
	resp_ := k.request("", "/currency-rate/v1/rates", "GET")
	res := currency{}
	json2.Unmarshal(resp_, &res)
	ans := math.Floor( f*res.Rates.Rub*100)/100
	return ans
}
//"https://api.dmarket.com/exchange/v1/history?offset=0&limit=80