package base

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/golang-jwt/jwt"
	"github.com/rigdev/rig-go-api/api/v1/project"
	"github.com/rigdev/rig-go-sdk"
	"github.com/rigdev/rig/pkg/service"
	"github.com/rigdev/rig/pkg/uuid"
	"go.uber.org/fx"
)

var clientModule = fx.Module("client",
	fx.Supply(&http.Client{}),
	fx.Provide(func(s *Service, cfg *Config) rig.Client {
		ai := &authInterceptor{cfg: cfg}
		nc := rig.NewClient(
			rig.WithHost(s.Server),
			rig.WithInterceptors(ai, &userAgentInterceptor{}),
			rig.WithSessionManager(&configSessionManager{cfg: cfg}),
		)
		ai.nc = nc
		return nc
	}),
	fx.Provide(func(cfg *Config) []connect.Interceptor {
		return []connect.Interceptor{&userAgentInterceptor{}, &authInterceptor{cfg: cfg}}
	}),
)

type userAgentInterceptor struct{}

func (i *userAgentInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		i.setUserAgent(ar.Header())
		return next(ctx, ar)
	}
}

func (i *userAgentInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, s)
		i.setUserAgent(conn.RequestHeader())
		return conn
	}
}

func (i *userAgentInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		i.setUserAgent(shc.RequestHeader())
		return next(ctx, shc)
	}
}

func (i *userAgentInterceptor) setUserAgent(h http.Header) {
	h.Set("User-Agent", "Rig-CLI/v0.0.1")
}

type configSessionManager struct {
	cfg *Config
}

func (s *configSessionManager) GetAccessToken() string {
	return s.cfg.Auth().AccessToken
}

func (s *configSessionManager) GetRefreshToken() string {
	return s.cfg.Auth().RefreshToken
}

func (s *configSessionManager) SetAccessToken(accessToken, refreshToken string) {
	s.cfg.Auth().AccessToken = accessToken
	s.cfg.Auth().RefreshToken = refreshToken
	if err := s.cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error saving config: %v\n", err)
	}
}

type authInterceptor struct {
	cfg *Config
	nc  rig.Client
}

func (i *authInterceptor) handleAuth(ctx context.Context, h http.Header, method string) {
	if _, ok := service.OmitProjectToken[method]; !ok {
		i.setProjectToken(ctx, h)
	}
}

func (i *authInterceptor) setProjectToken(ctx context.Context, h http.Header) {
	if i.cfg.Context().Project.ProjectToken == "" {
		return
	}

	c := jwt.StandardClaims{}
	p := jwt.Parser{
		SkipClaimsValidation: true,
	}
	_, _, err := p.ParseUnverified(
		i.cfg.Context().Project.ProjectToken,
		&c,
	)
	if err != nil {
		return
	}

	// Don't use if invalid user id.
	if i.cfg.Auth().UserID.String() != c.Subject {
		return
	}

	if !c.VerifyExpiresAt(time.Now().Add(30*time.Second).Unix(), true) && i.cfg.Context().Project.ProjectID != uuid.Nil {
		res, err := i.nc.Project().Use(ctx, &connect.Request[project.UseRequest]{
			Msg: &project.UseRequest{
				ProjectId: i.cfg.Context().Project.ProjectID.String(),
			},
		})
		if err == nil {
			i.cfg.Context().Project.ProjectToken = res.Msg.GetProjectToken()
			i.cfg.Save()
		}
	}

	h.Set(service.RigProjectTokenHeader, fmt.Sprint(i.cfg.Context().Project.ProjectToken))
}

func (i *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		i.handleAuth(ctx, ar.Header(), ar.Spec().Procedure)
		return next(ctx, ar)
	}
}

func (i *authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, s connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, s)
		i.handleAuth(ctx, conn.RequestHeader(), s.Procedure)
		return conn
	}
}

func (i *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		i.handleAuth(ctx, shc.RequestHeader(), shc.Spec().Procedure)
		return next(ctx, shc)
	}
}