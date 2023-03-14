resource "example_server" "my-server-name" {
	uuid_count = "1"
}

resource "example_did" "did" {
	method = "web"
	url = "antunovic.nz"
}
