package service

import (
	"context"
	"sort"
	"sync"

	"github.com/Amierza/e-wallet/dto"
	"github.com/Amierza/e-wallet/entity"
	"github.com/Amierza/e-wallet/helpers"
	"github.com/Amierza/e-wallet/repository"
	"github.com/google/uuid"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
		LoginUser(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
		TopUpUser(ctx context.Context, req dto.TopUpRequest) (dto.TopUpResponse, error)
		PaymentUser(ctx context.Context, req dto.PaymentRequest) (dto.PaymentResponse, error)
		TransferUser(ctx context.Context, req dto.TransferRequest) (dto.TransferResponse, error)
		GetAllTransactionWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.TransactionPaginationResponse, error)
		UpdateProfileUser(ctx context.Context, req dto.UpdateProfileRequest) (dto.UserResponse, error)
	}
	userService struct {
		userRepo   repository.UserRepository
		jwtService JWTService
	}
)

var (
	mu sync.Mutex
)

const (
	LOCAL_URL          = "http://localhost:8080"
	VERIFY_EMAIL_ROUTE = "register/verify_email"
)

func NewUserService(userRepo repository.UserRepository, jwtService JWTService) UserService {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	_, flag, err := s.userRepo.CheckPhoneNumber(ctx, nil, req.PhoneNumber)
	if err == nil || flag {
		return dto.UserResponse{}, dto.ErrPhoneNumberAlreadyExists
	}

	user := entity.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Pin:         req.Pin,
		Balance:     0,
	}

	userReg, err := s.userRepo.RegisterUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.UserResponse{
		ID:          userReg.ID.String(),
		FirstName:   userReg.FirstName,
		LastName:    userReg.LastName,
		PhoneNumber: userReg.PhoneNumber,
		Address:     userReg.Address,
		Pin:         userReg.Pin,
	}, nil
}

func (s *userService) LoginUser(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	user, flag, err := s.userRepo.CheckPhoneNumber(ctx, nil, req.PhoneNumber)
	if err != nil || !flag {
		return dto.UserLoginResponse{}, dto.ErrPhoneNumberNotFound
	}

	checkPin, err := helpers.ChcekPin(user.Pin, []byte(req.Pin))
	if err != nil || !checkPin {
		return dto.UserLoginResponse{}, dto.ErrPinNotMatch
	}

	accessToken, refreshToken, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return dto.UserLoginResponse{}, err
	}

	return dto.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) TopUpUser(ctx context.Context, req dto.TopUpRequest) (dto.TopUpResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	token := ctx.Value("Authorization").(string)

	userID, err := s.jwtService.GetUserIDByToken(token)
	if err != nil {
		return dto.TopUpResponse{}, dto.ErrGetUserFromToken
	}

	user, err := s.userRepo.FindUserByID(ctx, nil, userID)
	if err != nil {
		return dto.TopUpResponse{}, dto.ErrGetUserFromUserID
	}

	balanceBefore := user.Balance
	user.Balance += req.Amount

	if err := s.userRepo.UpdateUser(ctx, nil, user); err != nil {
		return dto.TopUpResponse{}, dto.ErrUpdateUserBalance
	}

	newTopup := entity.TopUp{
		ID:            uuid.New(),
		UserID:        user.ID,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
	}

	if err := s.userRepo.CreateTopUp(ctx, nil, newTopup); err != nil {
		return dto.TopUpResponse{}, dto.ErrCreateTopUp
	}

	return dto.TopUpResponse{
		ID:            newTopup.ID.String(),
		AmountTopUp:   req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
	}, nil
}

func (s *userService) PaymentUser(ctx context.Context, req dto.PaymentRequest) (dto.PaymentResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	token := ctx.Value("Authorization").(string)

	userID, err := s.jwtService.GetUserIDByToken(token)
	if err != nil {
		return dto.PaymentResponse{}, dto.ErrGetUserFromToken
	}

	user, err := s.userRepo.FindUserByID(ctx, nil, userID)
	if err != nil {
		return dto.PaymentResponse{}, dto.ErrGetUserFromUserID
	}

	if user.Balance < req.Amount {
		return dto.PaymentResponse{}, dto.ErrInsufficientBalance
	}

	balanceBefore := user.Balance
	user.Balance -= req.Amount

	if err := s.userRepo.UpdateUser(ctx, nil, user); err != nil {
		return dto.PaymentResponse{}, dto.ErrUpdateUserBalance
	}

	newPayment := entity.Payment{
		ID:            uuid.New(),
		UserID:        user.ID,
		Amount:        req.Amount,
		Remarks:       req.Remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
	}

	if err := s.userRepo.CreatePayment(ctx, nil, newPayment); err != nil {
		return dto.PaymentResponse{}, dto.ErrCreatePayment
	}

	return dto.PaymentResponse{
		ID:            newPayment.ID.String(),
		AmountPayment: req.Amount,
		Remarks:       req.Remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
	}, nil
}

