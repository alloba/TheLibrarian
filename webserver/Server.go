package webserver

import "net/http"

type requestFunction func (http.ResponseWriter, *http.Request)

type serverManager struct {
    preHandleChain []requestFunction
}

func registerPrehandle(manager serverManager, hook func (http.ResponseWriter, *http.Request)){
    manager.preHandleChain = append(manager.preHandleChain, hook) 
}


func New() manager {
   var manager := serverManager {make([]requestFunction, 0)}
   return manager  
}
