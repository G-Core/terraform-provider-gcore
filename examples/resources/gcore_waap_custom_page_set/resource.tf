provider "gcore" {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_cdn_resource" "example" {
  cname  = "api.example.com"
  origin = "origin.example.com"
  options {
    waap { value = true }
  }
}

resource "gcore_waap_domain" "mydomain" {
  name = gcore_cdn_resource.example.cname
}

resource "gcore_waap_custom_page_set" "example" {
  name    = "example-custom-pages"
  domains = [gcore_waap_domain.mydomain.id]

  block {
    logo    = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAACXBIWXMAAA7EAAAOxAGVKw4bAAAA0ElEQVQ4y62TMQ6DMAxFH1GlDkgsbJXYOAecowsnYO4F2DpxGdZeghEpS4+QqerS5YOiNEVtwZKl6Nvfib+dhM92BAqd78CDL60GBsABk9wJq9eIBugBCzRA6sVSYVY5JlagB25AtnJJppw+9mzrkU/BLUbYXMSG7Qx64kx+Aq0Xb4XNRRpxFrWd17MRofAKFMKMp4kTl1JK/2oTUB7YwcIWzkAXEbFT7K2FUMRcoxqBq3wUlsdEjI0RoAIu8irYBRvbyk2LtMsq//WZkq3f+QXg4Ta9dWcCCAAAAABJRU5ErkJggg=="
    header  = "Initial Block Header"
    title   = "Initial Block Title"
    text    = "This is the initial block page text."
    enabled = true
  }

  block_csrf {
    header  = "Initial CSRF Headersss"
    title   = "Initial CSRF Title"
    text    = "This is the initial CSRF block page text."
    enabled = true
  }

  captcha {
    header  = "Initial Captcha Header"
    title   = "Initial Captcha Title"
    text    = "This is the initial captcha page texts."
    error   = "Initial captcha error message."
    enabled = true
  }

  cookie_disabled {
    header  = "Initial Cookie Header"
    text    = "Initial cookie disabled text."
    enabled = true
  }

  handshake {
    header  = "Initial Handshake Header"
    title   = "Initial Handshake Title"
    enabled = true
  }

  javascript_disabled {
    header  = "Initial JS Header"
    text    = "Initial JavaScript disabled text."
    enabled = true
  }
}
