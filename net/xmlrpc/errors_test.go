package xmlrpc

import (
	"testing"
)

// const specialFaultRespXml = `
// <?xml version="1.0" encoding="UTF-8"?>
// <methodResponse>
//   <fault>
//     <faultstring>You must log in before using this part of Bugzilla.</faultstring>
//     <faultcode>410</faultcode>
//   </fault>
// </methodResponse>`

// func TestSpecialFailedResponse(t *testing.T) {
// 	resp := Response([]byte(specialFaultRespXml))

// 	rerr := resp.Err()
// 	if rerr == nil {
// 		t.Fatal("Err() error: expected error, got nil")
// 	}

// 	fe, ok := AsFaultError(rerr)
// 	if !ok {
// 		t.Fatalf("Not FaultError: %v", rerr)
// 	}

// 	if fe.FaultCode != 410 && fe.FaultString != "You must log in before using this part of Bugzilla." {
// 		t.Fatal("Err() error: got wrong error")
// 	}
// }

const standardFaultRespXml = `
<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
  <fault>
    <value>
      <struct>
        <member>
          <name>faultString</name>
          <value>
            <string>You must log in before using this part of Bugzilla.</string>
          </value>
        </member>
        <member>
          <name>faultCode</name>
          <value>
            <int>410</int>
          </value>
        </member>
      </struct>
    </value>
  </fault>
</methodResponse>`

func TestStandardFailedResponse(t *testing.T) {
	var v map[string]any

	rerr := DecodeString(standardFaultRespXml, v)
	if rerr == nil {
		t.Fatal("expected error, got nil")
	}

	fe, ok := AsFaultError(rerr)
	if !ok {
		t.Fatalf("Not FaultError: %v", rerr)
	}

	if fe.FaultCode != 410 && fe.FaultString != "You must log in before using this part of Bugzilla." {
		t.Fatalf("got wrong error: %s", fe)
	}
}

const emptyValResp = `
<?xml version="1.0" encoding="UTF-8"?>
<methodResponse>
	<params>
		<param>
			<value>
				<struct>
					<member>
						<name>user</name>
						<value><string>Joe Smith</string></value>
					</member>
					<member>
						<name>token</name>
						<value/>
					</member>
				</struct>
			</value>
		</param>
	</params>
</methodResponse>`

func TestResponseWithEmptyValue(t *testing.T) {
	result := struct {
		User  string `xmlrpc:"user"`
		Token string `xmlrpc:"token"`
	}{}

	if err := DecodeString(emptyValResp, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if result.User != "Joe Smith" || result.Token != "" {
		t.Fatalf("unexpected result: %v", result)
	}
}
