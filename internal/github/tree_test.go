package github

import (
	"context"
	"fmt"
	"testing"
)

func TestGetTree(t *testing.T) {
	repository, err := NewRepository(".")
	if err != nil {
		t.Fatal(err)
	}
	objectIter, err := repository.NewBatchObjectIter(context.TODO())
	oid, _ := NewOID("4e21ae6ca34ea8767b87689c739524c8794c39fd")
	objectIter.RequestObject(oid)
	object, _, err := objectIter.Next()
	fmt.Print(object)
}
