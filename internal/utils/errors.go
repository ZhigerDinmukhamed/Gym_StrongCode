package utils

import "errors"

// Common errors
var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidInput       = errors.New("invalid input")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// User errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
)

// Membership errors
var (
	ErrNoActiveMembership = errors.New("no active membership found")
	ErrMembershipExpired  = errors.New("membership has expired")
)

// Booking errors
var (
	ErrClassFull           = errors.New("class is full")
	ErrBookingNotFound     = errors.New("booking not found")
	ErrAlreadyBooked       = errors.New("already booked for this class")
	ErrCannotCancelBooking = errors.New("cannot cancel booking")
)

// Class errors
var (
	ErrClassNotFound = errors.New("class not found")
)

// Payment errors
var (
	ErrInvalidAmount        = errors.New("invalid payment amount")
	ErrPaymentFailed        = errors.New("payment processing failed")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)
