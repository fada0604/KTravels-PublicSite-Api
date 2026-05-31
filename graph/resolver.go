package graph

import "context"

type Resolver struct{}

func (r *mutationResolver) Noop(ctx context.Context) (bool, error) {
	return false, nil
}

func (r *queryResolver) Health(ctx context.Context) (string, error) {
	return "ok", nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() QueryResolver       { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
