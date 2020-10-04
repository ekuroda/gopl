package main

import (
	"fmt"
	"gopl/ch7/ex7.16/eval"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const indexTemplateText = `
<html>
<body>
<h1>Expression</h1>
<form action="/expr" method="post">
<input type="text" name="expr">
<input type="submit" value="submit">
</form>
</body>
</html>
`

const exprTemplateText = `
<html>
<body>
<h1>Expression</h1>
<div>
{{ .Expr }}
</div>
<h1>Env</h1>
<form action="/env" method="post">
<table>
  {{range $var := .Vars }}
  <tr>
    <td>{{$var}}</td>
    <td><input type="text" name="{{$var}}"></td>
  </tr>
  {{end}}
</table>
<input type="submit" value="submit">
</form>
</body>
</html>
`

const envTemplateText = `
<html>
<body>
<h1>Expression</h1>
<div>
{{ .Expr }}
</div>
<h1>Env</h1>
<table>
  {{range $var, $val := .Env }}
  <tr>
    <td>{{$var}}</td>
    <td>{{$val}}</td>
  </tr>
  {{end}}
</table>
<h1>Eval</h1>
<div>
{{ .Eval }}
</div>
</body>
</html>
`

var indexTemplate *template.Template = template.Must(template.New("index").Parse(indexTemplateText))
var exprTemplate *template.Template = template.Must(template.New("expr").Parse(exprTemplateText))
var envTemplate *template.Template = template.Must(template.New("env").Parse(envTemplateText))

var exprCache eval.Expr

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/expr", expr)
	http.HandleFunc("/env", env)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	indexTemplate.Execute(w, nil)
}

func expr(w http.ResponseWriter, req *http.Request) {
	e := req.FormValue("expr")
	var err error
	exprCache, err = eval.Parse(e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid expr: %q; %s", e, err)
		fmt.Fprintf(w, "invalid expr: %q\n", e)
		return
	}

	check := make(map[eval.Var]bool)
	if err := exprCache.Check(check); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("invalid expr: %q; %s", e, err)
		fmt.Fprintf(w, "invalid expr: %q\n", e)
		return
	}

	vars := exprCache.CollectVars(make([]eval.Var, 0))

	templateVarMap := struct {
		Expr string
		Vars []eval.Var
	}{
		Expr: exprCache.String(),
		Vars: vars,
	}

	w.WriteHeader(http.StatusOK)
	exprTemplate.Execute(w, templateVarMap)
}

func env(w http.ResponseWriter, req *http.Request) {
	vars := exprCache.CollectVars(make([]eval.Var, 0))
	env := make(map[eval.Var]float64)

	for _, vv := range vars {
		varQuery := req.FormValue(string(vv))
		val, err := strconv.ParseFloat(varQuery, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("failed to parse var: %s = %s; %s", vv, varQuery, err)
			fmt.Fprintf(w, "failed to parse var: %s = %s\n", vv, varQuery)
			return
		}
		env[eval.Var(vv)] = val
	}

	r := exprCache.Eval(env)

	templateVarMap := struct {
		Expr string
		Env  map[eval.Var]float64
		Eval float64
	}{
		Expr: exprCache.String(),
		Env:  env,
		Eval: r,
	}

	w.WriteHeader(http.StatusOK)
	envTemplate.Execute(w, templateVarMap)
}
