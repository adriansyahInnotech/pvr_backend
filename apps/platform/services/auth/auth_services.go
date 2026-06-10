package auth

import (
	"apps/platform/dtos"
	"context"
	"errors"
	"os"

	"pvr_backend/helper"
	responseDto "pvr_backend/helper/response/dto"
	"pvr_backend/middleware"
	jwtdto "pvr_backend/middleware/jwt/dto"
	"pvr_backend/models"
	"pvr_backend/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(traceCtx context.Context, c *fiber.Ctx, dto *dtos.Login) *responseDto.Response
	LoginDevice(traceCtx context.Context, c *fiber.Ctx, dto *dtos.LoginDevice) *responseDto.Response
	Register(traceCtx context.Context, c *fiber.Ctx, dto *dtos.Register) *responseDto.Response
}

type authService struct {
	authRepository     *repository.AuthRepository
	platformRepository *repository.PlatformRepository
	helper             *helper.Helper
	middleware         *middleware.Middleware
}

func NewAuthService(helper *helper.Helper, authRepository *repository.AuthRepository, platformRepository *repository.PlatformRepository, middleware *middleware.Middleware) AuthService {
	return &authService{
		helper:             helper,
		middleware:         middleware,
		authRepository:     authRepository,
		platformRepository: platformRepository,
	}
}

func (s *authService) Login(traceCtx context.Context, c *fiber.Ctx, dto *dtos.Login) *responseDto.Response {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(traceCtx, "auth-service", "Service.Login")
	defer span.End()

	user, err := s.authRepository.User.GetUserByUsername(dto.Username)
	if err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusUnauthorized, "silahkan masukan email atau password yang terdaftar")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusUnauthorized, "silahkan masukan email atau password yang terdaftar")
	}

	claims := &jwtdto.CustomClaim{
		UserID:   user.ID,
		Platform: "user",
		// Role:   string(user.Role),
	}

	token := s.middleware.JWT.CreateTokenJwt(claims)

	response := &dtos.ResponseLogin{
		Token: token,
	}

	return s.helper.Response.JSONResponseSuccess(response, 0, 0, "success")

}

func (s *authService) LoginDevice(traceCtx context.Context, c *fiber.Ctx, dto *dtos.LoginDevice) *responseDto.Response {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(traceCtx, "auth-service", "Service.LoginDevice")
	defer span.End()

	// 1. VALIDASI FACTORY KEY
	// Pastikan Anda sudah setting FACTORY_KEY_DEVICE di file .env
	if dto.Secreet != os.Getenv("FACTORY_KEY_DEVICE") {
		errAuth := fiber.NewError(fiber.StatusUnauthorized, "Factory Key tidak valid")
		s.helper.Utils.JaegerTracer.RecordSpanError(span, errAuth)
		return s.helper.Response.JSONResponseError(fiber.StatusUnauthorized, "akses ditolak: factory key tidak valid")
	}

	// 2. CEK APAKAH DEVICE TERDAFTAR DI DATABASE (DeviceKafka)
	deviceData, err := s.platformRepository.DeviceKafka.GetOneBySn(dto.Sn)
	if err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, err.Error())
	}

	if deviceData.SN == "" {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, errors.New("sn tidak terdaftar"))
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed sn tidak terdaftar")
	}

	// 3. GENERATE TOKEN JIKA SEMUA VALID
	claims := &jwtdto.CustomClaim{
		Sn:       deviceData.SN, // Inject ID alat ke token
		Platform: "device",
	}

	token := s.middleware.JWT.CreateTokenJwt(claims)

	response := &dtos.ResponseLogin{
		Token: token,
	}

	return s.helper.Response.JSONResponseSuccess(response, 0, 0, "success")

}

func (s *authService) Register(traceCtx context.Context, c *fiber.Ctx, dto *dtos.Register) *responseDto.Response {

	_, span := s.helper.Utils.JaegerTracer.StartSpan(traceCtx, "auth-service", "Service.RegisterAdmin")
	defer span.End()

	if dto.Secreet != os.Getenv("FACTORY_KEY_USER") {
		return s.helper.Response.JSONResponseError(fiber.StatusUnauthorized, "silahkan masukan secreet")
	}

	userModel := new(models.User)

	user, err := s.authRepository.User.GetUserByUsername(dto.Username)
	if err != nil {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusUnauthorized, "silahkan masukan email atau username atau password yang terdaftar")
	}

	if user.Username != "" {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusUnauthorized, "silahkan masukan email atau username atau password yang terdaftar")
	}

	if dto.Password != dto.Confirm_Password {
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "silahkan masukan email atau username atau password yang terdaftar")
	}

	if dto.Name == "" {
		return s.helper.Response.JSONResponseError(fiber.StatusBadRequest, "name required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(fiber.StatusInternalServerError, err.Error())
		span.SetAttributes(attribute.String("error", "true"), attribute.String("error.message", "failed hash password"))
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed")
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed")
	}
	userModel.ID = uuid
	userModel.Username = dto.Username
	userModel.Password = string(hashedPassword)
	userModel.Name = dto.Name

	registerUser, err := s.authRepository.User.Register(userModel)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(fiber.StatusInternalServerError, err.Error())
		span.SetAttributes(attribute.String("error", "true"), attribute.String("error.message", "failed register user"))
		s.helper.Utils.JaegerTracer.RecordSpanError(span, err)
		return s.helper.Response.JSONResponseError(fiber.StatusInternalServerError, "failed")
	}

	claims := &jwtdto.CustomClaim{
		UserID:   registerUser.ID,
		Platform: "web",
	}

	token := s.middleware.JWT.CreateTokenJwt(claims)

	response := &dtos.ResponseLogin{
		Token: token,
	}

	return s.helper.Response.JSONResponseSuccess(response, 0, 0, "success")

}
