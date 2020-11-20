package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/bwmarrin/discordgo"
)

var (
	svcEC2 *ec2.EC2
)

func init() {
	region := os.Getenv("AWS_REGION")
	if len(region) == 0 {
		fmt.Println("AWS_REGION not set; using ap-east-1 as default")
		region = "ap-east-1"
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("Error", err)
	}

	svcEC2 = ec2.New(sess)
}

func main() {
	key := os.Getenv("DISCORD_KEY")
	discord, err := discordgo.New("Bot " + key)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == session.State.User.ID {
		return
	}

	var reply string

	// TODO: Move this or the entire function to a separate file?
	switch msg.Content {
	case "good morning":
		reply = "<:kys:620483919774744596>"
	case "gsgo":
		reply = "<@!199174365408002049> <@!199462593264484352> <@!357042093199196162> <@!199917417953099778> <@!379870473481355264>"
	case "mc.start":
		if !stringInSlice(msg.Member.Roles, "779331911944372225") {
			reply = "You are not allowed to run this command."
			break
		}

		instanceID, err := getInstanceIDWithNameTag("minecraft-server")
		if err != nil {
			fmt.Println("Error", err)
			reply = "Error starting minecraft server: " + err.Error()
			break
		}
		result, err := startInstanceWithID(instanceID)
		if err != nil {
			fmt.Println("Error", err)
			reply = "Error starting minecraft server: " + err.Error()
			break
		}
		fmt.Println(result.StartingInstances)
		reply = "Starting Minecraft server..."
	}

	session.ChannelMessageSend(msg.ChannelID, reply)
}

func stringInSlice(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func getInstanceIDWithNameTag(name string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}

	result, err := svcEC2.DescribeInstances(input)
	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}
	if len(result.Reservations) == 0 {
		err := errors.New("Instance with Name \"" + name + "\" not found")
		fmt.Println("Error", err)
		return "", err
	}

	return *result.Reservations[0].Instances[0].InstanceId, nil
}

func startInstanceWithID(id string) (*ec2.StartInstancesOutput, error) {
	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}
	return svcEC2.StartInstances(input)
}
