package handlers

//
import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/hyqe/brigand/internal/storage"
	"net/http"
	"time"
)

func CheckSymlinkHash(symlink *storage.Symlink, hmacSecret string) bool {
	path := fmt.Sprintf("expiration=%s&id=%s&name=%s", symlink.Expiration, symlink.Id, symlink.Name)

	h := hmac.New(sha256.New, []byte(hmacSecret))
	_, err := h.Write([]byte(path))
	if err != nil {
		panic(err)
	}
	newhash := hex.EncodeToString(h.Sum(nil))
	if newhash == symlink.Hash {
		return true
	}
	return false
}

func checkIfSymlinkIsExpired(expiration string) (bool, error) {
	theTime, err := time.Parse(time.RFC3339, expiration)
	if err != nil {
		return true, err
	}

	if time.Now().After(theTime) {
		return true, nil
	}

	return false, nil
}

func timeError(oldTime string) string {
	realTime, err := time.Parse(time.RFC3339, oldTime)
	if err != nil {
		return err.Error() + "\n\t Failed to parse a time object for better errors, no idea how its possible"
	}

	old, now := realTime.Format(time.DateTime), time.Now().Format(time.DateTime)

	return fmt.Sprintf("your symlink has expired: symlink's time: <%s> || Time of request process: <%s>", old, now)
}

func TakeSymlink(fileDownloader storage.FileDownloader, symlinkSecret string, getSymlink func(r *http.Request) *storage.Symlink) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symlink := getSymlink(r)

		isExpired, err := checkIfSymlinkIsExpired(symlink.Expiration)
		if err != nil {
			http.Error(w, "The Time string could not be parsed back into a time.Time object", http.StatusInternalServerError)
			return
		}
		if isExpired {
			// TODO: return the time of the symlink in human-readable & the current time right now
			http.Error(w, timeError(symlink.Expiration), http.StatusNotAcceptable)
			return
		}

		if !CheckSymlinkHash(symlink, symlinkSecret) {
			http.Error(w, "We did not send you this link. UNACCEPTABLE!!!!", http.StatusForbidden)
			return
		}

		err = fileDownloader(w, symlink.Id)
		if err != nil {
			http.Error(w, "Failed to retrieve your file by the given id in the url parameter", http.StatusInternalServerError)
		}
	}
}
