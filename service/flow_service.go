package service

import (
	"github.com/google/uuid"
	"errors"
)

// Struktur response flow generik
// Komentar dalam Bahasa Indonesia

// Enum state flow
const (
	FlowStarted   = "started"
	FlowChallenge = "challenge"
	FlowSuccess   = "success"
	FlowFailed    = "failed"
)

// Struktur error flow
type FlowError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type FlowResponse struct {
	FlowID     string      `json:"flow_id"`
	State      string      `json:"state"`         // started, challenge, success, failed
	Error      *FlowError  `json:"error,omitempty"`
	NextAction string      `json:"next_action"`    // input_code, input_password, done
	Data       interface{} `json:"data,omitempty"`
}

// RegisterFlowService modular
func RegisterFlow(userService *UserService, email, password string) FlowResponse {
	flowID := uuid.NewString()
	user, err := userService.Register(email, password)
	if err != nil {
		return FlowResponse{
			FlowID: flowID,
			State: FlowFailed,
			Error: &FlowError{Code: "internal_error", Message: err.Error()},
			NextAction: "input_email_password",
		}
	}
	return FlowResponse{
		FlowID: flowID,
		State: FlowSuccess,
		NextAction: "request_verification",
		Data: user,
	}
}

// RecoveryFlowService modular
func RecoveryFlow(userService *UserService, email string) FlowResponse {
	flowID := uuid.NewString()
	code, err := userService.GenerateAndSaveRecoveryCode(email)
	if err != nil {
		return FlowResponse{
			FlowID: flowID,
			State: FlowFailed,
			Error: &FlowError{Code: "ERR_RECOVERY", Message: err.Error()},
			NextAction: "input_email",
		}
	}
	return FlowResponse{
		FlowID: flowID,
		State: FlowChallenge,
		NextAction: "input_code_new_password",
		Data: map[string]interface{}{
			"info": "Kode recovery dikirim ke email (simulasi)",
			"code": code, // hanya untuk simulasi/dev
		},
	}
}

// VerificationFlowService modular
func VerificationFlow(userService *UserService, email, code string) FlowResponse {
	flowID := uuid.NewString()
	// TODO: Implementasikan verifikasi kode email sesuai flow modular. Jika sudah ada di UserService, gunakan. Jika tidak, sesuaikan flow ini agar tidak error.
	ok := false
	err := errors.New("verifikasi email belum diimplementasikan")
	if err != nil {
		return FlowResponse{
			FlowID: flowID,
			State: FlowFailed,
			Error: &FlowError{Code: "ERR_VERIFICATION", Message: err.Error()},
			NextAction: "input_code",
		}
	}
	if !ok {
		return FlowResponse{
			FlowID: flowID,
			State: FlowFailed,
			Error: &FlowError{Code: "ERR_VERIFICATION_CODE", Message: "Kode verifikasi salah atau kadaluarsa"},
			NextAction: "input_code",
		}
	}
	return FlowResponse{
		FlowID: flowID,
		State: FlowSuccess,
		NextAction: "done",
		Data: map[string]interface{}{"message": "Email berhasil diverifikasi"},
	}
}

// LoginFlowService modular
func LoginFlow(userService *UserService, jwtService *JWTService, email, password string) FlowResponse {
	flowID := uuid.NewString()
	user, err := userService.Login(email, password)
	if err != nil {
		return FlowResponse{
			FlowID: flowID,
			State: FlowFailed,
			Error: &FlowError{Code: "ERR_LOGIN", Message: err.Error()},
			NextAction: "input_email_password",
		}
	}
	token, err := jwtService.GenerateToken(user.ID)
	if err != nil {
		return FlowResponse{
			FlowID: flowID,
			State: FlowFailed,
			Error: &FlowError{Code: "ERR_TOKEN", Message: "Gagal generate token"},
			NextAction: "input_email_password",
		}
	}
	return FlowResponse{
		FlowID: flowID,
		State: FlowSuccess,
		NextAction: "done",
		Data: map[string]interface{}{
			"token": token,
			"user": user,
		},
	}
}
