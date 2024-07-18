package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/gitsang/notation/pkg/soundslice"
	"github.com/joho/godotenv"
)

func f1(client *soundslice.Client) {
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

	err = client.DeleteNotation(ctx, sliceId)
	if err != nil {
		panic(err)
	}
	fmt.Println("DeleteNotation success.")
}

func f2(client *soundslice.Client) {
	scores, err := client.ListScores()
	if err != nil {
		panic(err)
	}
	scoresJSON, _ := json.Marshal(scores)
	fmt.Printf("ListScores success. scores: %+v\n", string(scoresJSON))

}

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

	f2(client)
}
