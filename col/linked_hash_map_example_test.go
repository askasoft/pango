package col

import (
	"encoding/json"
	"fmt"
)

func ExampleLinkedHashMap() {
	m := NewLinkedHashMap() // empty (keys are of type int)
	m.Set(2, "b")           // 2->b
	m.Set(1, "x")           // 2->b, 1->x (insertion-order)
	m.Set(1, "a")           // 2->b, 1->a (insertion-order)
	_, _ = m.Get(2)         // b, true
	_, _ = m.Get(3)         // nil, false
	_ = m.Values()          // []interface {}{"b", "a"} (insertion-order)
	_ = m.Keys()            // []interface {}{2, 1} (insertion-order)
	m.Delete(1)             // 2->b
	m.Clear()               // empty
	m.IsEmpty()             // true
	m.Len()                 // 0
}

func ExampleNewLinkedHashMap() {
	// initialize from a list of key-value pairs
	lm := NewLinkedHashMap(
		"country", "United States",
		"countryCode", "US",
		"region", "CA",
		"regionName", "California",
		"city", "Mountain View",
		"zip", "94043",
		"lat", 37.4192,
		"lon", -122.0574,
		"timezone", "America/Los_Angeles",
		"isp", "Google Cloud",
		"org", "Google Cloud",
		"as", "AS15169 Google Inc.",
		"mobile", true,
		"proxy", false,
		"query", "35.192.xx.xxx",
	)

	for it := lm.Iterator(); it.Next(); {
		fmt.Printf("%-12s: %v\n", it.Key(), it.Value())
	}

	// Output:
	// country     : United States
	// countryCode : US
	// region      : CA
	// regionName  : California
	// city        : Mountain View
	// zip         : 94043
	// lat         : 37.4192
	// lon         : -122.0574
	// timezone    : America/Los_Angeles
	// isp         : Google Cloud
	// org         : Google Cloud
	// as          : AS15169 Google Inc.
	// mobile      : true
	// proxy       : false
	// query       : 35.192.xx.xxx
}

func ExampleLinkedHashMap_UnmarshalJSON() {
	const jsonStream = `{
  "country"     : "United States",
  "countryCode" : "US",
  "region"      : "CA",
  "regionName"  : "California",
  "city"        : "Mountain View",
  "zip"         : "94043",
  "lat"         : 37.4192,
  "lon"         : -122.0574,
  "timezone"    : "America/Los_Angeles",
  "isp"         : "Google Cloud",
  "org"         : "Google Cloud",
  "as"          : "AS15169 Google Inc.",
  "mobile"      : true,
  "proxy"       : false,
  "query"       : "35.192.xx.xxx"
}`

	// compare with if using a regular generic map, the unmarshalled result
	//  is a map with unpredictable order of keys
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonStream), &m)
	if err != nil {
		fmt.Println("error:", err)
	}
	for key := range m {
		// fmt.Printf("%-12s: %v\n", key, m[key])
		_ = key
	}

	// use the LinkedHashMap to Unmarshal from JSON object
	lm := NewLinkedHashMap()
	err = json.Unmarshal([]byte(jsonStream), lm)
	if err != nil {
		fmt.Println("error:", err)
	}

	// loop over all key-value pairs,
	// it is ok to call Set append-modify new key-value pairs,
	// but not safe to call Delete during iteration.
	for it := lm.Iterator(); it.Next(); {
		fmt.Printf("%-12s: %v\n", it.Key(), it.Value())
		if it.Key() == "city" {
			lm.Set("mobile", false)
			lm.Set("extra", 42)
		}
	}

	// Output:
	// country     : United States
	// countryCode : US
	// region      : CA
	// regionName  : California
	// city        : Mountain View
	// zip         : 94043
	// lat         : 37.4192
	// lon         : -122.0574
	// timezone    : America/Los_Angeles
	// isp         : Google Cloud
	// org         : Google Cloud
	// as          : AS15169 Google Inc.
	// mobile      : false
	// proxy       : false
	// query       : 35.192.xx.xxx
	// extra       : 42
}
