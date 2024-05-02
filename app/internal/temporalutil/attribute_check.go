package temporalutil

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/operatorservice/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

// Checks for the existence of the specified Custom Search Attribute. If
// missing, it will attempt to create it. The current implementation only
// supports self-hosted deployments; doing this in Temporal Cloud involves
// a significantly different approach, which we've skipped for now. If the
// TEMPORAL_ADDRESS environment variable includes '.tmprl.cloud:', this
// function assumes that the application is using Temporal Cloud. In that
// case, it reminds the user (via a message to STDERR) that they must
// manually ensure the CSA exists.
func EnsureSearchAttributeExists(ctx context.Context, client client.Client, namespaceName string, temporalHostPort string, attr temporal.SearchAttributeKey) error {
	if isTemporalCloud(temporalHostPort) {
		log.Printf("Reminder: You must ensure that the '%s' Custom Search Attribute exists in your Temporal Cloud Namespace", attr.GetName())
		return nil
	}

	if namespaceName == "" {
		// Unlike ClientOptions creation, the AddSearchAttributes call
		// below requires that the Namespace is explicitly specified
		namespaceName = "default"
		fmt.Printf("Namespace name unspecified; using value '%s'\n", namespaceName)
	}

	attribMap := map[string]enums.IndexedValueType{
		attr.GetName(): attr.GetValueType(),
	}

	_, err := client.OperatorService().AddSearchAttributes(ctx,
		&operatorservice.AddSearchAttributesRequest{
			Namespace:        namespaceName,
			SearchAttributes: attribMap,
		})
	var deniedErr *serviceerror.PermissionDenied
	var alreadyErr *serviceerror.AlreadyExists

	if errors.As(err, &alreadyErr) {
		log.Printf("Search Attribute %s already exists", attr.GetName())
	} else if err != nil {
		log.Fatalf("Failed to add Search Attribute %s: %v", attr, err)

		if !errors.As(err, &deniedErr) {
			return err
		}
	} else {
		log.Printf("Search Attribute %s added", attr.GetName())
	}

	return nil
}

// returns true if the application appears to be configured for Temporal
// Cloud, returns false otherwise
func isTemporalCloud(temporalHostPort string) bool {
	return strings.Contains(temporalHostPort, ".tmprl.cloud:")
}
