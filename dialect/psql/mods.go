package psql

import (
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
	"github.com/stephenafamo/bob/query"
)

type withMod[Q interface {
	AppendWith(expr.CTE)
	SetRecursive(bool)
}] struct{}

func (withMod[Q]) With(name string, columns ...string) cteChain[Q] {
	return cteChain[Q](func() expr.CTE {
		return expr.CTE{
			Name:    name,
			Columns: columns,
		}
	})
}

func (withMod[Q]) Recursive(r bool) query.Mod[Q] {
	return mods.Recursive[Q](r)
}

type fromItemMod struct{}

func (fromItemMod) Only() query.Mod[*expr.FromItem] {
	return mods.QueryModFunc[*expr.FromItem](func(q *expr.FromItem) {
		q.Only = true
	})
}

func (fromItemMod) Lateral() query.Mod[*expr.FromItem] {
	return mods.QueryModFunc[*expr.FromItem](func(q *expr.FromItem) {
		q.Lateral = true
	})
}

func (fromItemMod) WithOrdinality() query.Mod[*expr.FromItem] {
	return mods.QueryModFunc[*expr.FromItem](func(q *expr.FromItem) {
		q.WithOrdinality = true
	})
}

type joinChain[Q interface{ AppendJoin(expr.Join) }] func() expr.Join

func (j joinChain[Q]) Apply(q Q) {
	q.AppendJoin(j())
}

func (j joinChain[Q]) As(alias string) joinChain[Q] {
	jo := j()
	jo.Alias = alias

	return joinChain[Q](func() expr.Join {
		return jo
	})
}

func (j joinChain[Q]) Natural() query.Mod[Q] {
	jo := j()
	jo.Natural = true

	return mods.Join[Q](jo)
}

func (j joinChain[Q]) On(on ...any) query.Mod[Q] {
	jo := j()
	jo.On = append(jo.On, on)

	return mods.Join[Q](jo)
}

func (j joinChain[Q]) OnEQ(a, b any) query.Mod[Q] {
	jo := j()
	jo.On = append(jo.On, expr.EQ(a, b))

	return mods.Join[Q](jo)
}

func (j joinChain[Q]) Using(using ...any) query.Mod[Q] {
	jo := j()
	jo.Using = using

	return mods.Join[Q](jo)
}

type joinMod[Q interface{ AppendJoin(expr.Join) }] struct{}

func (j joinMod[Q]) InnerJoin(e any) joinChain[Q] {
	return joinChain[Q](func() expr.Join {
		return expr.Join{
			Type: expr.InnerJoin,
			To:   e,
		}
	})
}

func (j joinMod[Q]) LeftJoin(e any) joinChain[Q] {
	return joinChain[Q](func() expr.Join {
		return expr.Join{
			Type: expr.LeftJoin,
			To:   e,
		}
	})
}

func (j joinMod[Q]) RightJoin(e any) joinChain[Q] {
	return joinChain[Q](func() expr.Join {
		return expr.Join{
			Type: expr.RightJoin,
			To:   e,
		}
	})
}

func (j joinMod[Q]) FullJoin(e any) joinChain[Q] {
	return joinChain[Q](func() expr.Join {
		return expr.Join{
			Type: expr.FullJoin,
			To:   e,
		}
	})
}

func (j joinMod[Q]) CrossJoin(e any) query.Mod[Q] {
	return mods.Join[Q]{
		Type: expr.CrossJoin,
		To:   e,
	}
}

type orderBy[Q interface{ AppendOrder(expr.OrderDef) }] func() expr.OrderDef

func (s orderBy[Q]) Apply(q Q) {
	q.AppendOrder(s())
}

func (o orderBy[Q]) Asc() orderBy[Q] {
	order := o()
	order.Direction = "ASC"

	return orderBy[Q](func() expr.OrderDef {
		return order
	})
}

