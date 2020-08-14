package concurrent

type values struct {
	prev    *values
	next    *values
	data    string
	version int64
}
