package monke

type NavItem struct {
	Url  string
	Name string
}

type NavData struct {
	Items []NavItem
}

var Nav NavData

func NavInit() {
	Nav.Items = []NavItem{}

	for _, value := range Db.Categories {
		Nav.Items = append(Nav.Items, NavItem{
			Url:  value.Url,
			Name: value.Title,
		})
	}
}
