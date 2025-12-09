package main

import (
	"github.com/ride4Low/contracts/pkg/otel"
	"github.com/ride4Low/contracts/proto/trip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TripClient struct {
	conn       *grpc.ClientConn
	tripClient trip.TripServiceClient
}

func NewTripClient(target string) (*TripClient, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	dialOptions = append(dialOptions, otel.ClientOptions()...)

	conn, err := grpc.NewClient(target, dialOptions...)
	if err != nil {
		return nil, err
	}

	tripClient := trip.NewTripServiceClient(conn)

	return &TripClient{conn: conn, tripClient: tripClient}, nil
}

func (c *TripClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
