// internal/modules/history/test/service_test.go
package test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	"github.com/kiminodare/HOVARLAY-BE/internal/modules/history"
	mockhistory "github.com/kiminodare/HOVARLAY-BE/internal/modules/history/mock"

	"github.com/stretchr/testify/require"
)

func TestService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	userID := uuid.New()
	text := "Test text for TTS"
	voice := "en-US-Wavenet-D"
	rate := 1.0
	pitch := 0.0
	volume := 0.0

	expectedHistory := &generated.History{
		ID:     uuid.New(),
		Text:   text,
		Voice:  voice,
		Rate:   rate,
		Pitch:  pitch,
		Volume: volume,
	}

	// Setup mock expectation
	mockRepo.EXPECT().
		Create(gomock.Any(), userID, text, voice, rate, pitch, volume).
		Return(expectedHistory, nil)

	// Test
	result, err := service.Create(context.Background(), userID, text, voice, rate, pitch, volume)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedHistory.ID, result.ID)
	require.Equal(t, expectedHistory.Text, result.Text)
	require.Equal(t, expectedHistory.Voice, result.Voice)
}

func TestService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	userID := uuid.New()
	text := "Test text"
	voice := "en-US-Wavenet-D"
	rate := 1.0
	pitch := 0.0
	volume := 0.0

	// Setup mock expectation untuk error
	mockRepo.EXPECT().
		Create(gomock.Any(), userID, text, voice, rate, pitch, volume).
		Return(nil, errors.New("database error"))

	// Test
	result, err := service.Create(context.Background(), userID, text, voice, rate, pitch, volume)

	// Assertions
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "database error")
}

func TestService_GetByUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	userID := uuid.New()
	offset := 0
	limit := 10

	expectedHistories := []*generated.History{
		{
			ID:   uuid.New(),
			Text: "First history",
		},
		{
			ID:   uuid.New(),
			Text: "Second history",
		},
	}

	// Setup mock expectation
	mockRepo.EXPECT().
		GetByUser(gomock.Any(), userID, offset, limit).
		Return(expectedHistories, nil)

	// Test
	result, err := service.GetByUser(context.Background(), userID, offset, limit)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result, 2)
	require.Equal(t, expectedHistories[0].Text, result[0].Text)
	require.Equal(t, expectedHistories[1].Text, result[1].Text)
}

func TestService_GetByUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	userID := uuid.New()
	offset := 0
	limit := 10

	// Setup mock expectation untuk error
	mockRepo.EXPECT().
		GetByUser(gomock.Any(), userID, offset, limit).
		Return(nil, errors.New("database error"))

	// Test
	result, err := service.GetByUser(context.Background(), userID, offset, limit)

	// Assertions
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "database error")
}

func TestService_CountByUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	userID := uuid.New()
	expectedCount := 25

	// Setup mock expectation
	mockRepo.EXPECT().
		CountByUser(gomock.Any(), userID).
		Return(expectedCount, nil)

	// Test
	result, err := service.CountByUser(context.Background(), userID)

	// Assertions
	require.NoError(t, err)
	require.Equal(t, expectedCount, result)
}

func TestService_CountByUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	userID := uuid.New()

	// Setup mock expectation untuk error
	mockRepo.EXPECT().
		CountByUser(gomock.Any(), userID).
		Return(0, errors.New("database error"))

	// Test
	result, err := service.CountByUser(context.Background(), userID)

	// Assertions
	require.Error(t, err)
	require.Equal(t, 0, result)
	require.Contains(t, err.Error(), "database error")
}

func TestService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	historyID := uuid.New()
	expectedHistory := &generated.History{
		ID:   historyID,
		Text: "Test history",
	}

	// Setup mock expectation
	mockRepo.EXPECT().
		GetByID(gomock.Any(), historyID).
		Return(expectedHistory, nil)

	// Test
	result, err := service.GetByID(context.Background(), historyID)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, expectedHistory.ID, result.ID)
	require.Equal(t, expectedHistory.Text, result.Text)
}

func TestService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	historyID := uuid.New()

	// Setup mock expectation untuk not found
	mockRepo.EXPECT().
		GetByID(gomock.Any(), historyID).
		Return(nil, errors.New("history not found"))

	// Test
	result, err := service.GetByID(context.Background(), historyID)

	// Assertions
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "not found")
}

func TestService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	historyID := uuid.New()
	text := "Updated text"
	voice := "en-US-Wavenet-A"
	rate := 1.5
	pitch := 0.5
	volume := -1.0

	// Setup mock expectation
	mockRepo.EXPECT().
		Update(gomock.Any(), historyID, text, voice, rate, pitch, volume).
		Return(nil)

	// Test
	err := service.Update(context.Background(), historyID, text, voice, rate, pitch, volume)

	// Assertions
	require.NoError(t, err)
}

func TestService_Update_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	historyID := uuid.New()
	text := "Updated text"
	voice := "en-US-Wavenet-A"
	rate := 1.5
	pitch := 0.5
	volume := -1.0

	// Setup mock expectation untuk error
	mockRepo.EXPECT().
		Update(gomock.Any(), historyID, text, voice, rate, pitch, volume).
		Return(errors.New("update failed"))

	// Test
	err := service.Update(context.Background(), historyID, text, voice, rate, pitch, volume)

	// Assertions
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed")
}

func TestService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	historyID := uuid.New()

	// Setup mock expectation
	mockRepo.EXPECT().
		Delete(gomock.Any(), historyID).
		Return(nil)

	// Test
	err := service.Delete(context.Background(), historyID)

	// Assertions
	require.NoError(t, err)
}

func TestService_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockhistory.NewMockRepositoryInterface(ctrl)
	service := history.NewService(mockRepo)

	// Test data
	historyID := uuid.New()

	// Setup mock expectation untuk error
	mockRepo.EXPECT().
		Delete(gomock.Any(), historyID).
		Return(errors.New("delete failed"))

	// Test
	err := service.Delete(context.Background(), historyID)

	// Assertions
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed")
}
