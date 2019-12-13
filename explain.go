package orm

// Explain ..
func (o *Orm) Explain() string {
	switch o.op {
	case OpDelete:
		return o.delete.explain
	case OpInsert:
		return o.insert.explain
	case OpQuery:
		return o.query.explain
	case OpUpdate:
		return o.update.explain
	}
	return ""
}
