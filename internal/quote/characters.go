package quote

var CharacterNames = map[string]string{
	"00": "Servants",
	"01": "Kinzo",
	"02": "Krauss",
	"03": "Natsuhi",
	"04": "Battler",
	"05": "Rudolf",
	"06": "Hideyoshi",
	"07": "Kyrie",
	"08": "Eva",
	"09": "Rosa",
	"10": "Battler",
	"11": "Battler",
	"12": "Rosa",
	"13": "Maria",
	"14": "Genji",
	"15": "Shannon",
	"16": "Kanon",
	"17": "Kumasawa",
	"18": "Gohda",
	"19": "Nanjo",
	"20": "Amakusa",
	"21": "Kasumi",
	"22": "Lambdadelta",
	"23": "Okonogi",
	"24": "Kawabata",
	"25": "Tetsuro",
	"26": "Sabakichi",
	"27": "Beatrice",
	"28": "Ange",
	"29": "Bernkastel",
	"30": "Virgilia",
	"31": "Ronove",
	"32": "Goldsmith",
	"33": "Sakutaro",
	"34": "Eva-Beatrice",
	"35": "Chiester 45",
	"36": "Chiester 410",
	"37": "Chiester 00",
	"38": "Lucifer",
	"39": "Leviathan",
	"40": "Satan",
	"41": "Belphegor",
	"42": "Mammon",
	"43": "Beelzebub",
	"44": "Asmodeus",
	"45": "Gaap",
	"46": "Jessica",
	"47": "Dlanor",
	"48": "Gertrude",
	"49": "Cornelia",
	"50": "Featherine",
	"51": "Zepar",
	"52": "Furfur",
	"53": "Lion",
	"54": "Will",
	"55": "Confession",
	"56": "Tohya",
	"57": "Ikuko",
	"58": "Young Kinzo",
	"59": "Clair",
	"60": "Piece",
	"99": "Narrator",
}

func GetCharacterName(id string) string {
	if name, ok := CharacterNames[id]; ok {
		return name
	}
	return "Unknown"
}

func GetAllCharacters() map[string]string {
	return CharacterNames
}
