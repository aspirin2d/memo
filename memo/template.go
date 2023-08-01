package memo

import "text/template"

type Templates struct {
	// interference template
	Interference string `toml:"interference"`
	interference *template.Template

	OrderQuery string `toml:"order_query"`
	orderQuery *template.Template
}

func (ts *Templates) Parse() error {
	var err error
	ts.interference, err = template.New("interference").Parse(ts.Interference)
	if err != nil {
		return err
	}
	ts.orderQuery, err = ts.interference.New("order_query").Parse(ts.OrderQuery)
	if err != nil {
		return err
	}
	return nil
}
