## Overview

xmlrpc is an implementation of client side part of XML-RPC protocol in Go language.  
based on https://github.com/kolo/xmlrpc of Dmitry Maksimov (dmtmax@gmail.com).


## Usage

```golang
func main() {
	client := &xmlrpc.Client{
		Endpoint: "https://bugzilla.mozilla.org/xmlrpc.cgi",
	}

	result := struct{
		Version string `xmlrpc:"version"`
	}{}

	_ = client.Call("Bugzilla.version", nil, &result) // 20250903.1

	fmt.Printf("Version: %s\n", result.Version)
}
```

### Arguments encoding

xmlrpc package supports encoding of native Go data types to method
arguments.

Data types encoding rules:

* int, int8, int16, int32, int64 encoded to int;
* float32, float64 encoded to double;
* bool encoded to boolean;
* string encoded to string;
* time.Time encoded to datetime.iso8601;
* []byte encoded to base64;
* slice encoded to array;

Structs encoded to struct by following rules:

* all public field become struct members;
* field name become member name;
* if field has xmlrpc tag, its value become member name.
* for fields tagged with `",omitempty"`, empty values are omitted;
* fields tagged with `"-"` are omitted.

Server method can accept few arguments, to handle this case there is
special approach to handle slice of empty interfaces (`[]any`).
Each value of such slice encoded as separate argument.

### Result decoding

Result of remote function is decoded to native Go data type.

Data types decoding rules:

* int, i4 decoded to int, int8, int16, int32, int64;
* double decoded to float32, float64;
* boolean decoded to bool;
* string decoded to string;
* array decoded to slice;
* structs decoded following the rules described in previous section;
* datetime.iso8601 decoded as time.Time data type;
* base64 decoded to []byte.

