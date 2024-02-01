package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"mado/weather"
)

const (
	CmdPoll      = "poll"
	CmdPolList   = "pollist"
	CmdPollHelp  = "pollhelp"
	CmdClosePoll = "closepoll"
	CmdWeather   = "weather"
)

var commandList = []string{CmdPoll, CmdPolList, CmdPollHelp, CmdClosePoll, CmdWeather}

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        CmdWeather,
			Description: "Get the weather for a location",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "location",
					Description: "Location to get weather for",
					Required:    true,
				},
			},
		},
		{
			Name:        CmdPoll,
			Description: "basic command route for starting a poll",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "question for the poll",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "multiple-options",
					Description: "able to cast multiple votes",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer1",
					Description: "first answer",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer2",
					Description: "second answer",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-3",
					Description: "third answer",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-4",
					Description: "fourth answer",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "answer-5",
					Description: "fifth answer",
					Required:    false,
				},
			},
		},
		{
			Name:        CmdPolList,
			Description: "List all open polls",
		},
		{
			Name:        CmdPollHelp,
			Description: "get help on all commands",
		},
		{
			Name:        CmdClosePoll,
			Description: "Close a poll by id",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "poll-id",
					Description: "Poll id",
					Required:    true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		CmdPoll: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{}
			msgformat := "New poll: \n"
			if len(i.ApplicationCommandData().Options) >= 3 {
				for j, opt := range i.ApplicationCommandData().Options {
					if opt.Name == "question" {
						msgformat += "question: %s \n"
						margs = append(margs, opt.StringValue())
					} else if opt.Name == "multiple-options" { // Change here
						msgformat += "> multipleOptions: %v\n"
						margs = append(margs, opt.BoolValue())
					} else {
						msgformat += fmt.Sprintf("answer %d", j)
						msgformat += ": %v\n"
						margs = append(margs, opt.StringValue())
					}
				}
				margs = append(margs, i.ApplicationCommandData().Options[0].StringValue())
				msgformat += "> poll-id: <#%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		CmdPolList: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "List of all open polls",
				},
			})
		},
		CmdPollHelp: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msgFormat := "All available commands: \n"
			cmdFormat := "/%s \n"
			for _, c := range commandList {
				msgFormat += fmt.Sprintf(cmdFormat, c)
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msgFormat,
				},
			})
		},
		CmdClosePoll: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				// Here we need to convert raw interface{} value to wanted type.
				// Also, as you can see, here is used utility functions to convert the value
				// to particular type. Yeah, you can use just switch type,
				// but this is much simpler
				i.ApplicationCommandData().Options[0].StringValue(),
			}
			msgformat :=
				` Attempting to close:
				> poll-id: %s
`
			if len(i.ApplicationCommandData().Options) >= 2 {
				margs = append(margs, i.ApplicationCommandData().Options[0].StringValue())
				msgformat += "> poll-id: <#%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		CmdWeather: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			fmt.Println("HERE")
			location := i.ApplicationCommandData().Options[0].StringValue()
			resName, resTempC, resCondition := weather.GetWeather(location)
			response := fmt.Sprintf("Location: %s\nTemperature: %f°C\nCondition: %s", resName, resTempC, resCondition)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
	}
)