func (s *userService) TransferUser(ctx context.Context, req dto.TransferRequest) (dto.TransferResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	token := ctx.Value("Authorization").(string)

	userID, err := s.jwtService.GetUserIDByToken(token)
	if err != nil {
		return dto.TransferResponse{}, dto.ErrGetUserFromToken
	}

	user, err := s.userRepo.FindUserByID(ctx, nil, userID)
	if err != nil {
		return dto.TransferResponse{}, dto.ErrGetUserFromUserID
	}

	targetUser, flag, err := s.userRepo.CheckTargetUser(ctx, nil, req.TargetUser.String())
	if err != nil || !flag {
		return dto.TransferResponse{}, dto.ErrGetTargetUser
	}

	if targetUser.ID == user.ID {
		return dto.TransferResponse{}, dto.ErrCannotTransferToOwnAccount
	}

	if user.Balance < req.Amount {
		return dto.TransferResponse{}, dto.ErrInsufficientBalance
	}

	balanceBefore := user.Balance

	user.Balance -= req.Amount

	targetUser.Balance += req.Amount

	if err := s.userRepo.UpdateUser(ctx, nil, user); err != nil {
		return dto.TransferResponse{}, dto.ErrUpdateUserBalance
	}

	if err := s.userRepo.UpdateUser(ctx, nil, targetUser); err != nil {
		return dto.TransferResponse{}, dto.ErrUpdateUserBalance
	}

	newTransfer := entity.Transfer{
		ID:            uuid.New(),
		UserID:        user.ID,
		TargetUserID:  req.TargetUser,
		Amount:        req.Amount,
		Remarks:       req.Remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
	}

	if err := s.userRepo.CreateTransfer(ctx, nil, newTransfer); err != nil {
		return dto.TransferResponse{}, dto.ErrCreateTransfer
	}

	return dto.TransferResponse{
		ID:             newTransfer.ID.String(),
		TargetUserID:   req.TargetUser.String(),
		AmountTransfer: req.Amount,
		Remarks:        req.Remarks,
		BalanceBefore:  balanceBefore,
		BalanceAfter:   user.Balance,
	}, nil
}

func (s *userService) GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := s.userRepo.GetAllUsersWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var datas []dto.AllUserResponse
	for _, user := range dataWithPaginate.Users {
		data := dto.AllUserResponse{
			ID:          user.ID.String(),
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			Address:     user.Address,
			Pin:         user.Pin,
			Balance:     user.Balance,
		}

		datas = append(datas, data)
	}

	return dto.UserPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *userService) GetAllTransactionWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.TransactionPaginationResponse, error) {
	dataWithPaginate, err := s.userRepo.GetAllTransactionWithPagination(ctx, nil, req)
	if err != nil {
		return dto.TransactionPaginationResponse{}, err
	}

	var transactions []dto.AllTransactionResponse
	for _, topup := range dataWithPaginate.TopUps {
		transaction := dto.AllTransactionResponse{
			TopUpID:       topup.ID.String(),
			UserID:        topup.UserID.String(),
			Amount:        topup.Amount,
			BalanceBefore: &topup.BalanceBefore,
			BalanceAfter:  topup.BalanceAfter,
			Timestamp: entity.Timestamp{
				CreatedAt: topup.CreatedAt,
				UpdatedAt: topup.UpdatedAt,
				DeletedAt: topup.DeletedAt,
			},
		}

		transactions = append(transactions, transaction)
	}

	for _, payment := range dataWithPaginate.Payments {
		transaction := dto.AllTransactionResponse{
			PaymentID:     payment.ID.String(),
			UserID:        payment.UserID.String(),
			Amount:        payment.Amount,
			Remarks:       payment.Remarks,
			BalanceBefore: &payment.BalanceBefore,
			BalanceAfter:  payment.BalanceAfter,
			Timestamp: entity.Timestamp{
				CreatedAt: payment.CreatedAt,
				UpdatedAt: payment.UpdatedAt,
				DeletedAt: payment.DeletedAt,
			},
		}

		transactions = append(transactions, transaction)
	}

	for _, transfer := range dataWithPaginate.Transfers {
		transaction := dto.AllTransactionResponse{
			TransferID:    transfer.ID.String(),
			UserID:        transfer.UserID.String(),
			TargetUserID:  transfer.TargetUserID.String(),
			Amount:        transfer.Amount,
			Remarks:       transfer.Remarks,
			BalanceBefore: &transfer.BalanceBefore,
			BalanceAfter:  transfer.BalanceAfter,
			Timestamp: entity.Timestamp{
				CreatedAt: transfer.CreatedAt,
				UpdatedAt: transfer.UpdatedAt,
				DeletedAt: transfer.DeletedAt,
			},
		}

		transactions = append(transactions, transaction)
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.Before(transactions[j].CreatedAt)
	})

	return dto.TransactionPaginationResponse{
		Data: transactions,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *userService) UpdateProfileUser(ctx context.Context, req dto.UpdateProfileRequest) (dto.UserResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	token := ctx.Value("Authorization").(string)

	userID, err := s.jwtService.GetUserIDByToken(token)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserFromToken
	}

	user, err := s.userRepo.FindUserByID(ctx, nil, userID)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserFromUserID
	}

	if req.FirstName == "" {
		req.FirstName = user.FirstName
	}

	if req.LastName == "" {
		req.LastName = user.LastName
	}

	if req.PhoneNumber == "" {
		req.PhoneNumber = user.PhoneNumber
	} else {
		_, flag, err := s.userRepo.CheckPhoneNumber(ctx, nil, req.PhoneNumber)
		if err == nil || flag {
			return dto.UserResponse{}, dto.ErrPhoneNumberAlreadyExists
		}
	}

	if req.Address == "" {
		req.Address = user.Address
	}

	if req.Pin == "" {
		req.Pin = user.Pin
	}

	updatedUser := entity.User{
		ID:          user.ID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Pin:         req.Pin,
		Balance:     user.Balance,
	}

	if err := s.userRepo.UpdateUser(ctx, nil, updatedUser); err != nil {
		return dto.UserResponse{}, dto.ErrUpdateUser
	}

	return dto.UserResponse{
		ID:          updatedUser.ID.String(),
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		PhoneNumber: updatedUser.PhoneNumber,
		Address:     updatedUser.Address,
		Pin:         updatedUser.Pin,
	}, nil
}
