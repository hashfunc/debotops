package v1alpha1

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"

	"k8s.io/apimachinery/pkg/util/rand"
)

const (
	RevisionKeyFormat = "%s/revision"
)

func RevisionKey() string {
	return fmt.Sprintf(RevisionKeyFormat, GroupVersion.Group)
}

func Revision(spec interface{}) (string, error) {
	data, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}

	hash64a := fnv.New64a()

	_, err = hash64a.Write(data)
	if err != nil {
		return "", err
	}

	revision := rand.SafeEncodeString(
		strconv.FormatUint(hash64a.Sum64(), 10),
	)

	return revision, nil
}
