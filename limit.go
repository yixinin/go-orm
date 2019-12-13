package orm

//Limit ..
func (o *Orm) Limit(limit ...int64) *Orm {
	if o.op == OpInsert {
		return o
	}
	if l := len(limit); l >= 2 {
		o.limit = "limit ?,?"
		o.limitArgs = []interface{}{
			limit[0],
			limit[1],
		}
	} else if l == 1 {
		o.limit = "limit ?"
		o.limitArgs = []interface{}{
			limit[0],
		}
	}
	return o
}
