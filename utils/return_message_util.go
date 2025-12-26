package utils

func MessageDataCreated(name, status string) string {
	return status + " to create " + name
}
func MessageDataUpdated(name, status string) string {
	return status + " to update " + name
}
func MessageDataDeleted(name, status string) string {
	return status + " to delete " + name
}
func MessageDataRead(name, status string) string {
	return status + " to retrieve " + name
}
func MessageLogin(status string) string {
	return "Login " + status
}
