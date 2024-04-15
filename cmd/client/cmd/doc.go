// Client command tool allow user to interract
// with the speech-to-text server

/*
Interaction with the server requires user to be registered and
and logged in.

# Register
Use subcommand `register` with parameters to create new user at the
server side, e.g:
` register --username u --email e --password p `

# Login
After successful registration user allowed to log in with his
user/password pair. Server on successful log in return an JWT token
expiring in 3 hours.

# Transcribe
Subcommand `transcribe` transcribes audio to text

# History
Subcommand `history` returns text of all the previous transcriptions for
a current user
*/

package cmd
