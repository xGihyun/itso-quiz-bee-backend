package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type session struct {
	SessionID string    `json:"sessionId"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (r *repository) generateSessionToken() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	token := encoder.EncodeToString(bytes)

	return token, nil
}

func (r *repository) createSession(ctx context.Context, token, userID string) (session, error) {
	hash := sha256.Sum256([]byte(token))
	sessionID := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	ses := session{
		SessionID: sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	byt, err := json.Marshal(ses)
	if err != nil {
		return session{}, err
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	if err := r.redisClient.Set(ctx, sessionKey, string(byt), time.Until(expiresAt)).Err(); err != nil {
		return session{}, err
	}

	userSessionsKey := fmt.Sprintf("user_sessions:%s", userID)
	if err := r.redisClient.SAdd(ctx, userSessionsKey, sessionID).Err(); err != nil {
		return session{}, err
	}

	return ses, nil
}

type sessionValidationResponse struct {
	User    UserResponse `json:"user"`
	Session session      `json:"session"`
}

func (r *repository) validateSessionToken(
	ctx context.Context,
	token string,
) (sessionValidationResponse, error) {
	hash := sha256.Sum256([]byte(token))
	sessionID := hex.EncodeToString(hash[:])

	sessionKey := fmt.Sprintf("session:%s", sessionID)

	data, err := r.redisClient.Get(ctx, sessionKey).Result()
	if err != nil {
		return sessionValidationResponse{}, err
	}

	var ses session

	if err := json.Unmarshal([]byte(data), &ses); err != nil {
		return sessionValidationResponse{}, err
	}

	now := time.Now()
	if now.After(ses.ExpiresAt) || now.Equal(ses.ExpiresAt) {
		if err := r.redisClient.Del(ctx, sessionKey).Err(); err != nil {
			return sessionValidationResponse{}, err
		}

		userSessionsKey := fmt.Sprintf("user_sessions:%s", ses.UserID)
		if err := r.redisClient.SRem(ctx, userSessionsKey, sessionID).Err(); err != nil {
			return sessionValidationResponse{}, err
		}
	}

	// If session is close to expiration (3 days), extend it
	beforeExpiry := ses.ExpiresAt.Add(-3 * 24 * time.Hour)
	if now.After(beforeExpiry) || now.Equal(beforeExpiry) {
		ses.ExpiresAt = now.Add(7 * 24 * time.Hour)

		if err := r.redisClient.Set(
			ctx,
			sessionKey,
			ses,
			time.Until(ses.ExpiresAt),
		).Err(); err != nil {
			return sessionValidationResponse{}, err
		}
	}

	user, err := r.GetByID(ctx, ses.UserID)
	if err != nil {
		return sessionValidationResponse{}, err
	}

	res := sessionValidationResponse{
		Session: ses,
		User:    user,
	}

	return res, nil
}

func (r *repository) invalidateSession(ctx context.Context, token, userID string) error {
	hash := sha256.Sum256([]byte(token))
	sessionID := hex.EncodeToString(hash[:])

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	if err := r.redisClient.Del(ctx, sessionKey).Err(); err != nil {
		return err
	}

	userSessionsKey := fmt.Sprintf("user_sessions:%s", userID)
	if err := r.redisClient.SRem(ctx, userSessionsKey, sessionID).Err(); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(result), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
