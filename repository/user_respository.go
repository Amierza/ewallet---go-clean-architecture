package repository

import (
	"context"
	"math"
	"strings"

	"github.com/Amierza/e-wallet/dto"
	"github.com/Amierza/e-wallet/entity"
	"gorm.io/gorm"
)

type (
	UserRepository interface {
		RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
		CheckPhoneNumber(ctx context.Context, tx *gorm.DB, phoneNumber string) (entity.User, bool, error)
		GetAllUsersWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error)
		FindUserByID(ctx context.Context, tx *gorm.DB, userID string) (entity.User, error)
		CheckTargetUser(ctx context.Context, tx *gorm.DB, userID string) (entity.User, bool, error)
		UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) error
		CreateTopUp(ctx context.Context, tx *gorm.DB, topup entity.TopUp) error
		CreatePayment(ctx context.Context, tx *gorm.DB, payment entity.Payment) error
		CreateTransfer(ctx context.Context, tx *gorm.DB, transfer entity.Transfer) error
		GetAllTransactionWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllTransactionRepositoryResponse, error)
	}

	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) RegisterUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) CheckPhoneNumber(ctx context.Context, tx *gorm.DB, phoneNumber string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("phone_number = ?", phoneNumber).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

func (r *userRepository) FindUserByID(ctx context.Context, tx *gorm.DB, userID string) (entity.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("id = ?", userID).Take(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Save(&user).Error
}

func (r *userRepository) CreateTopUp(ctx context.Context, tx *gorm.DB, topup entity.TopUp) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&topup).Error
}

func (r *userRepository) CreatePayment(ctx context.Context, tx *gorm.DB, payment entity.Payment) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&payment).Error
}

func (r *userRepository) CreateTransfer(ctx context.Context, tx *gorm.DB, transfer entity.Transfer) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Create(&transfer).Error
}

func (r *userRepository) CheckTargetUser(ctx context.Context, tx *gorm.DB, userID string) (entity.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user entity.User
	if err := tx.WithContext(ctx).Where("id = ?", userID).Take(&user).Error; err != nil {
		return entity.User{}, false, err
	}

	return user, true, nil
}

func (r *userRepository) GetAllUsersWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var users []entity.User
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.User{})

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(phone_number) LIKE ? OR LOWER(address) LIKE ?",
			searchValue, searchValue, searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	if err := query.Order("created_at DESC").Scopes(Paginate(req.Page, req.PerPage)).Find(&users).Error; err != nil {
		return dto.GetAllUserRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.GetAllUserRepositoryResponse{
		Users: users,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *userRepository) GetAllTransactionWithPagination(ctx context.Context, tx *gorm.DB, req dto.PaginationRequest) (dto.GetAllTransactionRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var TopUps []entity.TopUp
	var Payments []entity.Payment
	var Transfers []entity.Transfer
	var err error
	var countTopup, countPayment, countTransfer int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	topupQuery := tx.WithContext(ctx).Model(&entity.TopUp{})
	if err := topupQuery.Count(&countTopup).Error; err != nil {
		return dto.GetAllTransactionRepositoryResponse{}, err
	}

	if err := topupQuery.Order("created_at DESC").Scopes(Paginate(req.Page, req.PerPage)).Find(&TopUps).Error; err != nil {
		return dto.GetAllTransactionRepositoryResponse{}, err
	}

	paymentQuery := tx.WithContext(ctx).Model(&entity.Payment{})
	if err := paymentQuery.Count(&countPayment).Error; err != nil {
		return dto.GetAllTransactionRepositoryResponse{}, err
	}

	if err := paymentQuery.Order("created_at DESC").Scopes(Paginate(req.Page, req.PerPage)).Find(&Payments).Error; err != nil {
		return dto.GetAllTransactionRepositoryResponse{}, err
	}

	transferQuery := tx.WithContext(ctx).Model(&entity.Transfer{})
	if err := transferQuery.Count(&countTransfer).Error; err != nil {
		return dto.GetAllTransactionRepositoryResponse{}, err
	}

	if err := transferQuery.Order("created_at DESC").Scopes(Paginate(req.Page, req.PerPage)).Find(&Transfers).Error; err != nil {
		return dto.GetAllTransactionRepositoryResponse{}, err
	}

	totalCount := countTopup + countPayment + countTransfer
	totalPage := int64(math.Ceil(float64(totalCount) / float64(req.PerPage)))

	return dto.GetAllTransactionRepositoryResponse{
		TopUps:    TopUps,
		Payments:  Payments,
		Transfers: Transfers,
		PaginationResponse: dto.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   totalCount,
			MaxPage: totalPage,
		},
	}, err
}
