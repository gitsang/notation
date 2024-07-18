package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gitsang/notation/pkg/soundslice"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	sesn := os.Getenv("SESN")

	client := soundslice.NewClient(
		soundslice.WithLogHandler(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			})),
		soundslice.WithAddr("https://www.soundslice.com"),
		soundslice.WithSesn(sesn),
	)

	ctx := context.Background()

	sliceId, err := client.CreateNotation(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("CreateNotation success. sliceId:", sliceId)

	// open file
	fh, err := os.Open("./Yellow.gp")
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	uploadResp, err := client.UploadNotation(ctx, sliceId, fh)
	if err != nil {
		panic(err)
	}
	fmt.Println("UploadNotation success.", uploadResp)
}
