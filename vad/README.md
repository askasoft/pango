Package validator
=================

Package vad implements value validations for structs and individual fields based on tags.

It has the following **unique** features:

-   Cross Field and Cross Struct validations by using validation tags or custom validators.
-   Slice, Array and Map diving, which allows any or all levels of a multidimensional field to be validated.
-   Ability to dive into both map keys and values for validation
-   Handles type interface by determining it's underlying type prior to validation.
-   Handles custom field types such as sql driver Valuer see [Valuer](https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29)
-   Alias validation tags, which allows for mapping of several validations to a single tag for easier defining of validations on structs
-   Extraction of custom defined Field Name e.g. can specify to extract the JSON name while validating and have it available in the resulting FieldError

Installation
------------

Use go get.

	go get github.com/askasoft/pango

Then import the validator package into your own code.

	import "github.com/askasoft/pango/vad"


Error Return Value
-------------------

Validation functions return type error

They return type error to avoid the issue discussed in the following, where err is always != nil:

Validator returns only InvalidValidationError for bad validation input, nil or ValidationErrors as type error; so, in your code all you need to do is check if the error returned is not nil, and if it's not check if error is InvalidValidationError ( if necessary, most of the time it isn't ) type cast it to type ValidationErrors like so:

```go
err := validate.Struct(mystruct)
validationErrors := err.(vad.ValidationErrors)
 ```

Usage and documentation
------------------------

##### Examples:

- [Simple](./_examples/simple/main.go)
- [Custom Field Types](./_examples/custom/main.go)
- [Struct Level](./_examples/struct-level/main.go)

Baked-in Validations
---------------------

### Fields:

| Tag | Description |
| - | - |
| eqcsfield | Field Equals Another Field (relative)|
| eqfield | Field Equals Another Field |
| fieldcontains | NOT DOCUMENTED IN doc.go |
| fieldexcludes | NOT DOCUMENTED IN doc.go |
| gtcsfield | Field Greater Than Another Relative Field |
| gtecsfield | Field Greater Than or Equal To Another Relative Field |
| gtefield | Field Greater Than or Equal To Another Field |
| gtfield | Field Greater Than Another Field |
| ltcsfield | Less Than Another Relative Field |
| ltecsfield | Less Than or Equal To Another Relative Field |
| ltefield | Less Than or Equal To Another Field |
| ltfield | Less Than Another Field |
| necsfield | Field Does Not Equal Another Field (relative) |
| nefield | Field Does Not Equal Another Field |

### Network:

| Tag | Description |
| - | - |
| cidr | Classless Inter-Domain Routing CIDR |
| cidrv4 | Classless Inter-Domain Routing CIDRv4 |
| cidrv6 | Classless Inter-Domain Routing CIDRv6 |
| datauri | Data URL |
| fqdn | Full Qualified Domain Name (FQDN) |
| hostname | Hostname RFC 952 |
| hostname_port | HostPort |
| hostname_rfc1123 | Hostname RFC 1123 |
| ip | Internet Protocol Address IP |
| ipv4 | Internet Protocol Address IPv4 |
| ipv6 | Internet Protocol Address IPv6 |
| mac | Media Access Control Address MAC |
| uri | URI String |
| url | URL String |
| httpurl | URL (http://) String |
| httpsurl | URL (https://) String |
| httpxurl | URL (https?://) String |

### Strings:

| Tag | Description |
| - | - |
| ascii | ASCII |
| boolean | Boolean |
| number | ASCII Number |
| numeric | Numeric |
| decimal | Decimal |
| letter | ASCII Letter Only |
| letternum | ASCII Letter or Number |
| utfletter | Unicode Letter |
| utfletternum | Unicode Letter or Number |
| printable | Printable Unicode |
| printascii | Printable ASCII |
| lowercase | Lowercase |
| uppercase | Uppercase |
| multibyte | Multi-Byte Characters |
| contains | Contains |
| containsany | Contains Any |
| endsnotwith | Ends Not With |
| endswith | Ends With |
| excludes | Excludes |
| excludesall | Excludes All |
| startsnotwith | Starts Not With |
| startswith | Starts With |
| rematch | Regular Expression Match |
| wcmatch | Wildcard Expression Match |

### Format:
| Tag | Description |
| - | - |
| base64 | Base64 String |
| base64url | Base64URL String |
| btc_addr | Bitcoin Address |
| btc_addr_bech32 | Bitcoin Bech32 Address (segwit) |
| datetime | Datetime |
| duration | Duration (day support, e.g. "30d") |
| cron | Cron Expression |
| periodic | Periodic Expression |
| e164 | e164 formatted phone number |
| email | E-mail String
| hexadecimal | Hexadecimal String |
| hexcolor | Hexcolor String |
| hsl | HSL String |
| hsla | HSLA String |
| html | HTML Tags |
| html_encoded | HTML Encoded |
| isbn | International Standard Book Number |
| isbn10 | International Standard Book Number 10 |
| isbn13 | International Standard Book Number 13 |
| json | JSON |
| jsonobject | JSON Object |
| jsonarray | JSON Array |
| jwt | JSON Web Token (JWT) |
| latitude | Latitude |
| longitude | Longitude |
| postcode_iso3166_alpha2 | Postcode |
| postcode_iso3166_alpha2_field | Postcode |
| regexp | Regular Expression |
| rgb | RGB String |
| rgba | RGBA String |
| ssn | Social Security Number SSN |
| timezone | Timezone |
| uuid | Universally Unique Identifier UUID |
| uuid3 | Universally Unique Identifier UUID v3 |
| uuid4 | Universally Unique Identifier UUID v4 |
| uuid5 | Universally Unique Identifier UUID v5 |
| semver | Semantic Versioning 2.0.0 |
| swiftcode | Business Identifier Code (ISO 9362) |
| ulid | Universally Unique Lexicographically Sortable Identifier ULID |

### Comparisons:
| Tag | Description |
| - | - |
| eq | Equals |
| ne | Not Equal |
| gt | Greater than|
| gte | Greater than or equal |
| min | Greater than or equal |
| lt | Less Than |
| lte | Less Than or Equal |
| max | Less Than or Equal |
| btw | Between |

### Length:
| Tag | Description |
| - | - |
| len | (string's rune count, slice/map length) Equal |
| maxlen | (string's rune count, slice/map length) Maximum |
| minlen | (string's rune count, slice/map length) Minimum |
| btwlen | (string's rune count, slice/map length) Between |

### Other:
| Tag | Description |
| - | - |
| isempty | Is Default |
| oneof | One Of |
| required | Required |
| required_if | Required If |
| required_unless | Required Unless |
| required_with | Required With |
| required_with_all | Required With All |
| required_without | Required Without |
| required_without_all | Required Without All |
| unique | Unique |

#### Aliases:
| Tag | Description |
| - | - |
| alpha | letter |
| alphanum | letternum |
| alphaunicode | utfletter |
| alphanumunicode | utfletternum |
| iscolor | hexcolor\|rgb\|rgba\|hsl\|hsla |
| bic | swiftcode |

