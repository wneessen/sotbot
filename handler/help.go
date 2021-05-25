package handler

// Let's the bot provide a help message in the DMs
func Help() string {
	helpMsg := `General channel commands:
   !version                  - Provide some info about myself
   !play <sound_name>        - Jump into the voice chat of the requesting user and play the requested sound
   !time                     - Replies with the current time
   !fact                     - Replies with a random useless fact
   !code                     - Replies with a random SoT pirate code article

Admin-only channel commands:
   !register <@nick>         - Register a user in the database for access to advances features
   !unregister <@nick>       - Unregister a registered user

Registered-only channel commands:
   !balance                  - Replies your current SoT gold balance (registered-only)
   !achievement              - Replies with your latest SoT achievement (registered-only)
   !season                   - Replies with your current SoT season progress (registered-only)
   !rep <faction>            - Replies with your current reputation in the requested faction (registered-only)

DM-only commands:
   !userinfo                 - Tells you, if you are registered with the myself
   !setrat                   - Set/update your SoT RAT cookie (registered-only)

Note about the RAT cookie:
To extract your cookie for the API, you need the SoT-RAT-Extractor (https://github.com/echox/sot-rat-extractor).
Please keep in mind, that you are providing your MS Live login information to the RAT Extractor and once you use
the !setrat feature, your cookie will be stored unencrypted for our database.`

	return "\n`" + helpMsg + "`"
}
