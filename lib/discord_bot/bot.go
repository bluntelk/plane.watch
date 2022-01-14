package discord_bot

// handles the discord bot integration

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type (
	PwBot struct {
		botId    string
		session  *discordgo.Session
		user     *discordgo.User
		commands map[string]*discordgo.ApplicationCommand

		userDms sync.Map
	}
)

const (
	slashCmdRegisterAddress  = "register-address"
	slashCmdRegisterLocation = "register-location"
	slashCmdRemoveLocation   = "remove-location"
	slashCmdListLocations    = "list-locations"
)

func NewBot(token string) (*PwBot, error) {
	b := PwBot{
		commands: make(map[string]*discordgo.ApplicationCommand),
	}
	var err error

	b.session, err = discordgo.New("Bot " + token)
	if nil != err {
		return nil, fmt.Errorf("failed to log into discord: %s", err)
	}

	b.user, err = b.session.User("@me")
	if nil != err {
		return nil, fmt.Errorf("failed to get my own user information: %s", err)
	}
	b.botId = b.user.ID
	log.Printf("I am %s", b.user)

	b.session.AddHandler(b.handleMessageCreate)
	//MessageReactionAdd, MessageReactionRemove, MessageUpdate, MessageDelete
	// MessageAck

	return &b, nil
}

func (b *PwBot) RegisterCommandsForGuild(guild *discordgo.UserGuild, nukeExisting bool) {
	log.Printf("Registering slash commands for %s", guild.Name)
	b.commands[slashCmdRegisterAddress] = &discordgo.ApplicationCommand{
		Name:        slashCmdRegisterAddress,
		Description: "We use the address you give to do a lookup to get your geo location",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "A name for you to know what this place is",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "address",
				Description: "The address of the place you want to register",
				Required:    true,
			},
		},
	}

	b.commands[slashCmdRegisterLocation] = &discordgo.ApplicationCommand{
		Name:        slashCmdRegisterLocation,
		Description: "Add a known Latitude/Longitude for your alert",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "A name for you to know what this place is",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "latitude",
				Description: "The latitude component",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "longitude",
				Description: "The longitude component",
				Required:    true,
			},
		},
	}

	b.commands[slashCmdRemoveLocation] = &discordgo.ApplicationCommand{
		Name:        slashCmdRemoveLocation,
		Description: "Removes a location given its name",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "name",
				Description:  "A name you gave this location",
				Required:     true,
				Autocomplete: true,
			},
		},
	}

	b.commands[slashCmdListLocations] = &discordgo.ApplicationCommand{
		Name:        slashCmdListLocations,
		Description: "Lists all your alert locations",
	}

	existingCommands, err := b.session.ApplicationCommands(b.session.State.User.ID, guild.ID)
	if nil != err {
		log.Printf("Unable to list existing slash commands: %s", err)
	}
	if nukeExisting {
		log.Printf("Nuking Existing Commands")
		for _, existingCommand := range existingCommands {
			if err = b.session.ApplicationCommandDelete(b.session.State.User.ID, guild.ID, existingCommand.ID); nil != err {
				log.Printf("Failed to remove command %s - %s", existingCommand.Name, err)
			}
		}
		existingCommands = []*discordgo.ApplicationCommand{}
	}

	for _, cmd := range b.commands {
		exists := false
		for _, existingCommand := range existingCommands {
			if existingCommand.Name == cmd.Name {
				//b.session.ApplicationCommandDelete(b.session.State.User.ID, guild.ID, existingCommand.ID)
				log.Printf("[%s] Command exists: %s - no need to create", guild.Name, cmd.Name)
				exists = true
				break
			}
		}
		if !exists {
			_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, guild.ID, cmd)
			if nil != err {
				log.Printf("[%s] Failed to register command %s: %s", guild.Name, cmd.Name, err)
			} else {
				log.Printf("[%s] Registered Command: %s", guild.Name, cmd.Name)
			}
		}
	}

	// todo: remove commands that we no longer want

	// now register our slash command listener
	b.session.AddHandler(b.handleSlashCommands)
}

func (b *PwBot) RegisterCommands(nukeExisting bool) {
	for _, guild := range b.GuildList() {
		b.RegisterCommandsForGuild(guild, nukeExisting)
	}
}

