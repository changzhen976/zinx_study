package ziface

type IRouter interface {
	PreHandle(requst IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}
