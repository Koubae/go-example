// Created with ChatGPT
package db

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDuplicateCredential         = errors.New("duplicate credential")
	ErrDuplicateProviderForAccount = errors.New("provider already linked to account")
	ErrAccountNotFound             = errors.New("account not found")
	ErrProviderNotFound            = errors.New("provider not found")
	ErrInvalidCredentialState      = errors.New("invalid credential state")
	ErrPersistenceInvariant        = errors.New("persistence invariant violation")
	ErrRetryableTransaction        = errors.New("retryable transaction error")
)

func classifyPostgresError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505": // unique_violation
		switch pgErr.ConstraintName {
		case "account_credentials_provider_id_credential_key":
			return fmt.Errorf("%w: %w", ErrDuplicateCredential, err)
		case "account_credentials_account_id_provider_id_key":
			return fmt.Errorf("%w: %w", ErrDuplicateProviderForAccount, err)
		default:
			return fmt.Errorf("unique violation (%s): %w", pgErr.ConstraintName, err)
		}

	case "23503": // foreign_key_violation
		switch pgErr.ConstraintName {
		case "account_credentials_account_id_fkey":
			return fmt.Errorf("%w: %w", ErrAccountNotFound, err)
		case "account_credentials_provider_id_fkey":
			return fmt.Errorf("%w: %w", ErrProviderNotFound, err)
		default:
			return fmt.Errorf("foreign key violation (%s): %w", pgErr.ConstraintName, err)
		}

	case "23514": // check_violation
		return fmt.Errorf("%w: %w", ErrInvalidCredentialState, err)

	case "23502": // not_null_violation
		return fmt.Errorf("%w: %w", ErrPersistenceInvariant, err)

	case "40001", "40P01": // serialization_failure, deadlock_detected
		return fmt.Errorf("%w: %w", ErrRetryableTransaction, err)

	default:
		return err
	}
}
