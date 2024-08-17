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

# resource "hyperfabric_bearer_token" "token1" {
#   name       = "token1"
# }

# resource "hyperfabric_bearer_token" "token2" {
#   name       = "token2"
#   description = "My New Second Token"
#   not_after = "2025-09-03T08:00:00.000Z"
#   not_before = "2024-09-03T08:00:00.000Z"
#   scope = "ADMIN"
# }

# data "hyperfabric_bearer_token" "token1" {
#   name      = hyperfabric_bearer_token.token1.name
# }

# output "datasource_bearer_token_description" {
#     value = data.hyperfabric_bearer_token.token1.description
# }