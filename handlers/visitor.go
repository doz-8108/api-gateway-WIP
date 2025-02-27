package handlers

import (
	"flag"
	"regexp"
	"strings"

	"github.com/doz-8108/api-gateway/pb"
	"github.com/doz-8108/api-gateway/utils"
	"github.com/gofiber/fiber/v3"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	VisitorHandlers struct {
		Utils utils.Utils
	}
	IncrementVisitorCountReqBody struct {
		IpAddr string `json:"ip_addr"`
	}
)

var (
	visitorCountSvc = flag.String("visitor-counter-svc", "localhost:8081", "address of visitor counting service")
	ipv4Regex       = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	ipv6Regex       = regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}([0-9a-fA-F]{1,4}|:)$`)
)

func NewVisitorHandlers(utils utils.Utils) *VisitorHandlers {
	return &VisitorHandlers{Utils: utils}
}

func (v *VisitorHandlers) GetVisitorCounts(f fiber.Ctx) error {
	conn, ctx, cancel := v.Utils.GrpcConnect(*visitorCountSvc)
	defer v.Utils.GrpcDisConnect(conn, cancel)

	client := pb.NewVisitorCounterServiceClient(conn)
	resp, err := client.GetVisitorCounts(ctx, &emptypb.Empty{})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch visitor counts")
	}

	return f.Status(fiber.StatusOK).JSON(fiber.Map{"visitor_counts": resp.VisitorCounts})
}

func (v *VisitorHandlers) IncrementVisitorCount(f fiber.Ctx) error {
	incrementVisitorCountReqBody := new(IncrementVisitorCountReqBody)
	if err := f.Bind().Body(incrementVisitorCountReqBody); err != nil {
		v.Utils.CatchError(err, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if strings.Trim(incrementVisitorCountReqBody.IpAddr, " ") == "" {
		v.Utils.CatchError(nil, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	if !(ipv4Regex.MatchString(incrementVisitorCountReqBody.IpAddr) || ipv6Regex.MatchString(incrementVisitorCountReqBody.IpAddr)) {
		v.Utils.CatchError(nil, fiber.NewError(fiber.StatusBadRequest, "invalid request body"))
	}

	conn, ctx, cancel := v.Utils.GrpcConnect(*visitorCountSvc)
	defer v.Utils.GrpcDisConnect(conn, cancel)

	client := pb.NewVisitorCounterServiceClient(conn)
	_, err := client.IncrementVisitorCount(ctx, &pb.IncrementVisitorCountRequest{IpAddr: incrementVisitorCountReqBody.IpAddr})
	if err != nil {
		v.Utils.CatchError(err, fiber.NewError(fiber.StatusInternalServerError, "failed to increment visitor count"))
	}

	return f.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}
