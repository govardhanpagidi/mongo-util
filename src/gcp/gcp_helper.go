package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
)

// saveSecret adds a new secret version to the given secret with the
// provided payload.
func saveSecret(config Config, user MongoUser, secretStr string) error {
	secretId := user.ProjectID + "-" + user.Username
	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)

	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	// Create the request to create the secret.
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", config.GCP.ProjectID),
		SecretId: secretId,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	parent := ""
	secret, err := client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		//log.Println("WARNING: %v", err)
		parent = createSecretReq.Parent + "/secrets/" + secretId
	} else {
		parent = secret.Name
	}

	// Declare the payload to store.
	payload := []byte(secretStr)

	// Build the request.
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	version, err := client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		log.Fatalf("failed to add secret version: %v", err)
	}

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: version.Name,
	}

	// Call the API.
	_, err = client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	// Print the secret payload.
	//
	// WARNING: Do not print the secret in a production environment - this
	// snippet is showing how to access the secret material.

	//log.Printf("Read secret froom GCP: %s", result.Payload.Data)
	fmt.Println("")
	return nil
}
