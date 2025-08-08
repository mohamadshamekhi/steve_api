package fetcherservice

import (
	"context"
	"sync"

	mongodbarticle "github.com/abolfazlnorzad/steve/repository/mongo/article"
	mongodbbackoffice "github.com/abolfazlnorzad/steve/repository/mongo/backoffice"
	mongodbproduct "github.com/abolfazlnorzad/steve/repository/mongo/product"
	mongodbuser "github.com/abolfazlnorzad/steve/repository/mongo/user"
	commentservice "github.com/abolfazlnorzad/steve/service/comment"
)

type CommentableFetcherImpl struct {
	productRepo mongodbproduct.DB
	articleRepo mongodbarticle.DB
}

func NewCommentableFetcher(productRepo mongodbproduct.DB, articleRepo mongodbarticle.DB) *CommentableFetcherImpl {
	return &CommentableFetcherImpl{productRepo: productRepo, articleRepo: articleRepo}
}

func (f *CommentableFetcherImpl) FindCommentablesByIDs(ctx context.Context, idsByType map[string][]string) (map[string]commentservice.Commentable, error) {
	result := make(map[string]commentservice.Commentable)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	if ids, ok := idsByType["product"]; ok && len(ids) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			products, err := f.productRepo.FindManyByIDs(ctx, ids)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			for _, p := range products {
				result[p.ID] = &p
			}
			mu.Unlock()
		}()
	}

	if ids, ok := idsByType["article"]; ok && len(ids) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			articles, err := f.articleRepo.FindManyByIDs(ctx, ids)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			for _, a := range articles {
				result[a.ID] = &a
			}
			mu.Unlock()
		}()
	}

	wg.Wait()
	close(errChan)

	for fetchErr := range errChan {
		if fetchErr != nil {
			return nil, fetchErr
		}
	}

	return result, nil
}

type CommenterFetcherImpl struct {
	userRepo       mongodbuser.DB
	backofficeRepo mongodbbackoffice.DB
}

func NewCommenterFetcher(userRepo mongodbuser.DB, backofficeRepo mongodbbackoffice.DB) *CommenterFetcherImpl {
	return &CommenterFetcherImpl{userRepo: userRepo, backofficeRepo: backofficeRepo}
}
func (f *CommenterFetcherImpl) FindCommentersByIDs(ctx context.Context, idsByType map[string][]string) (map[string]commentservice.Commenter, error) {
	result := make(map[string]commentservice.Commenter)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	if ids, ok := idsByType["user"]; ok && len(ids) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			users, err := f.userRepo.FindManyByIDs(ctx, ids)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			for _, u := range users {
				result[u.ID] = &u
			}
			mu.Unlock()
		}()
	}

	if ids, ok := idsByType["backoffice"]; ok && len(ids) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			boUsers, err := f.backofficeRepo.FindManyByIDs(ctx, ids)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			for _, bo := range boUsers {
				result[bo.ID] = &bo
			}
			mu.Unlock()
		}()
	}

	wg.Wait()
	close(errChan)

	for fetchErr := range errChan {
		if fetchErr != nil {
			return nil, fetchErr
		}
	}

	return result, nil
}
