package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

const listTemplateText = `
<html>
<body>
<table>
  <tr style='text-align:left'>
    <th>Item</th>
    <th>Price</th>
  </tr>
  {{range $item, $price := .}}
  <tr>
    <td>{{$item}}</td>
    <td>{{$price}}</td>
  </tr>
  {{end}}
</table>
</body>
</html>
`

var listTemplate *template.Template = template.Must(template.New("list").Parse(listTemplateText))

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/", db.get)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	listTemplate.Execute(w, db)
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	item := q.Get("item")
	priceStr := q.Get("price")

	var price float32
	_, err := fmt.Sscanf(priceStr, "%f", &price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid price: %q; %s\n", priceStr, err)
		fmt.Fprintf(w, "invalid price: %q\n", priceStr)
		return
	}

	if _, ok := db[item]; ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid item: %q; item already exists\n", item)
		fmt.Fprintf(w, "invalid item: %q; item already exists\n", item)
		return
	}

	db[item] = dollars(price)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK\n")
}

func (db database) get(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	item := q.Get("item")

	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "item: %q not found\n", item)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s: %s\n", item, price)
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	item := q.Get("item")
	priceStr := q.Get("price")

	var price float32
	_, err := fmt.Sscanf(priceStr, "%f", &price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid price: %q; %s\n", priceStr, err)
		fmt.Fprintf(w, "invalid price: %q\n", priceStr)
		return
	}

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("item: %q not exists\n", item)
		fmt.Fprintf(w, "item: %q not exists\n", item)
		return
	}

	db[item] = dollars(price)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s: %s\n", item, db[item])
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	item := q.Get("item")

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("item: %q not exists\n", item)
		fmt.Fprintf(w, "item: %q not exists\n", item)
		return
	}

	delete(db, item)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK\n")
}
