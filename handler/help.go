package handler

// Let's the bot provide a help message in the DMs
func Help() []string {
	helpMsg := make([]string, 0)
	helpMsg = append(helpMsg, `General channel commands:
   [!userinfo|!info]             - Tells you, if you are registered with the myself
   [!register|!reg] <@nick>      - Registers <@nick> as user in the bot's database for them to access
                                   some advances features (like the SoT commands)
   [!unregister|!unreg] <@nick>  - Delete <@nick> from registered user's database`)

	return helpMsg
}
