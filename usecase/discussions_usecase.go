package usecase

import (
	"EthioGuide/domain"
	"context"
)

type PostUseCase struct {
	PostRepo domain.IPostRepository
}



func NewPostUseCase(dr domain.IPostRepository) domain.IPostUseCase {
	return &PostUseCase{
		PostRepo: dr,
	}
}

func (d *PostUseCase) CreatePost(ctx context.Context, Post *domain.Post) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return d.PostRepo.CreatePost(ctx, Post)
}

// // DeletePost implements domain.IPostUseCase.
// func (d *PostUseCase) DeletePost(ctx context.Context, id int) error {
// 	panic("unimplemented")
// }

// // GetPostByID implements domain.IPostUseCase.
// func (d *PostUseCase) GetPostByID(ctx context.Context, id int) (*domain.Post, error) {
// 	panic("unimplemented")
// }

// // GetPosts implements domain.IPostUseCase.
// func (d *PostUseCase) GetPosts(ctx context.Context) ([]*domain.Post, error) {
// 	panic("unimplemented")
// }

// // UpdatePost implements domain.IPostUseCase.
// func (d *PostUseCase) UpdatePost(ctx context.Context, Post *domain.Post) error {
// 	panic("unimplemented")
// }
