// Package db provides database operations.
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
)

func (p *Pool) Login(ctx context.Context, email, password string) (string, error) {
	// check pool
	if p == nil || p.inner == nil {
		return "", fmt.Errorf("connection cannot be stablished")
	}

	row := p.inner.QueryRow(ctx, "select id, password_hash from users where email = $1", email)
	var id string
	var passwordHash string
	err := row.Scan(&id, &passwordHash)
	if err != nil {
		// as the user does not exist, we are going to use a dummy hash
		// to avoid timing attacks
		argon2id.ComparePasswordAndHash("dummy", "dummy")
		return "", fmt.Errorf("email/password does not match")
	}

	match, err := argon2id.ComparePasswordAndHash(password, passwordHash)
	if err != nil {
		return "", fmt.Errorf("email/password does not match")
	}
	if !match {
		return "", fmt.Errorf("email/password does not match")
	}
	return id, nil
}

func (p *Pool) Signup(ctx context.Context, email, password string) (string, error) {
	// check pool
	if p == nil || p.inner == nil {
		return "", fmt.Errorf("connection cannot be stablished")
	}

	// we try to get the user and if this exists we return the same id
	row := p.inner.QueryRow(ctx, "select id from users where email = $1", email)
	var id string
	err := row.Scan(&id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		// if the user does not exist we create it
		hashPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			return "", fmt.Errorf("signup failed")
		}
		err = p.inner.QueryRow(ctx, "insert into users (email, password_hash) values ($1, $2) returning id", email, hashPassword).Scan(&id)
		if err != nil {
			return "", fmt.Errorf("signup failed")
		}
		return id, nil
	}

	return id, nil
}
