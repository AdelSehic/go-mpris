package mpris

import "github.com/godbus/dbus/v5"

type Metadata struct {
	TrackID string   `json:"track_id"`
	Title   string   `json:"title"`
	Album   string   `json:"album"`
	Artist  []string `json:"artist"`
	Length  int64    `json:"length"` // microseconds
	ArtURL  string   `json:"art_url"`
	URL     string   `json:"url"`
}

func parseMetadata(variantMap map[string]dbus.Variant) *Metadata {
	getString := func(v dbus.Variant) string {
		if val, ok := v.Value().(string); ok {
			return val
		}
		return ""
	}

	getStringArray := func(v dbus.Variant) []string {
		if val, ok := v.Value().([]string); ok {
			return val
		}
		return nil
	}

	getInt64 := func(v dbus.Variant) int64 {
		switch val := v.Value().(type) {
		case int64:
			return val
		case uint64:
			return int64(val)
		}
		return 0
	}

	return &Metadata{
		TrackID: getString(variantMap["mpris:trackid"]),
		Title:   getString(variantMap["xesam:title"]),
		Album:   getString(variantMap["xesam:album"]),
		Artist:  getStringArray(variantMap["xesam:artist"]),
		Length:  getInt64(variantMap["mpris:length"]),
		ArtURL:  getString(variantMap["mpris:artUrl"]),
		URL:     getString(variantMap["xesam:url"]),
	}
}
