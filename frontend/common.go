package frontend

var Commands = []string{
	"play,add-play the song | play <song name>",
	"search,s-search the song and display search result | search <song name>",
	"radio-start radio for song | radio <song name>",
	"pause,resume,p-toggle pause/resume",
	"showq,q-display song queue",
	"curr,c-display current song",
	"skipn,n-skip to next song",
	"skipb,b-skip to previous song",
	"skip-skip to the specified index, default is 1 | skip <index>",
	"remove,rem-remove song at specified index, default is last | remove <index>",
	"removeAll,reml-remove all songs stating from at specified index, default is current+1 | removeAll <index>",
	"forward,f-forwads playback by 10s ** | forward <seconds>",
	"rewind,r-rewinds playback by 10s ** | rewind <seconds>",
	"setVol,v-sets the volume by amount (0-100) | setVol <volume>",
	"stop-resets the player",
	"listSongs,ls-displays list of songs based on criteria (recent,likes,plays) | listSongs <criteria>",
	"checkApi-check the current piped api",
	"setApi-set new piped api | setApi <piped api>",
	"listApi-display all available instances",
	"randApi-randomly select an piped instance",
	"version-display application details",
	"quit-quit application",
}

var Configs = []string{
	"config.piped.apiUrl-default piped api to be used",
	"config.piped.instanceListApi-default instance list api to be used",
	"config.cache.enabled-enable/disable audio caching, enabled by default",
	"config.cache.path-path to audio caching",
	"config.database.path-path to db",
	"config.source.isPiped-enable piped as default source for audio searching",
}
