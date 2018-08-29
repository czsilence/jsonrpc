package server

import (
	"reflect"

	"github.com/czsilence/jsonrpc/jsonrpc2/object"
)

// rpc method interface
type RPCMethod interface {
	Invoke() (result interface{}, err object.Err)
	InvokeA(params []interface{}) (result interface{}, err object.Err)
}

// rpc method data
type _RPCMethod struct {
	// register name
	n string
	// func object
	rf reflect.Value
	// type of func object
	rft reflect.Type
	// param type list
	rfpt []reflect.Type
}

func (m *_RPCMethod) Invoke() (result interface{}, err object.Err) {
	if m.rft.NumIn() != 0 {
		err = object.ErrMethod_ParamsNumNotMatch
	} else {
		var result_vals = m.rf.Call(nil)
		result = m.return_values(result_vals)
	}
	return
}
func (m *_RPCMethod) InvokeA(params []interface{}) (result interface{}, err object.Err) {
	if m.rft.NumIn() != len(params) {
		err = object.ErrMethod_ParamsNumNotMatch
	} else {
		var param_vals = make([]reflect.Value, len(params))
		for i, p := range params {
			if reflect.TypeOf(p) != m.rfpt[i] {
				err = object.ErrMethod_WrongParamsType
				return
			}
			param_vals[i] = reflect.ValueOf(p)
		}
		var result_vals = m.rf.Call(param_vals)
		result = m.return_values(result_vals)
	}
	return
}

func (m *_RPCMethod) return_values(vals []reflect.Value) (result interface{}) {
	if result_num := m.rf.Type().NumOut(); result_num == 0 {
		result = nil
	} else if result_num == 1 {
		result = vals[0].Interface()
	} else {
		var results = make([]interface{}, result_num)
		for i, rv := range vals {
			results[i] = rv.Interface()
		}
		result = results
	}
	return
}

////////////////////////////////////////////////
type RPCMethodMapper interface {
	RegisterMethod(name string, method interface{}) (err error)
	Get(name string) (method RPCMethod, exist bool)
}

type _RPCMethodMapper struct {
	method_map map[string]RPCMethod
}

// register method to rpc server.
// method MUST be a func. (except variadic function)
// params of method with number type (int/int32/...) MUST use float64
func (mapper *_RPCMethodMapper) RegisterMethod(name string, method interface{}) (err error) {
	if mapper.method_map == nil {
		mapper.method_map = make(map[string]RPCMethod)
	}

	rf := reflect.ValueOf(method)
	if !rf.IsValid() || rf.IsNil() || rf.Kind() != reflect.Func {
		err = Error_Server_InvalidRPCMethod
	} else {
		mapper.map_rpc_method(name, rf)
	}
	return
}

func (mapper *_RPCMethodMapper) map_rpc_method(name string, reflect_func reflect.Value) (err error) {
	if _, ex := mapper.method_map[name]; ex {
		err = Error_Server_DuplicatedRPCMethod
	} else {
		var method = &_RPCMethod{
			n:   name,
			rf:  reflect_func,
			rft: reflect_func.Type(),
		}

		method.rfpt = make([]reflect.Type, method.rft.NumIn())
		for i := 0; i < method.rft.NumIn(); i++ {
			method.rfpt[i] = method.rft.In(i)
		}

		mapper.method_map[name] = method
	}
	return
}

func (mapper *_RPCMethodMapper) Get(name string) (method RPCMethod, exist bool) {
	method, exist = mapper.method_map[name]
	return
}

func NewRPCMethodMapper() RPCMethodMapper {
	return &_RPCMethodMapper{}
}
