package eval

// CollectVars ...
func (v Var) CollectVars(vars []Var) []Var {
	for _, vv := range vars {
		if v == vv {
			return vars
		}
	}

	vars = append(vars, v)
	return vars
}

// CollectVars ...
func (literal) CollectVars(vars []Var) []Var {
	return vars
}

// CollectVars ...
func (u unary) CollectVars(vars []Var) []Var {
	return u.x.CollectVars(vars)
}

// CollectVars ...
func (b binary) CollectVars(vars []Var) []Var {
	vars = b.x.CollectVars(vars)
	return b.y.CollectVars(vars)
}

// CollectVars ...
func (c call) CollectVars(vars []Var) []Var {
	for _, arg := range c.args {
		vars = arg.CollectVars(vars)
	}
	return vars
}
