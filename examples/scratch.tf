locals {
  settings = {
    "title" = "Example Forum"
    "site_description" = "An enthusiast community of examples"
    "short_site_description" = "Hello World"
    "contact_email" = "noreply@example.com"
    "contact_url" = "https://altavista.com"
    "notification_email" = "noreply@example.com"
    "site_contact_username" = "system"
    "company_name" = "Example Company"
    "city_for_disputes" = "Sanctuary"
    "top_menu" = "categories|new|latest"
    "base_font" = "roboto"
    "heading_font" = "roboto"
    "login_required" = "true"
    "enable_discord_logins" = "true"
    "discord_client_id" = "abcd1234"
    "discord_secret" = "1234abcd"
    "discord_trusted_guilds" = "987654321"
    "bootstrap_mode_min_users" = "0"

  }
}

resource "discourse_group" "approved_members" {
  name = "SecondGroup"
}

resource "discourse_setting" "settings" {
  for_each = local.settings
  key   = each.key
  value = each.value
}
