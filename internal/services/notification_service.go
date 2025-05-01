package services

import (
	"context"

	"github.com/google/uuid"
)

// NotificationService handles user notifications
type NotificationService interface {
	NotifyVotingPeriodStart(ctx context.Context, categoryID uuid.UUID) error
	NotifyVotingPeriodEnd(ctx context.Context, categoryID uuid.UUID) error
	SendNewNomineeNotification(ctx context.Context, categoryID uuid.UUID, nomineeID uuid.UUID) error
	SendVoteConfirmation(ctx context.Context, userID uuid.UUID, voteID uuid.UUID) error
	SendVotingReminders(ctx context.Context) error
	AnnounceResults(ctx context.Context, categoryID uuid.UUID) error
}

type notificationService struct {
	// Dependencies
}