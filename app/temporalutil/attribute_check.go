package temporalutil

import (
	"context"
	"errors"
	"log"
	"strings"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/operatorservice/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
)

// EnsureSearchAttributeExists checks for the existence of the specified
// Custom Search Attribute. If missing, it will attempt to create it if
// using a self-hosted deployment. If using Temporal Cloud, it will emit
// a reminder stating that the user must create the attribute.
func EnsureSearchAttributeExists(ctx context.Context, client client.Client, clientOptions client.Options, attr temporal.SearchAttributeKey) error {
	if IsTemporalCloud(clientOptions.HostPort) {
		log.Printf("Reminder: You must ensure that the '%s' Custom Search Attribute exists in your Temporal Cloud Namespace", attr.GetName())
		return nil
	}

	attribMap := map[string]enums.IndexedValueType{
		attr.GetName(): attr.GetValueType(),
	}

	_, err := client.OperatorService().AddSearchAttributes(ctx,
		&operatorservice.AddSearchAttributesRequest{
			Namespace:        clientOptions.Namespace,
			SearchAttributes: attribMap,
		})
	var alreadyErr *serviceerror.AlreadyExists

	if errors.As(err, &alreadyErr) {
		log.Printf("Required Search Attribute %s is present", attr.GetName())
	} else if err != nil {
		log.Fatalf("Failed to add Search Attribute %s: %v", attr, err)
		return err
	} else {
		log.Printf("Search Attribute %s added", attr.GetName())
	}

	return nil
}

// IsTemporalCloud returns true if the application appears to be
// configured for Temporal Cloud based on the provided address,
// or false otherwise
func IsTemporalCloud(temporalHostPort string) bool {
	return strings.Contains(temporalHostPort, ".tmprl.cloud:")
}
