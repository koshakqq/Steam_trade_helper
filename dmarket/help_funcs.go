package dmarket

func getRootUrl() string {
	return "https://api.dmarket.com"
}
func Parse_id(s *string) string{
	res := ""
	was := false
	for i := 0; i<len(*s);i++{
		if (was == true){
			res+=string((*s)[i])
		}
		if ((*s)[i] == ':'){
			was = true
		}
	}
	return res
}
