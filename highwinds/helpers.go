package highwinds

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// Errors and string checks
const (
	ErrBadImportParse = "unexpected format of import ID (%s), expected account_hash/ID"
)

func getContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, 8*time.Second)
}

// ResourceImportParseHashID Highwinds requires an account_hash in addition to
// the resource ID, so on imports we must input account_hash/ID unless
// someone knows a way to get the Importer func to get it from the resource
// definition
// https://www.terraform.io/docs/extend/resources/import.html
func ResourceImportParseHashID(input string) (string, string, error) {
	parts := strings.SplitN(input, "/", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf(ErrBadImportParse, input)
	}

	return parts[0], parts[1], nil
}

func devLog(message string, opts ...interface{}) {
	msg := fmt.Sprintf("===== [DEV] %s", message)
	log.Printf(msg, opts...)
}
