package mpris

const (
	SERVICE_DBUS  = "org.freedesktop.DBus"
	SERVICE_MPRIS = "org.mpris.MediaPlayer2"

	PLAYER_OBJECT = SERVICE_MPRIS + ".Player"

	PLAYER_PLAY     = PLAYER_OBJECT + ".Play"
	PLAYER_PAUSE    = PLAYER_OBJECT + ".Pause"
	PLAYER_STATUS   = PLAYER_OBJECT + ".PlaybackStatus"
	PLAYER_METADATA = PLAYER_OBJECT + ".Metadata"

	MPRIS_LIST_PLAYERS = SERVICE_DBUS + ".ListNames"

	PROPERTY_OWNER_CHANGE = "NameOwnerChanged"
	OBJECT_MPRIS          = "/org/mpris/MediaPlayer2"

	STATUS_PLAYING = "Playing"
	STATUS_PAUSED  = "Paused"
)
