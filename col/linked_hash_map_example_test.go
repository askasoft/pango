package col

import (
	"encoding/json"
	"fmt"
)

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

	for mi := lm.Back(); mi != nil; mi = mi.Prev() {
		fmt.Printf("%-12s: %v\n", mi.Key(), mi.Value())
	}

	// Output:
	// query       : 35.192.xx.xxx
	// proxy       : false
	// mobile      : true
	// as          : AS15169 Google Inc.
	// org         : Google Cloud
	// isp         : Google Cloud
	// timezone    : America/Los_Angeles
	// lon         : -122.0574
	// lat         : 37.4192
	// zip         : 94043
	// city        : Mountain View
	// regionName  : California
	// region      : CA
	// countryCode : US
	// country     : United States
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
	for mi := lm.Front(); mi != nil; mi = mi.Next() {
		fmt.Printf("%-12s: %v\n", mi.Key(), mi.Value())
		if mi.Key() == "city" {
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