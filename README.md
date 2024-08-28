# go-mongo-uic
Provides a simple, reusable way to implement optimistic concurrency control in MongoDB for GoLang. This package is inspired by Mongoose's *update-if-current plugin*, allowing you to safely update documents in a concurrent environment.

## Installation
To install the package, use go get:
```
go get github.com/agumbe-ai/go-mongo-uic
```

## Usage
1. Import the package.
   ```
    import (
        "github.com/agumbe-ai/go-mongo-uic/versioned"
    )
   ```

2. Define your Document Struct. 
   
   Your MongoDB document should include a **Version** field to manage versioning:
   ```
    type Workspace struct {
        ID           string    `json:"id" bson:"_id,omitempty"`
        Name         string    `json:"name" bson:"name"`
        Description  string    `json:"description" bson:"description"`
        Version      int       `json:"version" bson:"version"`
    }

   ```
3. Set Initial Version when creating a Document.
   
   Use the **SetInitialVersion** function to ensure that the Version field is initialized properly when creating a new document:
   ```
    workspace := Workspace{
        Name:        "Example Workspace",
        Description: "This is a workspace example",
    }
    versioned.SetInitialVersion(&workspace.Version)

   ```
4. Insert the Document into MongoDB.
   ```
    collection := db.Collection("workspaces")
    result, err := collection.InsertOne(context.TODO(), workspace)
    if err != nil {
        log.Fatalf("Failed to insert workspace: %v", err)
    }

    workspace.ID = result.InsertedID.(primitive.ObjectID).Hex()
   ```
5. Update the Document with Version Control.
   
   When updating a document, ensure that the Version field matches the current version in the database. Use the **UpdateIfCurrent** function to perform an atomic update:
   ```
    filter := bson.M{
        "_id":     objectID,
        "version": workspace.Version, // Ensure the current version matches
    }

    update := bson.M{
        "$set": bson.M{
            "name":        "Updated Workspace Name",
            "description": "Updated description",
            "updated_at":  time.Now(),
        },
    }

    result, err := versioned.UpdateIfCurrent(context.TODO(), collection, filter, update, workspace.Version)
    if err != nil {
        if err == versioned.ErrVersionConflict {
            log.Println("Version conflict detected")
        } else {
            log.Fatalf("Failed to update workspace: %v", err)
        }
    } else {
        err = result.Decode(&workspace)
        if err != nil {
            log.Fatalf("Failed to decode updated workspace: %v", err)
        }
    }

   ```

6. Handle Version Conflicts.
   
   If the version in the database does not match the version provided in the update, the update will not be applied, and the function will return an **ErrVersionConflict**. This allows you to handle concurrency issues safely.
   ```
    if err == versioned.ErrVersionConflict {
        log.Println("Version conflict detected")
        // Implement your conflict resolution logic here
    }

   ```

## API Reference
`SetInitialVersion(version *int)`

Initializes the version field to 1 if it's currently 0. Use this function when creating new documents.

`UpdateIfCurrent(ctx context.Context, collection *mongo.Collection, filter bson.M, update bson.M, version int) (*mongo.SingleResult, error)`

Attempts to update a document if the current version matches the provided version. If successful, it increments the version field atomically.

`ErrVersionConflict`
An error returned when the version in the database does not match the expected version, indicating that the document has been modified by another process.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue to discuss changes.

## License
This project is licensed under the Apache 2.0 License.