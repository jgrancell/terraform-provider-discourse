terraform {
  required_providers {
    discourse = {
      source = "jgrancell/discourse"
      version = "0.1.0"
    }
  }
}

// You must either set `host`, `username`, and `token` here
// or set DISCOURSE_HOST, DISCOURSE_USERNAME, and DISCOURSE_TOKEN
// in your environment
provider "discourse" {}
