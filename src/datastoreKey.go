package main

import (
	"appengine/datastore"
	"bytes"
	"datastorekey"
	"fmt"
	"os"
	_ "quiet"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	result, err := process()
	if err == nil {
		fmt.Println(result)
	} else {
		fmt.Fprintln(os.Stderr, "Failure")
		fmt.Fprintln(os.Stderr, err)
	}
}

func usage() {
	fmt.Println("Usage ")
	fmt.Printf(" %s decode sourceKey\n", os.Args[0])
	fmt.Printf(" %s extract field sourceKey\n", os.Args[0])
	fmt.Printf(" %s translate [-appID=a] [-namespace=ns] sourceKey\n", os.Args[0])
	fmt.Printf(" %s encode kind appID namespace name id [parentkey]\n", os.Args[0])
	os.Exit(1)
}

func process() (string, error) {
	action := os.Args[1]
	switch action {
	case "decode":
		sourceKeyString := os.Args[2]
		key, err := datastore.DecodeKey(sourceKeyString)
		if err != nil {
			return "", err
		}
		return keyContent(key), nil
	case "extract":
		if len(os.Args) < 4 {
			usage()
		}
		field := os.Args[2]
		sourceKeyString := os.Args[3]
		key, err := datastore.DecodeKey(sourceKeyString)
		if err != nil {
			return "", err
		}
		return extract(key, field)
	case "translate":
		if len(os.Args) < 4 {
			usage()
		}
		translateAppId := false
		targetAppId := ""
		translateNamespace := false
		targetNamespace := ""
		var sourceKeyString string

		if strings.HasPrefix(os.Args[2], "-appID=") {
			translateAppId = true
			targetAppId = os.Args[2][len("-appID="):]
			if strings.HasPrefix(os.Args[3], "-namespace=") {
				translateNamespace = true
				targetNamespace = os.Args[3][len("-namespace="):]
				if len(os.Args) < 4 {
					usage()
				}
				sourceKeyString = os.Args[4]
			} else {
				sourceKeyString = os.Args[3]
			}
		} else if strings.HasPrefix(os.Args[2], "-namespace=") {
			translateNamespace = true
			targetNamespace = os.Args[2][len("-namespace="):]
			sourceKeyString = os.Args[3]
		}
		return translateKeyString(sourceKeyString, translateAppId, targetAppId, translateNamespace, targetNamespace)
	case "encode":
		if len(os.Args) < 7 {
			usage()
		}
		kind := os.Args[2]
		appID := os.Args[3]
		namespace := os.Args[4]
		name := os.Args[5]
		intID := os.Args[6]
		parentKeyStr := ""
		if len(os.Args) >= 8 {
			parentKeyStr = os.Args[7]
		}
		key, err := encode(kind, appID, namespace, name, intID, parentKeyStr)
		if err == nil {
			return key.Encode(), nil
		} else {
			return "", err
		}
	default:
		usage()
		return "", fmt.Errorf("Will not happen.")
	}
}

func keyContent(key *datastore.Key) string {
	var buffer bytes.Buffer

	indent := " "
	k := key

	for k != nil {
		if k != key {
			buffer.WriteString(indent)
			buffer.WriteString(" Parent key:")
			buffer.WriteString("\n")
		}

		buffer.WriteString(indent)
		buffer.WriteString(" Kind=")
		buffer.WriteString(k.Kind())
		buffer.WriteString("\n")

		buffer.WriteString(indent)
		buffer.WriteString(" AppID=")
		buffer.WriteString(k.AppID())
		buffer.WriteString("\n")

		buffer.WriteString(indent)
		buffer.WriteString(" Namespace=")
		buffer.WriteString(k.Namespace())
		buffer.WriteString("\n")

		buffer.WriteString(indent)
		buffer.WriteString(" Name=")
		buffer.WriteString(k.StringID())
		buffer.WriteString("\n")

		buffer.WriteString(indent)
		buffer.WriteString(" Id=")
		buffer.WriteString(fmt.Sprint(k.IntID()))
		buffer.WriteString("\n")

		indent = "   " + indent
		k = k.Parent()
	}

	return buffer.String()
}

func extract(key *datastore.Key, field string) (string, error) {
	switch strings.ToLower(field) {
	case "kind":
		return key.Kind(), nil
	case "appid":
		return key.AppID(), nil
	case "namespace":
		return key.Namespace(), nil
	case "name":
		return key.StringID(), nil
	case "id":
		return fmt.Sprintf("%v", key.IntID()), nil
	case "parent":
		if key.Parent() == nil {
			return "", nil
		} else {
			return key.Parent().Encode(), nil
		}
	default:
		return "", fmt.Errorf("Unsupported field [%v]. Supported fields are kind, appID, namespace, name, id, parentkey.", field)
	}
}

func encode(kind, appID, namespace, stringID, intIDStr, parentKeyString string) (*datastore.Key, error) {
	var parent *datastore.Key = nil
	if parentKeyString != "" {
		var err error
		parent, err = datastore.DecodeKey(parentKeyString)
		if err != nil {
			return nil, err
		}
	}
	return datastorekey.CreateKey(nil, appID, namespace, kind, stringID, intID64(intIDStr), parent)
}

func translate(sourceKey *datastore.Key, translateAppId bool, targetAppId string, translateNamespace bool, targetNamespace string) (*datastore.Key, error) {
	if !translateAppId {
		targetAppId = sourceKey.AppID()
	}
	if !translateNamespace {
		targetNamespace = sourceKey.Namespace()
	}
	var translatedParent *datastore.Key = nil
	if sourceKey.Parent() != nil {
		var err error
		translatedParent, err = translate(sourceKey.Parent(), translateAppId, targetAppId, translateNamespace, targetNamespace)
		if err != nil {
			return nil, err
		}
	}
	return datastorekey.CreateKey(nil, targetAppId, targetNamespace, sourceKey.Kind(), sourceKey.StringID(), sourceKey.IntID(), translatedParent)
}

func translateKeyString(sourceKeyString string, translateAppId bool, targetAppId string, translateNamespace bool, targetNamespace string) (string, error) {
	sourceKey, err := datastore.DecodeKey(sourceKeyString)
	if err != nil {
		return "", err
	}
	translatedKey, err := translate(sourceKey, translateAppId, targetAppId, translateNamespace, targetNamespace)
	if err == nil {
		return translatedKey.Encode(), nil
	} else {
		return "", err
	}
}

func intID64(intIDstr string) int64 {
	if intIDstr == "" {
		return 0
	}
	intID64, _ := strconv.ParseInt(intIDstr, 10, 64)
	return intID64
}
