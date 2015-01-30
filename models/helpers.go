// Copyright G2G Market Inc, 2015

package models

import (
	"fmt"
	"strings"
)

const IMAGE_URL_PREFIX string = "https://s3.amazonaws.com/g2gcdn"

// SetThumbnailURL computes the thumbnail URL to display to the user for
// assistance visually identifying the product.
func ThumbnailURL(sku string, ma_id int) string {
	//replace - with "" in sku, add leading 0
	sku = strings.Replace(sku, "-", "", -1)

	//result will look like https://s3.amazonaws.com/g2gcdn/68/00046000820118_200x200.jpg
	//add leading zero because they're all prefixed with that on S3 at the moment
	return fmt.Sprintf("%s/%d/0%s_200x200.jpg",
		IMAGE_URL_PREFIX, ma_id, sku)
}
