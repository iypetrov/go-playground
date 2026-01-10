package main

import (
	"context"
	"fmt"
)

type CartesianCoordinateSystem interface {
	GetCoordinates() (x, y, z float64, err error)
}

func GetCoordinates(ctx context.Context) (x, y, z float64, err error) {
	if ctx.Err() != nil {
		return 0, 0, 0, err
	}
	return 1.0, 2.0, 3.0, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	x, y, z, err := GetCoordinates(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Coordinates: x=%.2f, y=%.2f, z=%.2f\n", x, y, z)
}
