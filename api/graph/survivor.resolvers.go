package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/failuretoload/datamonster/ent"
)

// CreateSurvivor is the resolver for the createSurvivor field.
func (r *mutationResolver) CreateSurvivor(ctx context.Context, input ent.CreateSurvivorInput) (*ent.Survivor, error) {
	return ent.FromContext(ctx).Survivor.Create().SetInput(input).Save(ctx)
}

// UpdateSurvivor is the resolver for the updateSurvivor field.
func (r *mutationResolver) UpdateSurvivor(ctx context.Context, id int, input ent.UpdateSurvivorInput) (*ent.Survivor, error) {
	return r.client.Survivor.UpdateOneID(id).SetInput(input).Save(ctx)
}

// DeleteSurvivor is the resolver for the deleteSurvivor field.
func (r *mutationResolver) DeleteSurvivor(ctx context.Context, id int) (*bool, error) {
	return nil, r.client.Survivor.DeleteOneID(id).Exec(ctx)
}

// Survivors is the resolver for the survivors field.
func (r *queryResolver) Survivors(ctx context.Context, filter *ent.SurvivorWhereInput, order *ent.SurvivorOrder) ([]*ent.Survivor, error) {
	query := r.client.Survivor.Query()
	var err error
	if filter != nil {
		query, err = filter.Filter(query)
		if err != nil {
			return nil, err
		}
	}

	if order != nil {
		orderFunc := survivorOrderFunc(order)
		query = query.Order(orderFunc(order.Direction.OrderTermOption()))
	}

	return query.All(ctx)
}
