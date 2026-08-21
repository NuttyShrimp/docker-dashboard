package model

type Resource struct {
	Name      string
	Ports     []int
	Subdomain *Subdomain `json:"-"`
}
