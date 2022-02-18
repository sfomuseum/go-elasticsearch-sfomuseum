# go-elasticsearch-sfomuseum

Go package for indexing SFO Museum Who's On First style records in Elasticsearch.

## Tools

## es-sfomuseum-index

es-sfomuseum-index bulk indexes one or more whosonfirst/go-whosonfirst-iterate/v2 sources in an Elasticsearch database.

```
$> ./bin/es-sfomuseum-index -h
  -append-spelunker-v1-properties
    	Append and index auto-generated Whos On First Spelunker properties.
  -elasticsearch-endpoint string
    	A fully-qualified Elasticsearch endpoint. (default "http://localhost:9200")
  -elasticsearch-index string
    	A valid Elasticsearch index. (default "millsfield")
  -index-alt-files
    	Index alternate geometries.
  -index-only-properties
    	Only index GeoJSON Feature properties (not geometries).
  -index-spelunker-v1
    	Index GeoJSON Feature properties inclusive of auto-generated Whos On First Spelunker properties.
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterator/emitter URI. Supported emitter URI schemes are: directory://,featurecollection://,file://,filelist://,geojsonl://,git://,null://,repo:// (default "repo://")
  -workers int
    	The number of concurrent workers to index data using. Default is the value of runtime.NumCPU().
```

## See also

* https://github.com/sfomuseum/go-elasticsearch-whosonfirst