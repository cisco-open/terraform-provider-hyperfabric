resource "hyperfabric_user" "user1" {
  email       = "test@test.be"
}

resource "hyperfabric_user" "user2" {
  email       = "test2@test.be"
  role = "ADMIN"
  labels = ["blue", "green", "red"]
  enabled = false
}

data "hyperfabric_user" "user1" {
  email      = hyperfabric_user.user1.email
}

output "datasource_user_auth_provider" {
    value = data.hyperfabric_user.user1.auth_provider
}
