package handler

// Let's the bot provide a help message in the DMs
func Help() []string {
	helpMsg := make([]string, 0)
	helpMsg = append(helpMsg, `General channel commands:
   !movie <query>                - Look up the <query> in the TMDb movie database. If <query> is omitted, a
                                   random movie is chosen
   !tv <query>                   - Look up the <query> in the TMDb TV series database. If <query> is omitted, a
                                   random TV series is chosen
   [!urban|!ud] <term>           - Looks up the <term> in Urban Dictionary. If <term> is omitted, a random
                                   term is chosen
   [!userinfo|!info]             - Tells you, if you are registered with the myself
   [!register|!reg] <@nick>      - Registers <@nick> as user in the bot's database for them to access
                                   some advances features (like the SoT commands)
   [!unregister|!unreg] <@nick>  - Delete <@nick> from registered user's database
   [!setrat|!rat] <cookie>       - Set/update your SoT RAT cookie in the bot's DB[1][3]

Note about the RAT cookie:
To extract your cookie for the API, you need the SoT-RAT-Extractor (https://github.com/echox/sot-rat-extractor).
Please keep in mind, that you are providing your MS Live login information to the RAT Extractor and once you use
the !setrat feature, your cookie will be stored unencrypted for our database.

[1] = Registered users only
[2] = SoT RAT cookie is set in the DB
[3] = Command needs to be issued in a DM, not in a public channel`)

	return helpMsg
}
