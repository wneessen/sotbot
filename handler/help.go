package handler

// Let's the bot provide a help message in the DMs
func Help() []string {
	helpMsg := make([]string, 0)
	helpMsg = append(helpMsg, `General channel commands:
   !play <sound_name>            - Jump into the voice chat of the requesting user and play the requested sound
   !fact                         - Replies with a random useless fact
   !movie <query>                - Look up the <query> in the TMDb movie database. If <query> is omitted, a
                                   random movie is choosen
   !time                         - Replies with the current time
   !tv <query>                   - Look up the <query> in the TMDb TV series database. If <query> is omitted, a
                                   random TV series is choosen
   [!urban|!ud] <term>           - Looks up the <term> in Urban Dictionary. If <term> is omitted, a random 
                                   term is choosen
   !weather <location>           - Look up the current weather in <location> and provide it to the user
   [!uptime|!up]                 - Will show you my current uptime
   [!userinfo|!info]             - Tells you, if you are registered with the myself
   [!version|!ver]               - Provide some version info about myself`, `Admin-only channel commands:
   [!register|!reg] <@nick>      - Registers <@nick> as user in the bot's database for them to access 
                                   some advances features (like the SoT commands)
   [!unregister|!unreg] <@nick>  - Delete <@nick> from registered user's database
   [!memory|!mem]                - Show some memory usage information`, `Sea of Thieves specific commands:
   [!achievement|!achieve]       - Replies with your latest SoT achievement[1][2]
   [!balance|!bal]               - Replies your current SoT gold balance[1][2]
   !code                         - Replies with a random SoT pirate code article
   [!deed|!dd]                   - Replies with the currently active "daily deed" in SoT
   [!ledger|!led] <fac.>         - Replies with your current ledger position/rank in the requested faction[1][2]
   [!reputation|!rep] <fac.>     - Replies with your current reputation in the requested faction[1][2]
   !season                       - Replies with your current SoT season progress[1][2]
   !tr				 - Replies with the current traderoutes from rarethief
   [!setrat|!rat] <cookie>       - Set/update your SoT RAT cookie in the bot's DB[1][3]
   [!stats|!stat]                - Provides some general SoT stats for your user (killed kraken, sold chests, etc.)[1][2]

Note about the RAT cookie:
To extract your cookie for the API, you need the SoT-RAT-Extractor (https://github.com/echox/sot-rat-extractor).
Please keep in mind, that you are providing your MS Live login information to the RAT Extractor and once you use
the !setrat feature, your cookie will be stored unencrypted for our database.

[1] = Registered users only
[2] = SoT RAT cookie is set in the DB
[3] = Command needs to be issued in a DM, not in a public channel`)

	return helpMsg
}
