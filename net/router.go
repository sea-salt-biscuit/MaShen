package net

type HandlerFunc func()

// acoount login//logout
type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc // Handler处理器
}

type router struct {
	group []*group
}
