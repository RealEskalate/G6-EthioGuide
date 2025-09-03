package usecase

import (
	"EthioGuide/domain"
	"context"
	"time"
)

type PostUseCase struct {
	PostRepo domain.IPostRepository
	timeout  time.Duration
}



func NewPostUseCase(dr domain.IPostRepository, contextTimeout time.Duration) domain.IPostUseCase {

	return &PostUseCase{
		PostRepo: dr,
		timeout:  contextTimeout,

	}
}

func (d *PostUseCase) CreatePost(ctx context.Context, Post *domain.Post) (*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()
	return d.PostRepo.CreatePost(ctx, Post)
}

// GetPosts implements domain.IPostUseCase.
func (d *PostUseCase) GetPosts(ctx context.Context, opts domain.PostFilters) ([]*domain.Post,int64, error) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()
	return d.PostRepo.GetPosts(ctx, opts)
}


func (d *PostUseCase) DeletePost(ctx context.Context, id, UserID, role string) error {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()
	return d.PostRepo.DeletePost(ctx, id, UserID, role)
}

func (d *PostUseCase) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()
	return d.PostRepo.GetPostByID(ctx, id)
}


func (d *PostUseCase) UpdatePost(ctx context.Context, Post *domain.Post) (*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()
	return d.PostRepo.UpdatePost(ctx, Post)
}
