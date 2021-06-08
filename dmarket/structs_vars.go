package dmarket

type Keys struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

type History struct {
	Objects []history_item `json:"objects"`
}

type history_item struct {
	Type string `json:"type"`
	Name string `json:"subject"`
	Id string `json:"customId"`
	Changes []struct{
		Money struct{
			Amount float64 `json:"amount,string"`
		} `json:"money"`
	} `json:"changes"`
}
type currency struct {
	Rates struct{
		Rub float64 `json:"RUB"`
	}`json:"Rates"`
}

var Dmarket_req_volume = 500;

/*type
{
"private": "f004bf5b4a8d308c229b7c2ea1bbd55a993a1fcee1f8772ba4261068f02c2cc445a511212fcb92114018c5d4e07127573e09801a22f3a10bc7f888878169a7c7",
"public": "45a511212fcb92114018c5d4e07127573e09801a22f3a10bc7f888878169a7c7"
}*/
