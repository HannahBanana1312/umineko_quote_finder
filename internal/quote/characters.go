package quote

var CharacterNames = map[string]string{
	"00": "GroupVoices",
	"01": "Kinzo",
	"02": "Krauss",
	"03": "Natsuhi",
	"04": "Jessica",
	"05": "Eva",
	"06": "Hideyoshi",
	"07": "George",
	"08": "Rudolf",
	"09": "Kyrie",
	"10": "Battler",
	"11": "Ange",
	"12": "Rosa",
	"13": "Maria",
	"14": "Genji",
	"15": "Shannon",
	"16": "Kanon",
	"17": "Gohda",
	"18": "KumasawaChiyo",
	"19": "NanjoTerumasa",
	"20": "Amakusa",
	"21": "Okonogi",
	"22": "Kasumi",
	"23": "ProfessorOotsuki",
	"24": "CaptainKawabata",
	"25": "NanjoMasayuki",
	"26": "KumasawaSabakichi",
	"27": "Beatrice",
	"28": "Bernkastel",
	"29": "Lambdadelta",
	"30": "Virgilia",
	"31": "Ronove",
	"32": "Gaap",
	"33": "Sakutarou",
	"34": "Evatrice",
	"35": "Chiester45",
	"36": "Chiester410",
	"37": "Chiester00",
	"38": "Lucifer",
	"39": "Leviathan",
	"40": "Satan",
	"41": "Belphegor",
	"42": "Mammon",
	"43": "Beelzebub",
	"44": "Asmodeus",
	"45": "Goat",
	"46": "Erika",
	"47": "Dlanor",
	"48": "Gertrude",
	"49": "Cornelia",
	"50": "Featherine",
	"51": "Zepar",
	"52": "Furfur",
	"53": "Lion",
	"54": "Willard",
	"55": "Claire",
	"56": "Ikuko",
	"57": "Tohya",
	"58": "KinzoYoung",
	"59": "BiceChickBeato",
	"60": "BeatoElder",
	"99": "MiscVoices",
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
