package repository

import (
	"context"
	"time"

	"github.com/ablaze/gonexttemp-backend/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&refreshToken, "token = ? AND expires_at > ?", token, time.Now()).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *tokenRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "token = ?", token).Error
}

func (r *tokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "user_id = ?", userID).Error
}

func (r *tokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Delete(&model.RefreshToken{}, "expires_at < ?", time.Now()).Error
}
