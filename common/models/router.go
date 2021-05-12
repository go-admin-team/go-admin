package models

type Routers struct {
	List []Router
}

type Router struct {
	HttpMethod, RelativePath string
}

func (e *Routers) Add(HttpMethod, RelativePath string) *Routers {
	e.List = append(e.List, Router{HttpMethod: HttpMethod,RelativePath: RelativePath})
	return e
}