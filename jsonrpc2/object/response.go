package object

var ()

type Response interface {
	json_object
}

type response_object struct {
	// A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
	jsonrpc string

	// This member is REQUIRED on success.
	// This member MUST NOT exist if there was an error invoking the method.
	// The value of this member is determined by the method invoked on the Server.
	result interface{}

	// This member is REQUIRED on error.
	// This member MUST NOT exist if there was no error triggered during invocation.
	// The value for this member MUST be an Object as defined in section 5.1.
	err Err

	// This member is REQUIRED.
	// It MUST be the same as the value of the id member in the Request Object.
	// If there was an error in detecting the id in the Request object (e.g. Parse error/Invalid Request), it MUST be Null.
	id interface{}
}

func (resp *response_object) JsonObject() JsonObject {
	obj := JsonObject{
		"jsonrpc": resp.jsonrpc,
		"id":      resp.id,
	}
	if resp.err != nil {
		obj["error"] = resp.err.JsonObject()
	} else {
		obj["result"] = resp.result
	}
	return obj
}

func CreateResponse(obj map[string]interface{}) (Response, error) {
	var resp = &response_object{}
	// TODO: fill request object
	return resp, nil
}

func NewResponse(result interface{}, err Err, id interface{}) (Response, error) {
	var resp = &response_object{
		jsonrpc: "2.0",
		result:  result,
		err:     err,
		id:      id,
	}
	return resp, nil
}
