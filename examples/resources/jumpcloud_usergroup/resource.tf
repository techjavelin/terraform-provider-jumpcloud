resource "jumpcloud_usergroup" "example" {
  name        = "example"
  description = "Example Description"
  email       = "example@example.com"

  sudo = {
    enabled           = false
    passwpasswordless = false
  }

  ldap = {
    groups = [
      { name = "example" }
    ]
  }

  posix = [
    {
      id   = 1000
      name = "example"
    }
  ]

  radius = [
    {
      name  = "example"
      value = "example"
    }
  ]

  samba = false

  member_queries = [
    {
      query = {
        field    = "example"
        operator = "eq"
        value    = "example"
      }
    }
  ]

  notify = false
  auto   = false
}