func (o orderBy[Q]) Desc() orderBy[Q] {
	order := o()
	order.Direction = "DESC"

	return orderBy[Q](func() expr.OrderDef {
		return order
	})
}

func (o orderBy[Q]) Using(operator string) orderBy[Q] {
	order := o()
	order.Direction = "USING " + operator

	return orderBy[Q](func() expr.OrderDef {
		return order
	})
}

func (o orderBy[Q]) NullsFirst() orderBy[Q] {
	order := o()
	order.Nulls = "FIRST"

	return orderBy[Q](func() expr.OrderDef {
		return order
	})
}

func (o orderBy[Q]) NullsLast() orderBy[Q] {
	order := o()
	order.Nulls = "LAST"

	return orderBy[Q](func() expr.OrderDef {
		return order
	})
}

func (o orderBy[Q]) Collate(collation string) orderBy[Q] {
	order := o()
	order.CollationName = collation

	return orderBy[Q](func() expr.OrderDef {
		return order
	})
}

type cteChain[Q interface{ AppendWith(expr.CTE) }] func() expr.CTE

func (c cteChain[Q]) Apply(q Q) {
	q.AppendWith(c())
}

func (c cteChain[Q]) Name(tableName string, columnNames ...string) cteChain[Q] {
	cte := c()
	cte.Name = tableName
	cte.Columns = columnNames
	return cteChain[Q](func() expr.CTE {
		return cte
	})
}

func (c cteChain[Q]) As(q query.Query) cteChain[Q] {
	cte := c()
	cte.Query = q
	return cteChain[Q](func() expr.CTE {
		return cte
	})
}

func (c cteChain[Q]) NotMaterialized() cteChain[Q] {
	var b = false
	cte := c()
	cte.Materialized = &b
	return cteChain[Q](func() expr.CTE {
		return cte
	})
}

func (c cteChain[Q]) Materialized() cteChain[Q] {
	var b = true
	cte := c()
	cte.Materialized = &b
	return cteChain[Q](func() expr.CTE {
		return cte
	})
}

func (c cteChain[Q]) SearchBreadth(setCol string, searchCols ...string) cteChain[Q] {
	cte := c()
	cte.Search = expr.CTESearch{
		Order:   expr.SearchDepth,
		Columns: searchCols,
		Set:     setCol,
	}
	return cteChain[Q](func() expr.CTE {
		return cte
	})
}

func (c cteChain[Q]) SearchDepth(setCol string, searchCols ...string) cteChain[Q] {
	cte := c()
	cte.Search = expr.CTESearch{
		Order:   expr.SearchDepth,
		Columns: searchCols,
		Set:     setCol,
	}
	return cteChain[Q](func() expr.CTE {
		return cte
	})
}

func (c cteChain[Q]) Cycle(set, using string, cols ...string) cteChain[Q] {
	cte := c()
	cte.Cycle.Set = set
	cte.Cycle.Using = using
	cte.Cycle.Columns = cols
	return cteChain[Q](func() expr.CTE {
		return cte
	})
}

func (c cteChain[Q]) CycleValue(value, defaultVal any) cteChain[Q] {
	cte := c()
	cte.Cycle.SetVal = value
	cte.Cycle.DefaultVal = defaultVal
	return cteChain[Q](func() expr.CTE {
		return cte
	})
}

type lockChain[Q interface{ SetFor(expr.For) }] func() expr.For

func (l lockChain[Q]) Apply(q Q) {
	q.SetFor(l())
}

func (l lockChain[Q]) NoWait() lockChain[Q] {
	lock := l()
	lock.Wait = expr.LockWaitNoWait
	return lockChain[Q](func() expr.For {
		return lock
	})
}

func (l lockChain[Q]) SkipLocked() lockChain[Q] {
	lock := l()
	lock.Wait = expr.LockWaitSkipLocked
	return lockChain[Q](func() expr.For {
		return lock
	})
}
