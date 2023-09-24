package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/paas"
)

// ExpandUsers converts terraform representation of list of users to api representation.
func (s service) ExpandUsers(tfList []interface{}, forDatabase bool) []*paas.User {
	if len(tfList) == 0 {
		return nil
	}

	var users []*paas.User

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		user := s.toInterface().ExpandUser(tfMap, forDatabase)
		if user == nil {
			continue
		}

		users = append(users, user)
	}

	return users
}

// ExpandUser converts terraform representation of service user to api representation.
// If forDatabase is true, user is considered a database user.
func (s service) ExpandUser(tfMap map[string]interface{}, forDatabase bool) *paas.User {
	if tfMap == nil {
		return nil
	}

	user := &paas.User{}

	if v, ok := tfMap["name"].(string); ok && v != "" {
		user.Name = aws.String(v)
		delete(tfMap, "name")
	}

	if forDatabase {
		user.Parameters = s.toInterface().expandDatabaseUserParameters(tfMap)
	} else {
		user.Parameters = s.toInterface().expandUserParameters(tfMap)
	}

	return user
}

// ExpandDatabases converts terraform representation of list of databases to api representation.
func (s service) ExpandDatabases(tfList []interface{}) []*paas.Database {
	if len(tfList) == 0 {
		return nil
	}

	var databases []*paas.Database

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		database := s.toInterface().ExpandDatabase(tfMap)
		if database == nil {
			continue
		}

		databases = append(databases, database)
	}

	return databases
}

// ExpandDatabase converts terraform representation of database to api representation.
func (s service) ExpandDatabase(tfParameters map[string]interface{}) *paas.Database {
	if tfParameters == nil {
		return nil
	}

	database := &paas.Database{}

	if v, ok := tfParameters["backup_enabled"].(bool); ok {
		database.BackupEnabled = aws.Bool(v)
		delete(tfParameters, "backup_enabled")
	}

	if v, ok := tfParameters["name"].(string); ok && v != "" {
		database.Name = aws.String(v)
		delete(tfParameters, "name")
	}

	if v, ok := tfParameters["user"].([]interface{}); ok && len(v) > 0 {
		database.Users = s.toInterface().ExpandUsers(v, true)
		delete(tfParameters, "user")
	}

	database.Parameters = s.toInterface().expandDatabaseParameters(tfParameters)

	return database
}

// expandUserParameters converts terraform representation of service-specific user parameters
// // to api representation.
//
// If PaaS service has specific user parameters, it should override this method.
func (s service) expandUserParameters(_ map[string]interface{}) UserParameters {
	return nil
}

// expandDatabaseParameters converts terraform representation of service-specific database parameters
// to api representation.
//
// If PaaS service has specific database parameters, it should override this method.
func (s service) expandDatabaseParameters(_ map[string]interface{}) DatabaseParameters {
	return nil
}

// expandDatabaseUserParameters converts terraform representation of service-specific database user parameters
// to api representation.
//
// If PaaS service has specific database users parameters, it should override this method.
func (s service) expandDatabaseUserParameters(_ map[string]interface{}) DatabaseUserParameters {
	return nil
}

// FlattenServiceParametersUsersDatabases converts all blocks of service-specific parameters
// from api to terraform representation.
//
// It's Expand analogue is represented by three separate methods:
// ExpandServiceParameters (overridden on service level), ExpandUsers and ExpandDatabases,
// because these blocks are separated in api representation.
func (s service) FlattenServiceParametersUsersDatabases(
	serviceParameters ServiceParameters,
	users []*paas.User,
	databases []*paas.Database,
) map[string]interface{} {
	tfMap := s.toInterface().flattenServiceParameters(serviceParameters)

	if s.usersEnabled {
		tfMap["user"] = s.toInterface().FlattenUsers(users, false)
	}

	if s.databasesEnabled {
		tfMap["database"] = s.toInterface().FlattenDatabases(databases)
	}

	return tfMap
}

// FlattenUsers converts api representation of list of users to terraform representation.
func (s service) FlattenUsers(users []*paas.User, forDatabase bool) []interface{} {
	if len(users) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, user := range users {
		if user == nil {
			continue
		}

		tfList = append(tfList, s.toInterface().FlattenUser(user, forDatabase))
	}

	return tfList
}

// FlattenUser converts api representation of service user to terraform representation.
// If forDatabase is true, user is considered a database user.
func (s service) FlattenUser(user *paas.User, forDatabase bool) map[string]interface{} {
	if user == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v := user.Name; v != nil {
		tfMap["name"] = v
	}

	var parameters map[string]interface{}
	if forDatabase {
		parameters = s.toInterface().flattenDatabaseParameters(user.Parameters)
	} else {
		parameters = s.toInterface().flattenUserParameters(user.Parameters)
	}

	for k, v := range parameters {
		tfMap[k] = v
	}

	return tfMap
}

// FlattenDatabases converts api representation of list of databases to terraform representation.
func (s service) FlattenDatabases(databases []*paas.Database) []interface{} {
	if len(databases) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, database := range databases {
		if database == nil {
			continue
		}

		tfList = append(tfList, s.toInterface().FlattenDatabase(database))
	}

	return tfList
}

// FlattenDatabase converts api representation of database to terraform representation.
func (s service) FlattenDatabase(database *paas.Database) map[string]interface{} {
	if database == nil {
		return map[string]interface{}{}
	}

	tfMap := map[string]interface{}{}

	if v := database.BackupEnabled; v != nil {
		tfMap["backup_enabled"] = v
	}

	if v := database.Name; v != nil {
		tfMap["name"] = v
	}

	if v := database.Users; v != nil {
		tfMap["user"] = s.toInterface().FlattenUsers(database.Users, true)
	}

	for k, v := range s.toInterface().flattenDatabaseParameters(database.Parameters) {
		tfMap[k] = v
	}

	return tfMap
}

// flattenUserParameters converts api representation of service-specific user parameters
// to terraform representation.
//
// If PaaS service has specific user parameters, it should override this method.
func (s service) flattenUserParameters(_ UserParameters) map[string]interface{} {
	return nil
}

// flattenDatabaseParameters converts api representation of service-specific database parameters
// to terraform representation.
//
// If PaaS service has specific database parameters, it should override this method.
func (s service) flattenDatabaseParameters(_ DatabaseParameters) map[string]interface{} {
	return nil
}

// flattenDatabaseUserParameters converts api representation of service-specific database user parameters
// to terraform representation.
//
// If PaaS service has specific database users parameters, it should override this method.
func (s service) flattenDatabaseUserParameters(_ DatabaseUserParameters) map[string]interface{} {
	return nil
}
