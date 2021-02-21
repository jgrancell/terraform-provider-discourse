package discourse

import (
  "context"
  "encoding/json"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSetting() *schema.Resource {
  return &schema.Resource{
    CreateContext: resourceSettingCreate,
    ReadContext:   resourceSettingRead,
    UpdateContext: resourceSettingUpdate,
    DeleteContext: resourceSettingDelete,
    Schema: map[string]*schema.Schema{
      "key": &schema.Schema{
        Type: schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "value": &schema.Schema{
        Type: schema.TypeString,
        Required: true,
      },
      "default": &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
        Computed: true,
      },
      "description": &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
        Computed: true,
      },
    },
    SchemaVersion: 1,
  }
}

func resourceSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics
  conf := m.(Config)

  // Building our request map
  newSetting := make(map[string]string)
  newSetting[d.Get("key").(string)] = d.Get("value").(string)

  // Turning resource struct into jsonBody
  settingJson, err := json.Marshal(newSetting)
  if err != nil {
    return diag.FromErr(err)
  }

  request := &ApiRequest{
    Method: "PUT",
    Endpoint: "/admin/site_settings/"+d.Get("key").(string)+".json",
    JsonBody: string(settingJson),
    Config: conf,
  }

  _, reqErr, ok := request.Call()
  if ! ok {
    diags = append(diags, reqErr)
    return diags
  }

  // Terraform requires strings for Ids, so we have to convert the integer Discourse gives us
	d.SetId(d.Get("key").(string))
  return resourceSettingRead(ctx, d, m)
}

func resourceSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  var diags diag.Diagnostics
  conf := m.(Config)

  request := &ApiRequest{
    Method: "GET",
    Endpoint: "/admin/site_settings.json",
    Config: conf,
  }

  body, reqErr, ok := request.Call()
  if ! ok {
    diags = append(diags, reqErr)
    return diags
  }

  var settings SettingsResponse
  j := json.Unmarshal([]byte(body), &settings)
  if j != nil {
    return diag.FromErr(j)
  }

  var value string
  var def string
  var desc string
  for _, s := range settings.Settings {
    if s.Name == d.Get("key") {
      value = s.Value
      def = s.Default
      desc = s.Description
    }
  }

  d.Set("value", value)
  d.Set("default", def)
  d.Set("description", desc)
  return diags
}

func resourceSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  var diags diag.Diagnostics
  conf := m.(Config)

  if d.HasChange("value") {
    // Building our request map
    newSetting := make(map[string]string)
    newSetting[d.Get("key").(string)] = d.Get("value").(string)

    // Turning resource struct into jsonBody
    settingJson, err := json.Marshal(newSetting)
    if err != nil {
      return diag.FromErr(err)
    }

    request := &ApiRequest{
      Method: "PUT",
      Endpoint: "/admin/site_settings/"+d.Get("key").(string)+".json",
      JsonBody: string(settingJson),
      Config: conf,
    }

    _, reqErr, ok := request.Call()
    if ! ok {
      diags = append(diags, reqErr)
      return diags
    }
  }
  return resourceSettingRead(ctx, d, m)
}

func resourceSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  var diags diag.Diagnostics
  conf := m.(Config)

  // Building our request map
  newSetting := make(map[string]string)
  newSetting[d.Get("key").(string)] = d.Get("default").(string)

  // Turning resource struct into jsonBody
  settingJson, err := json.Marshal(newSetting)
  if err != nil {
    return diag.FromErr(err)
  }

  request := &ApiRequest{
    Method: "PUT",
    Endpoint: "/admin/site_settings/"+d.Id()+".json",
    JsonBody: string(settingJson),
    Config: conf,
  }

  _, reqErr, ok := request.Call()
  if ! ok {
    diags = append(diags, reqErr)
    return diags
  }
  d.SetId("")
  return diags
}
