package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	bytewords               = "ableacidalsoapexaquaarchatomauntawayaxisbackbaldbarnbeltbetabiasbluebodybragbrewbulbbuzzcalmcashcatschefcityclawcodecolacookcostcruxcurlcuspcyandarkdatadaysdelidicedietdoordowndrawdropdrumdulldutyeacheasyechoedgeepicevenexamexiteyesfactfairfernfigsfilmfishfizzflapflewfluxfoxyfreefrogfuelfundgalagamegeargemsgiftgirlglowgoodgraygrimgurugushgyrohalfhanghardhawkheathelphighhillholyhopehornhutsicedideaidleinchinkyintoirisironitemjadejazzjoinjoltjowljudojugsjumpjunkjurykeepkenokeptkeyskickkilnkingkitekiwiknoblamblavalazyleaflegsliarlimplionlistlogoloudloveluaulucklungmainmanymathmazememomenumeowmildmintmissmonknailnavyneednewsnextnoonnotenumbobeyoboeomitonyxopenovalowlspaidpartpeckplaypluspoempoolposepuffpumapurrquadquizraceramprealredorichroadrockroofrubyruinrunsrustsafesagascarsetssilkskewslotsoapsolosongstubsurfswantacotasktaxitenttiedtimetinytoiltombtoystriptunatwinuglyundouniturgeuservastveryvetovialvibeviewvisavoidvowswallwandwarmwaspwavewaxywebswhatwhenwhizwolfworkyankyawnyellyogayurtzapszerozestzinczonezoom"
	bytewordsNum            = 256
	bytewordLength          = 4
	minimalBytewordLength   = 2
	dim                     = 26
	invalidBytewordsMessage = "Invalid Bytewords: "
)

var bytewordsLookUpTable []int

func getWord(index int) string {
	return bytewords[index*bytewordLength : (index*bytewordLength)+bytewordLength]
}

func getMinimalWord(index int) string {
	byteword := getWord(index)
	return fmt.Sprintf("%c%c", byteword[0], byteword[bytewordLength-1])
}

func addCRC(str string) string {
	crc := getCRCHex([]byte(str))
	return str + crc
}

func encodeWithSeparator(word string, separator string) string {
	crcAppendedWord := addCRC(word)
	crcWordBuff, _ := hex.DecodeString(crcAppendedWord)
	result := make([]string, len(crcWordBuff))
	for i, w := range crcWordBuff {
		result[i] = getWord(int(w))
	}
	return strings.Join(result, separator)
}

func encodeMinimal(word string) string {
	crcAppendedWord := addCRC(word)
	crcWordBuff, _ := hex.DecodeString(crcAppendedWord)
	var result strings.Builder
	for _, w := range crcWordBuff {
		result.WriteString(getMinimalWord(int(w)))
	}
	return result.String()
}

func decodeWord(word string, wordLength int) (string, error) {
	if len(word) != wordLength {
		return "", errors.New(invalidBytewordsMessage + "word.length does not match wordLength provided")
	}

	if len(bytewordsLookUpTable) == 0 {
		arrayLen := dim * dim
		bytewordsLookUpTable = make([]int, arrayLen)
		for i := range bytewordsLookUpTable {
			bytewordsLookUpTable[i] = -1
		}

		for i := 0; i < bytewordsNum; i++ {
			byteword := getWord(i)
			x := byteword[0] - 'a'
			y := byteword[3] - 'a'
			offset := int(y)*dim + int(x)
			bytewordsLookUpTable[offset] = i
		}
	}

	x := word[0] - 'a'
	y := word[1]
	if wordLength == bytewordLength {
		y = word[3]
	}
	y -= 'a'

	if x < 0 || x >= dim || y < 0 || y >= dim {
		return "", errors.New(invalidBytewordsMessage + "invalid word")
	}

	offset := int(y)*dim + int(x)
	value := bytewordsLookUpTable[offset]

	if value == -1 {
		return "", errors.New(invalidBytewordsMessage + "value not in lookup table")
	}

	if wordLength == bytewordLength {
		byteword := getWord(value)
		if word[1] != byteword[1] || word[2] != byteword[2] {
			return "", errors.New(invalidBytewordsMessage + "invalid middle letters of word")
		}
	}

	return fmt.Sprintf("%02x", value), nil
}

func _decode(str string, separator string, wordLength int) (string, error) {
	var words []string
	if wordLength == bytewordLength {
		words = strings.Split(str, separator)
	} else {
		words = partition(str, 2)
	}
	var decodedString strings.Builder
	for _, word := range words {
		decodedWord, err := decodeWord(word, wordLength)
		if err != nil {
			return "", err
		}
		decodedString.WriteString(decodedWord)
	}

	decodedHex := decodedString.String()
	if len(decodedHex) < 10 {
		return "", errors.New(invalidBytewordsMessage + "invalid decoded string length")
	}

	body, bodyChecksum := split([]byte(decodedHex), 4)
	checksum := getCRCHex(body)
	if checksum != hex.EncodeToString(bodyChecksum) {
		return "", errors.New(invalidBytewordsMessage + "invalid checksum")
	}

	return hex.EncodeToString(body), nil
}

func BytewordDecode(str string, style string) (string, error) {
	switch style {
	case "standard":
		return _decode(str, " ", bytewordLength)
	case "uri":
		return _decode(str, "-", bytewordLength)
	case "minimal":
		return _decode(str, "", minimalBytewordLength)
	default:
		return "", fmt.Errorf("invalid style %s", style)
	}
}

func BytewordEncode(str string, style string) (string, error) {
	switch style {
	case "standard":
		return encodeWithSeparator(str, " "), nil
	case "uri":
		return encodeWithSeparator(str, "-"), nil
	case "minimal":
		return encodeMinimal(str), nil
	default:
		return "", fmt.Errorf("invalid style %s", style)
	}
}
