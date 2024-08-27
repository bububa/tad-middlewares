# tad-middlewares

腾讯广告SDK 中间件

[![Go Reference](https://pkg.go.dev/badge/github.com/bububa/tad-middlewares.svg)](https://pkg.go.dev/github.com/bububa/tad-middlewares)
[![Go](https://github.com/bububa/tad-middlewares/actions/workflows/go.yml/badge.svg)](https://github.com/bububa/tad-middlewares/actions/workflows/go.yml)
[![goreleaser](https://github.com/bububa/tad-middlewares/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/bububa/tad-middlewares/actions/workflows/goreleaser.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/bububa/tad-middlewares.svg)](https://github.com/bububa/tad-middlewares)
[![GoReportCard](https://goreportcard.com/badge/github.com/bububa/tad-middlewares)](https://goreportcard.com/report/github.com/bububa/tad-middlewares)
[![GitHub license](https://img.shields.io/github/license/bububa/tad-middlewares.svg)](https://github.com/bububa/tad-middlewares/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/bububa/tad-middlewares.svg)](https://GitHub.com/bububa/tad-middlewares/releases/)

## Opentelementry Middleware

```golang
import (
  "github.com/tencentad/marketing-api-go-sdk/pkg/ads/v3"
  "github.com/tencentad/marketing-api-go-sdk/pkg/config/v3"
)
func main() {
  tads:= ads.Init(&config.SDKConfig{})
  mw := tadmw.NewOtelMiddlware(clt, "")
  tads.AppendMiddleware(mw)
  // your client id
	clientId := int64(0)
	clientSecret := "your client secret"
	grantType := "authorization_code"
	oauthTokenOpts := &api.OauthTokenOpts{
		AuthorizationCode: optional.NewString("your authorization code"),
		RedirectUri: optional.NewString("your authorization code"),
	}
	ctx := *tads.Ctx
	// oauth/token接口即对应Oauth().Token()方法
	response, _, err := tads.Oauth().Token(ctx, clientId, clientSecret, grantType, oauthTokenOpts)

	if err != nil {
		if resErr, ok := err.(errors.ResponseError); ok {
			errStr, _ := json.Marshal(resErr)
			// TODO for api error
			fmt.Println("Response error:", string(errStr))
		} else {
			// TODO for other error
			fmt.Println("Error:", err)
		}
	}
	tads.SetAccessToken(response.AccessToken)
}
```

```

```
