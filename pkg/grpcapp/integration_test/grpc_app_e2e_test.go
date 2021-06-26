package integration

import (
	"context"
	"github.com/stretchr/testify/suite"
	"goboilerplate/pkg/grpcapp/pb"
	"goboilerplate/pkg/service/postservice/mocks"
	"goboilerplate/pkg/testhelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"
	"net"
	"testing"
)

var testDB *gorm.DB

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type e2eTestSuite struct {
	suite.Suite
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(e2eTestSuite))
}

func (suite *e2eTestSuite) SetupTest() {
	const bufSize = 1024 * 1024
	testDB = testhelper.InitSpecificTestDB("goboilerplate_grpc_e2e_test")
	testhelper.SeedFixtures(suite.T(), testDB, "testdata/fixtures")
	lis = bufconn.Listen(bufSize)
	grpcApp, err := InitializeGRPCApp(testDB, &mocks.UuidService{})
	if err != nil {
		suite.FailNow("Couldn't initialize GRPC app", err)
	}

	go grpcApp.Run(lis)
}

func setupConnection(ctx context.Context, suite *e2eTestSuite) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		suite.FailNow("Failed to dial bufnet", err)
	}

	return conn
}

func (suite *e2eTestSuite) TestGetAllPosts() {
	conn := setupConnection(context.Background(), suite)
	client := pb.NewPostsServiceClient(conn)
	req := &pb.GetAllPostsRequest{}
	resp, err := client.GetAllPosts(context.Background(), req)
	expected := &pb.GetAllPostsResponse{
		Posts: []*pb.Post{
			{
				Uuid:    "uuid-1",
				Content: "Post 1",
			},
			{
				Uuid:    "uuid-2",
				Content: "Post 2",
			},
		},
	}

	suite.Equal(expected, resp)
	suite.Nil(err)
}
