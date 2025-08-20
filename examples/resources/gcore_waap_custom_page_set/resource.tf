resource "gcore_waap_custom_page_set" "example_waap_custom_page_set" {
  name = "x"
  block = {
    enabled = true
    header = "xxx"
    logo = "logo"
    text = "xxxxxxxxxxxxxxxxxxxx"
    title = "xxx"
  }
  block_csrf = {
    enabled = true
    header = "xxx"
    logo = "logo"
    text = "xxxxxxxxxxxxxxxxxxxx"
    title = "xxx"
  }
  captcha = {
    enabled = true
    error = "xxxxxxxxxx"
    header = "xxx"
    logo = "logo"
    text = "xxxxxxxxxxxxxxxxxxxx"
    title = "xxx"
  }
  cookie_disabled = {
    enabled = true
    header = "xxx"
    text = "xxxxxxxxxxxxxxxxxxxx"
  }
  domains = [1]
  handshake = {
    enabled = true
    header = "xxx"
    logo = "logo"
    title = "xxx"
  }
  javascript_disabled = {
    enabled = true
    header = "xxx"
    text = "xxxxxxxxxxxxxxxxxxxx"
  }
}
