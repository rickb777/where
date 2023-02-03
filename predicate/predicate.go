package predicate

const (
	IsNull               = " IS NULL"
	IsNotNull            = " IS NOT NULL"
	EqualTo              = "=?"
	NotEqualTo           = "<>?"
	GreaterThan          = ">?"
	GreaterThanOrEqualTo = ">=?"
	LessThan             = "<?"
	LessThanOrEqualTo    = "<=?"
	Between              = " BETWEEN ? AND ?"
	Like                 = " LIKE ?"
)
