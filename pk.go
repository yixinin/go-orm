package orm

//Pk .
//usage: Pk() or or Pk(v) or Pk(k,v)
func (o *Orm) Pk(args ...interface{}) *Orm {
	o.usePk = true
	switch len(args) {
	case 0:
		o.parsePk()
	case 1:
		o.parsePk(args[0])
	default:
		pk, ok := args[0].(string)
		if ok {
			o.pk = pk
		}
		o.pkValue = args[1]
	}
	return o
}