// GuildList returns the list of guilds the current user is associated with (handling the error from the API)
func (b *PwBot) GuildList() []*discordgo.UserGuild {
	list, err := b.session.UserGuilds(200, "", "")
	if nil != err {
		log.Printf("Failed to get list of user guilds for bot: %s", err)
		return []*discordgo.UserGuild{}
	}
	return list
}

// handleMessageCreate is our event handler for when the bot gets a message from a user
func (b *PwBot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == b.botId {
		return
	}
	// TODO: Send helpful help message when spoken to directly
	isPm := m.GuildID == "" && m.Member == nil
	if !isPm {
		return
	}

	switch strings.ToLower(m.Content) {
	case "list":
		b.sendUserLocationList(m.Author.ID)
	case "hi", "hai", "omg", "omghai", "yo", "sup", "wassup":
		_, err := s.ChannelMessageSend(m.ChannelID, "Hello there.")
		if nil != err {
			log.Printf("Failed to send message. %s", err)
		}
	default:
		var buf strings.Builder
		buf.WriteString("Hello, I am Birdy :small_airplane: the Plane.Watch bot\n")
		buf.WriteString("I can let you know when a plane fly's over your house (or whatever).\n")
		buf.WriteString("You can use any of my commands to interact with me\n")
		buf.WriteString("```\n")
		for _, cmd := range b.commands {
			buf.WriteString(fmt.Sprintf("  * /%-17s - %s\n", cmd.Name, cmd.Description))
		}
		buf.WriteString("```")
		buf.WriteString("`list` will also show you the list of configured alert locations, in case you don't want to use the command\n")

		_, err := s.ChannelMessageSend(m.ChannelID, buf.String())
		if nil != err {
			log.Printf("Failed to send message. %s", err)
		}
	}
}

func (b *PwBot) handleSlashCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var discordUserName, discordUserId string
	if nil != i.User {
		discordUserId = i.User.ID
		discordUserName = i.User.Username
	} else if nil != i.Member {
		discordUserId = i.Member.User.ID
		discordUserName = i.Member.User.Username
	} else {
		log.Printf("Failed to determine who I am interacting with")
		return
	}

	userChan, err := s.UserChannelCreate(discordUserId)
	if nil != err {
		log.Printf("Failed to create DM channel for user: %s", err)
	}
	sendPm := func(msg string) {
		if nil == userChan {
			log.Printf("No User Channel - would have said: %s", msg)
			return
		}
		_, err = s.ChannelMessageSend(userChan.ID, msg)
	}
	respondToInteraction := func(msg string) {
		if nil == i {
			return
		}
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
				Flags:   1 << 6,
			},
		})
		if nil != err {
			log.Printf("Failed to respond to interaction. %s", err)
		}

	}

	switch i.ApplicationCommandData().Name {
	case slashCmdRegisterAddress:
		name := i.ApplicationCommandData().Options[0].StringValue()
		addr := i.ApplicationCommandData().Options[1].StringValue()
		log.Printf("User %s is registering address %s: %s", discordUserName, name, addr)
		sendPm(fmt.Sprintf("You are registering address for `%s` ```%s```", name, addr))
		lat, lon, err := geoCodeAddress(addr)
		if nil != err {
			sendPm(err.Error())
			log.Printf("Geolocation failed %s", err)
			respondToInteraction("Geolocation failed, See PM")
		} else {
			if err = addAlertLocation(discordUserId, discordUserName, name, lat, lon); nil != err {
				respondToInteraction("Saving Failed, See PM")
				sendPm(fmt.Sprintf("There was an error adding your address. %s", err))
			} else {
				sendPm("We have setup your alert")
				if err = setLocationAddress(discordUserId, name, addr); nil != err {
					log.Printf("Failed to update user address: %s", err)
					respondToInteraction("Saving Partially Failed, See PM")
					sendPm("Failed to set alert locations address")
				} else {
					respondToInteraction("Success!, More details in PM")
				}
			}
		}
	case slashCmdRegisterLocation:
		name := i.ApplicationCommandData().Options[0].StringValue()
		lat := i.ApplicationCommandData().Options[1].FloatValue()
		lon := i.ApplicationCommandData().Options[2].FloatValue()

		log.Printf("user %s is registering %s: %0.5f,%0.5f", discordUserName, name, lat, lon)
		sendPm(fmt.Sprintf("Adding Location `%s` ```Lat: %0.5f, Lon: %0.5f```", name, lat, lon))

		if err = addAlertLocation(discordUserId, discordUserName, name, lat, lon); nil != err {
			sendPm(fmt.Sprintf("Failed adding that location. ```%s```", err))
			respondToInteraction("Failed, More info in PM")
		} else {
			sendPm("Successfully added alert location")
			respondToInteraction("Success!, More details in PM")
		}

	case slashCmdRemoveLocation:
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			name := i.ApplicationCommandData().Options[0].StringValue()

			log.Printf("user %s is removing location %s", discordUserName, name)
			sendPm(fmt.Sprintf("Removing Location `%s`", name))

			if err = removeAlertLocation(discordUserId, name); nil != err {
				sendPm(fmt.Sprintf("Failed removing location. ```%s```", err))
				respondToInteraction("Failed, see PM")
			} else {
				sendPm("Successfully removed alert location")
				respondToInteraction("Location Removed")
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			var choices []*discordgo.ApplicationCommandOptionChoice
			for _, loc := range getLocationsForUser(discordUserId) {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  loc.LocationName,
					Value: loc.LocationName,
				})
			}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{Choices: choices},
			})
			if err != nil {
				log.Printf("Failed to respond to interaction. %s", err)
			}
		}
	case slashCmdListLocations:
		b.sendUserLocationList(discordUserId)
		respondToInteraction("See PM For List")
	default:
		respondToInteraction("Unknown Command")
	}
}

