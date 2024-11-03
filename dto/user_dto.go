package dto

import (
	"errors"

	"github.com/Amierza/e-wallet/entity"
	"github.com/google/uuid"
)

const (
	// Failed
	MESSAGE_FAILED_GET_DATA_FROM_BODY   = "failed get data from body"
	MESSAGE_FAILED_REGISTER_USER        = "failed create user"
	MESSAGE_FAILED_GET_LIST_USER        = "failed get list user"
	MESSAGE_FAILED_GET_LIST_TRANSACTION = "failed get list transaction"
	MESSAGE_FAILED_LOGIN_USER           = "failed login user"
	MESSAGE_FAILED_PROSES_REQUEST       = "failed proses request"
	MESSAGE_FAILED_TOKEN_NOT_FOUND      = "failed token not found"
	MESSAGE_FAILED_TOKEN_NOT_VALID      = "failed token not valid"
	MESSAGE_FAILED_TOKEN_DENIED_ACCESS  = "denied access"
	MESSAGE_FAILED_TOP_UP               = "failed top up"
	MESSAGE_FAILED_PAYMENT              = "failed payment"
	MESSAGE_FAILED_TRANSFER             = "failed transfer"
	MESSAGE_FAILED_UPDATE_PROFILE_USER  = "failed update profile user"

	// Success
	MESSAGE_SUCCESS_REGISTER_USER        = "success create user"
	MESSAGE_SUCCESS_GET_LIST_USER        = "success get list user"
	MESSAGE_SUCCESS_GET_LIST_TRANSACTION = "success get list transaction"
	MESSAGE_SUCCESS_LOGIN_USER           = "success login user"
	MESSAGE_SUCCESS_TOP_UP               = "success top up"
	MESSAGE_SUCCESS_PAYMENT              = "success payment"
	MESSAGE_SUCCESS_TRANSFER             = "success transfer"
	MESSAGE_SUCCESS_UPDATE_PROFILE_USER  = "success update profile user"
)

var (
	ErrPhoneNumberAlreadyExists   = errors.New("phone number is already exists")
	ErrPhoneNumberNotFound        = errors.New("phone number not found")
	ErrPinNotMatch                = errors.New("pin not match")
	ErrCreateUser                 = errors.New("failed to create user")
	ErrGetUserFromToken           = errors.New("failed to get user from token")
	ErrGetUserFromUserID          = errors.New("failed to get user from user ID")
	ErrUpdateUserBalance          = errors.New("failed to update user balance")
	ErrUpdateUser                 = errors.New("failed to update user")
	ErrInsufficientBalance        = errors.New("failed insufficient balance")
	ErrCreateTopUp                = errors.New("failed to create topup")
	ErrCreatePayment              = errors.New("failed to create payment")
	ErrGetTargetUser              = errors.New("failed to get target user")
	ErrCannotTransferToOwnAccount = errors.New("failed transfer to own account")
	ErrCreateTransfer             = errors.New("failed to create transfer")
)

type (
	UserCreateRequest struct {
		FirstName   string `json:"first_name" form:"first_name"`
		LastName    string `json:"last_name" form:"last_name"`
		PhoneNumber string `json:"phone_number" form:"phone_number"`
		Address     string `json:"address" form:"address"`
		Pin         string `json:"pin" form:"pin"`
	}

	UserResponse struct {
		ID          string `json:"user_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		Pin         string `json:"pin"`

		entity.Timestamp
	}

	AllUserResponse struct {
		ID          string `json:"user_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		Pin         string `json:"pin"`
		Balance     int64  `json:"balance"`

		entity.Timestamp
	}

	UserLoginRequest struct {
		PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required"`
		Pin         string `json:"pin" form:"pin" binding:"required"`
	}

	UserLoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	UserPaginationResponse struct {
		Data []AllUserResponse `json:"data"`
		PaginationResponse
	}

	AllTransactionResponse struct {
		TopUpID       string `json:"top_up_id,omitempty"`
		PaymentID     string `json:"payment_id,omitempty"`
		TransferID    string `json:"transfer_id,omitempty"`
		UserID        string `json:"user_id,omitempty"`
		TargetUserID  string `json:"target_user_id,omitempty"`
		Amount        int64  `json:"amount_top_up,omitempty"`
		Remarks       string `json:"remarks_payment,omitempty"`
		BalanceBefore *int64 `json:"balance_before_top_up,omitempty"`
		BalanceAfter  int64  `json:"balance_after_top_up,omitempty"`
		entity.Timestamp
	}

	TransactionPaginationResponse struct {
		Data []AllTransactionResponse `json:"data"`
		PaginationResponse
	}

	GetAllUserRepositoryResponse struct {
		Users []entity.User
		PaginationResponse
	}

	GetAllTransactionRepositoryResponse struct {
		TopUps    []entity.TopUp
		Payments  []entity.Payment
		Transfers []entity.Transfer
		PaginationResponse
	}

	TopUpRequest struct {
		Amount int64 `json:"amount" binding:"required"`
	}

	TopUpResponse struct {
		ID            string `json:"top_up_id"`
		AmountTopUp   int64  `json:"amount_top_up"`
		BalanceBefore int64  `json:"balance_before"`
		BalanceAfter  int64  `json:"balance_after"`
		entity.Timestamp
	}

	PaymentRequest struct {
		Amount  int64  `json:"amount" binding:"required"`
		Remarks string `json:"remarks"`
	}

	PaymentResponse struct {
		ID            string `json:"payment_id"`
		AmountPayment int64  `json:"amount_payment"`
		Remarks       string `json:"remarks"`
		BalanceBefore int64  `json:"balance_before"`
		BalanceAfter  int64  `json:"balance_after"`
		entity.Timestamp
	}

	TransferRequest struct {
		TargetUser uuid.UUID `json:"target_user" binding:"required"`
		Amount     int64     `json:"amount" binding:"required"`
		Remarks    string    `json:"remarks"`
	}

	TransferResponse struct {
		ID             string `json:"transfer_id"`
		TargetUserID   string `json:"target_user_id"`
		AmountTransfer int64  `json:"amount_transfer"`
		Remarks        string `json:"remarks"`
		BalanceBefore  int64  `json:"balance_before"`
		BalanceAfter   int64  `json:"balance_after"`
		entity.Timestamp
	}

	UpdateProfileRequest struct {
		FirstName   string `json:"first_name,omitempty"`
		LastName    string `json:"last_name,omitempty"`
		PhoneNumber string `json:"phone_number,omitempty"`
		Address     string `json:"address,omitempty"`
		Pin         string `json:"pin,omitempty"`
	}
)
