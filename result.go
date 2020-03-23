package simply

type Result struct {
	Success  bool
	Complete bool
	output   string
}

func (r Result) String() string {
	return r.output
}
