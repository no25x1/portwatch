// Package tagger provides a thread-safe registry for associating arbitrary
// key-value label sets (Tags) with scan targets identified by "host:port"
// strings.
//
// Tags travel alongside scan events and can be used by filters, alerting
// rules, and output formatters to enrich or route events based on metadata
// such as environment, region, or service role.
//
// Basic usage:
//
//	reg := tagger.New()
//	reg.Set("db.internal:5432", tagger.Tags{"env": "prod", "role": "db"})
//
//	if tags, ok := reg.Get("db.internal:5432"); ok {
//		fmt.Println(tags["env"]) // prod
//	}
package tagger
