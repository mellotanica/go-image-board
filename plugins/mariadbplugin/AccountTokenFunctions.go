package mariadbplugin

import (
	"bytes"
	"database/sql"
	"errors"
	"go-image-board/logging"

	uuid "github.com/satori/go.uuid"
)

//ValidateToken Validate a cookie token (true if valid cookie, false otherwise, error for reason or nil)
func (DBConnection *MariaDBPlugin) ValidateToken(userName string, tokenID string, ip string) error {
	var validTokenID sql.NullString
	var validTokenIP sql.NullString
	var userDisabled bool
	row := DBConnection.DBHandle.QueryRow("SELECT TokenID, IP, Disabled FROM Users WHERE Name = ?", userName)
	err := row.Scan(&validTokenID, &validTokenIP, &userDisabled)
	if userDisabled {
		return errors.New("Account disabled")
	}
	if err != nil && validTokenID.Valid && validTokenIP.Valid {
		//User's token in DB is blank
		logging.LogInterface.WriteLog("MariaDBPlugin", "ValidateToken", userName, "ERROR", []string{"Token Invalid", userName, tokenID, ip})
		return errors.New("Token invalid")
	}

	UUIDBytes := uuid.FromStringOrNil(tokenID)
	if uuid.Equal(UUIDBytes, uuid.UUID{}) == true {
		//Token provided is blank
		//logging.LogInterface.WriteLog("MariaDBPlugin", "ValidateToken", userName, "ERROR", []string{"Blank token provided", userName, tokenID, ip}) //This happens for ALL unauth users. Log spam.
		return errors.New("Token provided is blank")
	}

	if validTokenIP.String != ip {
		//Token is registered for a different IP
		logging.LogInterface.WriteLog("MariaDBPlugin", "ValidateToken", userName, "ERROR", []string{"Token for a different IP", userName, tokenID, ip})
		return errors.New("Token invalid")
	}

	if bytes.Equal(UUIDBytes.Bytes(), uuid.FromStringOrNil(validTokenID.String).Bytes()) == false {
		//Tokens do not match
		logging.LogInterface.WriteLog("MariaDBPlugin", "ValidateToken", userName, "ERROR", []string{"Tokens don't match", userName, tokenID, ip})
		return errors.New("Token invalid")
	}

	return nil
}

//GenerateToken Generate a cookie token (string token, or error)
func (DBConnection *MariaDBPlugin) GenerateToken(userName string, ip string) (string, error) {
	newToken, err := uuid.NewV4()
	if err != nil {
		logging.LogInterface.WriteLog("MariaDBPlugin", "GenerateToken", userName, "ERROR", []string{"Failed to generate a token!", userName, ip, err.Error()})
		return "", errors.New("failed to generate a token")
	}
	_, err = DBConnection.DBHandle.Exec("UPDATE Users SET TokenID=?, IP=? WHERE Name = ?", newToken.String(), ip, userName)
	if err != nil {
		logging.LogInterface.WriteLog("MariaDBPlugin", "GenerateToken", userName, "ERROR", []string{"Failed to save token", userName, ip, err.Error()})
		return "", errors.New("failed to generate a token, check if user exists")
	}
	return newToken.String(), nil
}

//RevokeToken Revokes a token (nil on success)
func (DBConnection *MariaDBPlugin) RevokeToken(userName string) error {
	_, err := DBConnection.DBHandle.Exec("UPDATE Users SET TokenID=NULL, IP=NULL WHERE Name = ?", userName)
	if err == nil {
		logging.LogInterface.WriteLog("MariaDBPlugin", "RevokeToken", userName, "SUCCESS", []string{"Token revoked!", userName})
	} else {
		logging.LogInterface.WriteLog("MariaDBPlugin", "RevokeToken", userName, "ERROR", []string{"Token not revoked", userName, err.Error()})
	}
	return err
}
