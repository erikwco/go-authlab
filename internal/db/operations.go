// Package db provides database operations.
package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
)

type user struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func (p *Pool) GetID(ctx context.Context, email string) (*user, bool) {
	u := &user{}
	// check pool
	if p == nil || p.inner == nil {
		return u, false
	}

	row := p.inner.QueryRow(ctx, "select id, password_hash from users where email = $1", email)
	err := row.Scan(&u.ID, &u.Password)
	if err != nil {
		return u, false
	}
	return u, true
}

func (p *Pool) Login(ctx context.Context, email, password string) (string, error) {
	// get user from db
	u, ok := p.GetID(ctx, email)
	if !ok {
		return "", fmt.Errorf("invalid_credentials")
	}

	match, err := argon2id.ComparePasswordAndHash(password, u.Password)
	if err != nil {
		return "", fmt.Errorf("invalid_credentials")
	}
	if !match {
		return "", fmt.Errorf("invalid_credentials")
	}
	return u.ID, nil
}

func (p *Pool) Signup(ctx context.Context, email, password string) (string, error) {
	_, ok := p.GetID(ctx, email)
	if ok {
		return "", errors.New("invalid_credentials")
	}
	// if the user does not exist we create it
	hashPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("invalid_credentials")
	}
	var generatedID string
	err = p.inner.QueryRow(ctx, "insert into users (email, password_hash) values ($1, $2) returning id", email, hashPassword).Scan(&generatedID)
	if err != nil {
		return "", fmt.Errorf("invalid_credentials")
	}
	return generatedID, nil
}
