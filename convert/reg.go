package convert

var (
	globalReg = makeConverterRegistry()
)

func AllConverterNames() []string {
	return globalReg.AllConverterNames()
}

type convReg struct {
	reg map[string]Converter
}

func makeConverterRegistry() *convReg {
	return &convReg{
		reg: map[string]Converter{},
	}
}

func (r *convReg) Register(c Converter) {
	r.reg[c.Name()] = c
}

func (r *convReg) Get(name string) Converter {
	return r.reg[name]
}

func (r *convReg) AllConverterNames() []string {
	var names []string
	for n := range r.reg {
		names = append(names, n)
	}
	return names
}
