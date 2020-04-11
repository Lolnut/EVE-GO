package specifications

type Specification interface {
	IsSatisfiedBy(object interface{}) bool
}

type AbstractSpecification struct {
	Specification
}

func (c *AbstractSpecification) And(other Specification) Specification {
	return &andSpecification{c, other}
}

func (c *AbstractSpecification) Or(other Specification) Specification {
	return &orSpecification{c, other}
}

func (c *AbstractSpecification) Not() Specification {
	return &notSpecification{c}
}

type andSpecification struct {
	one Specification
	two Specification
}

func (a *andSpecification) IsSatisfiedBy(object interface{}) bool {
	return a.one.IsSatisfiedBy(object) && a.two.IsSatisfiedBy(object)
}

type orSpecification struct {
	one Specification
	two Specification
}

func (a *orSpecification) IsSatisfiedBy(object interface{}) bool {
	return a.one.IsSatisfiedBy(object) || a.two.IsSatisfiedBy(object)
}

type notSpecification struct {
	one Specification
}

func (a *notSpecification) IsSatisfiedBy(object interface{}) bool {
	return !a.one.IsSatisfiedBy(object)
}

