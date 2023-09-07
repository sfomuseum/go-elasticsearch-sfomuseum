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
	"github.com/sfomuseum/go-sfomuseum-instagram/media"
	"github.com/sfomuseum/go-whosonfirst-elasticsearch/index"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
	"strings"
	"time"
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

		im_rsp := gjson.GetBytes(body, "millsfield:images")

		if im_rsp.Exists() {

			count := len(im_rsp.Array())

			body, err = sjson.SetBytes(body, "millsfield:count_images", count)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign millsfield:count_images, %w", err)
			}
		}

		texts_rsp := gjson.GetBytes(body, "millsfield:images_texts")

		if texts_rsp.Exists() {

			texts_array := texts_rsp.Array()
			texts_count := len(texts_array)

			if texts_count > 0 {

				texts := make([]string, texts_count)

				for idx, t := range texts_array {
					texts[idx] = t.String()
				}

				body, err = sjson.SetBytes(body, "millsfield:images_texts", texts)

				if err != nil {
					return nil, fmt.Errorf("Failed to assign millsfield:images_texts, %w", err)
				}
			}
		}

		// Instagram stuff
		// tl;dr is "convert IG's goofy datetime strings in RFC3339 so that Elasticsearch isn't sad"
		// See also: sfomuseum/go-sfomuseum-instagram and sfomuseum/go-sfomuseum-instagram-publish

		ig_rsp := gjson.GetBytes(body, "instagram:post")

		if ig_rsp.Exists() {

			taken_rsp := gjson.GetBytes(body, "instagram:post.taken_at")

			t, err := time.Parse(media.TIME_FORMAT, taken_rsp.String())

			if err != nil {
				return nil, fmt.Errorf("Failed to parse '%s', %w", taken_rsp.String(), err)
			}

			body, err = sjson.SetBytes(body, "instagram:post.taken_at", t.Format(time.RFC3339))

			if err != nil {
				return nil, err
			}

			tags_rsp := gjson.GetBytes(body, "instagram:post.caption.hashtags")

			if tags_rsp.Exists() {

				hashtags := make([]string, 0)

				for _, t := range tags_rsp.Array() {
					hashtags = append(hashtags, strings.ToLower(t.String()))
				}

				body, err = sjson.SetBytes(body, "instagram:post.caption.hashtags", hashtags)

				if err != nil {
					return nil, fmt.Errorf("Failed to update IG hash tags, %w", err)
				}

			}

		}

		return body, nil
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
