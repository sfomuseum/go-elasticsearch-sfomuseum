// es-sfomuseum-index bulk indexes one or more whosonfirst/go-whosonfirst-iterate/v2 sources in an Elasticsearch database.
package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git/v2"
)

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-whosonfirst-elasticsearch/index"
	"log"
       "github.com/tidwall/gjson"
       "github.com/tidwall/sjson"
)

func main() {

	ctx := context.Background()

	fs, err := index.NewBulkIndexerFlagSet(ctx)

	if err != nil {
		log.Fatalf("Failed to create new flagset, %v", err)
	}

	flagset.Parse(fs)

	opts, err := index.RunBulkIndexerOptionsFromFlagSet(ctx, fs)

	if err != nil {
		log.Fatalf("Failed to create options, %v", err)
	}

	// START OF sfom stuff - this should eventually be moved in to
	// a github.com/sfomuseum/go-elasticsearch-sfomuseum/document
	// package

	sfom_f := func(ctx context.Context, body []byte) ([]byte, error) {

		var err error
		
		im_rsp := gjson.GetBytes(body, "millsfield:images")

		if im_rsp.Exists(){
			
			count :=  len(im_rsp.Array())
			
			body, err = sjson.SetBytes(body, "millsfield:count_images", count)
			
			if err != nil {
				return nil, fmt.Errorf("Failed to assign millsfield:count_images, %w", err)
			}
		}
		
		 return body, err
	}

	opts.PrepareFuncs = append(opts.PrepareFuncs, sfom_f)

	// END OF sfom stuff

	stats, err := index.RunBulkIndexer(ctx, opts)

	if err != nil {
		log.Fatalf("Failed to run bulk tool, %v", err)
	}

	enc_stats, err := json.Marshal(stats)

	if err != nil {
		log.Fatalf("Failed to marshal stats, %v", err)
	}

	fmt.Println(string(enc_stats))
}
