# ShortId

Sometimes, we couldn't use GUID/UUID as primary key in our application.  In this case, we have to generate a short id to resovle the problem.  The package generates a short id which  is 12 characters.  For instance, `15DE398FEB83`.  The first two characters represent to current year, and characters which between 3 to 10 are random.  The last two characters are machineId.  I run 1M times and there is no duplicate key.

## Usuage

```go
package main

import shortid "github.com/jasonsoft/go-short-id"

func main() {
	opt := Options{
		Number:        14,
		StartWithYear: true,
		EndWithHost:   false,
	}
	id := shortid.Generate(opt)
	println(id) // output: 17GuOMTTeMKZQpxZ
}
```

## Reference:
[.Net Version](https://github.com/jasonsoft/short-id)


 
