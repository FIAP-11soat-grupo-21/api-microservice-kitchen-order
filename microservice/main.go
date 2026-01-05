// @title Tech Challenge API
// @version 1.0
// @description This is an API for a tech challenge.
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8082
// @BasePath /v1
// @schemes http
//
//go:debug x509negativeserial=1
package main

import "tech_challenge/internal/shared/infra/api"

func main() {
	api.Init()
}
