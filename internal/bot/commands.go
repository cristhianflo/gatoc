package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	Metadata *discordgo.ApplicationCommand
	Handler  func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *BotContext) error
}

type SlashSubcommand struct {
	Metadata *discordgo.ApplicationCommandOption
	Handler  func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *BotContext) error
}

func GetInteractionFailedResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	message := "Ha ocurrido un error ejecutando el comando."

	if content != "" {
		message = content
	}
	return s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		},
	)
}

func DeferReply(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func EditDeferred(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	_, err := s.InteractionResponseEdit(
		i.Interaction,
		&discordgo.WebhookEdit{
			Content: &content,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func GetInteractionOptionString(name string, i *discordgo.InteractionCreate) (string, error) {
	options := i.ApplicationCommandData().Options

	for _, opt := range options {
		if opt.Name == name {
			return opt.StringValue(), nil
		}
	}
	return "", errors.New("Error option name does not exist in interaction options")
}

// 	videoID := strings.Split(query, "v=")[1]
// 	stream, err := getAudioStreamURL(videoID)

// 	if err != nil {
// 		fmt.Println("failed on youtube download:", err)
// 		content = "Ha ocurrido un error"
// 		s.InteractionResponseEdit(
// 			i.Interaction,
// 			&discordgo.WebhookEdit{
// 				Content: &content,
// 			},
// 		)
// 	}

// 	opusStream, err := convertMP4AToOpus(stream)
// 	if err != nil {
// 		fmt.Println("failed to convert to opus", err)
// 		content = "Ha ocurrido un error"
// 		s.InteractionResponseEdit(
// 			i.Interaction,
// 			&discordgo.WebhookEdit{
// 				Content: &content,
// 			},
// 		)
// 	}
// 	defer opusStream.Close()

// 	fmt.Println("Successfully started ffmpeg for Ogg Opus conversion. Reading stream...")
// 	ogg, _, err := oggreader.NewWith(opusStream)
// 	if err != nil {
// 		panic(err)
// 	}

// 	buffer := make([]byte, 4096)

// 	decoder := opus.NewDecoder()

// 	for {
// 		segments, _, err := ogg.ParseNextPage()

// 		if errors.Is(err, io.EOF) {
// 			break
// 		} else if bytes.HasPrefix(segments[0], []byte("OpusTags")) {
// 			continue
// 		}

// 		if err != nil {
// 			panic(err)
// 		}

// 		for i, segment := range segments {
// 			if _, _, err = decoder.Decode(segments[i], out); err != nil {
// 				panic(err)
// 			}
// 		}
// 	}

// 	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
// 		Content: &content,
// 	})

// }

// func getAudioStreamURL(videoID string) (io.ReadCloser, error) {
// 	client := youtube.Client{}

// 	video, err := client.GetVideo(videoID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	formats := video.Formats.WithAudioChannels() // only get videos with audio
// 	stream, _, err := client.GetStream(video, &formats[0])
// 	if err != nil {
// 		return nil, err
// 	}

// 	return stream, nil
// }

// func convertMP4AToOpus(input io.ReadCloser) (io.ReadCloser, error) {
// 	// Construct the ffmpeg command
// 	cmd := exec.Command(
// 		"ffmpeg",
// 		"-i", "-", // Read from stdin
// 		"-vn", // No video
// 		"-acodec", "libopus",
// 		"-f", "ogg", // Output to Ogg container
// 		"-", // Output to stdout
// 	)

// 	// Set the input pipe to the provided io.ReadCloser
// 	cmd.Stdin = input

// 	// Get the output pipe
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
// 	}

// 	// Get the error pipe (optional, but good for debugging)
// 	stderr, err := cmd.StderrPipe()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
// 	}

// 	// Start the ffmpeg command
// 	if err := cmd.Start(); err != nil {
// 		return nil, fmt.Errorf("failed to start ffmpeg: %w", err)
// 	}

// 	// Handle stderr in a goroutine to avoid blocking
// 	go func() {
// 		errOutput, _ := io.ReadAll(stderr)
// 		if len(errOutput) > 0 {
// 			fmt.Printf("ffmpeg stderr: %s\n", string(errOutput))
// 		}
// 	}()

// 	// Return the stdout pipe as the opus audio stream
// 	return stdout, nil
// }
