package routers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/influxdata/influxdb/v2/pkg/snowflake"
	"github.com/zercle/docker-snowflake/pkg/datamodels"
)

var snowGen *snowflake.Generator

func InitSnowflake() {
	rand.Seed(time.Now().UnixNano())
	snowGen = snowflake.New(rand.Intn(1023))
}

func SetRouters(app *fiber.App) {

	InitSnowflake()

	app.Get("/:machineID?", func(c *fiber.Ctx) (err error) {
		machineID, _ := strconv.Atoi(c.Params("machineID", "0"))
		if machineID < 0 || machineID > 1023 {
			return fiber.NewError(http.StatusBadRequest, "machineID must be a number between (inclusive) 0 and 1023")
		} else {
			snowGen = snowflake.New(machineID)
		}
		snowID := snowGen.Next()

		responseForm := datamodels.ResponseForm{
			Success: bool(err == nil),
			Result: fiber.Map{
				"id":  snowID,
				"hex": fmt.Sprintf("%x", snowID),
			},
		}

		return c.JSON(responseForm)
	})
}
