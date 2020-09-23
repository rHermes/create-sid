package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
)

var inputGUID = flag.String("guid", "", "the guid to convert")

// parseGUID parses a microsoft encoded style GUID.
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Encoding
func parseGUID(guid string) ([16]byte, error) {
	var ret [16]byte
	parts := strings.Split(guid, "-")
	if len(parts) != 5 ||
		len(parts[0]) != 8 ||
		len(parts[1]) != 4 ||
		len(parts[2]) != 4 ||
		len(parts[3]) != 4 ||
		len(parts[4]) != 12 {
		return ret, errors.New("Not a valid format")
	}
	k := 0
	for i, part := range parts {
		bts, err := hex.DecodeString(part)
		if err != nil {
			return ret, err
		}

		// The first 3 components are in little endian, so we must reverse them
		if i < 3 {
			for j := len(bts)/2 - 1; j >= 0; j-- {
				opp := len(bts) - 1 - j
				bts[j], bts[opp] = bts[opp], bts[j]
			}
		}

		for _, c := range bts {
			ret[k] = c
			k++
		}
	}
	return ret, nil
}

// toSid converts from a GUID to a sid string with a comment
// it's meant for use in the templates
func toSid(guid string) (string, error) {
	bs, err := parseGUID(guid)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("0x%x /* uuid = %s */", bs, guid), nil
}

func main() {
	flag.Parse()

	if inputGUID == nil || *inputGUID == "" {
		fmt.Println("No guid given!")
	} else {
		sid, err := toSid(*inputGUID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(sid)
	}
}
