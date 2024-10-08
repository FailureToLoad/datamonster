package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/failuretoload/datamonster/ent"
)

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return r.client.Noder(ctx, id)
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return r.client.Noders(ctx, ids)
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// CreateSettlementInput returns CreateSettlementInputResolver implementation.
func (r *Resolver) CreateSettlementInput() CreateSettlementInputResolver {
	return &createSettlementInputResolver{r}
}

// UpdateSettlementInput returns UpdateSettlementInputResolver implementation.
func (r *Resolver) UpdateSettlementInput() UpdateSettlementInputResolver {
	return &updateSettlementInputResolver{r}
}

type queryResolver struct{ *Resolver }
type createSettlementInputResolver struct{ *Resolver }
type updateSettlementInputResolver struct{ *Resolver }