func (b *PwBot) sendDirectMessage(discordUserId, message string) {
	if _, ok := b.userDms.Load(discordUserId); !ok {
		// need to create a user channel
		userChan, err := b.session.UserChannelCreate(discordUserId)
		if nil != err {
			log.Printf("Failed to create DM channel for user: %s", err)
			log.Printf("Would have told user [%s] %s", discordUserId, message)
			return
		}
		b.userDms.Store(discordUserId, userChan)
	}

	userChanInterface, ok := b.userDms.Load(discordUserId)
	if !ok {
		log.Printf("Failed to get the user DM channel from the map")
		log.Printf("Would have told user [%s] %s", discordUserId, message)
		return
	}

	userChan, ok := userChanInterface.(*discordgo.Channel)
	if !ok {
		b.userDms.Delete(discordUserId)
		log.Printf("What I got from the map was not the right thing!")
		log.Printf("Would have told user [%s] %s", discordUserId, message)
		return
	}

	_, err := b.session.ChannelMessageSend(userChan.ID, message)
	if nil != err {
		log.Printf("Failed to send message: %s", err)
		log.Printf("Would have told user [%s] %s", discordUserId, message)
		// if we are failing on sending, let's just nuke what we have and start again
		b.userDms.Delete(discordUserId)
	}
}

func (b *PwBot) sendUserLocationList(discordUserId string) {
	var buf strings.Builder
	buf.WriteString("*Configured Alert Locations*\n")
	buf.WriteString("```\n")
	tbl := tablewriter.NewWriter(&buf)
	tbl.SetColWidth(100)
	tbl.SetHeader([]string{"Name", "Address", "Lat", "Lon"})
	for _, loc := range getLocationsForUser(discordUserId) {
		tbl.Append([]string{
			loc.LocationName,
			loc.Address,
			fmt.Sprintf("%0.5f", loc.Lat),
			fmt.Sprintf("%0.5f", loc.Lon),
		})
	}

	tbl.Render()
	buf.WriteString("```")
	b.sendDirectMessage(discordUserId, buf.String())
}

func (b *PwBot) Run(c *cli.Context) error {
	log.Println("Running...")
	// load our existing alert config
	loadLocationsList()

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	err := b.session.Open()
	if nil != err {
		return err
	}

	b.RegisterCommands(c.Bool("nuke-commands"))

	if err = b.session.UpdateListeningStatus("ADSB"); nil != err {
		log.Printf("Unable to update listening to: %s", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Closing...")
	if err = saveLocationsList(); nil != err {
		log.Printf("Failed when saving locations list. %s", err)
	}
	return b.session.Close()
}
