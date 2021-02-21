package discourse

import (
  "context"
  "encoding/json"
  "strconv"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
  return &schema.Resource{
    CreateContext: resourceGroupCreate,
    ReadContext:   resourceGroupRead,
    UpdateContext: resourceGroupUpdate,
    DeleteContext: resourceGroupDelete,
    Schema: map[string]*schema.Schema{
      "name" : &schema.Schema{
        Type: schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "mentionable_level" : &schema.Schema{
        Type: schema.TypeInt,
        Optional: true,
        Computed: true,
      },
      "messageable_level" : &schema.Schema{
        Type: schema.TypeInt,
        Optional: true,
        Computed: true,
      },
      "visibility_level" : &schema.Schema{
        Type: schema.TypeInt,
        Optional: true,
        Computed: true,
      },
      "primary_group" : &schema.Schema{
        Type: schema.TypeBool,
        Optional: true,
        Computed: true,
      },
      "title" : &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
        Computed: true,
      },
      "grant_trust_level" : &schema.Schema{
        Type: schema.TypeInt,
        Optional: true,
        Computed: true,
      },
      "flair_url" : &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
        Computed: true,
      },
      "flair_bg_color" : &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
        Computed: true,
      },
      "flair_color" : &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
        Computed: true,
      },
      "public_admission" : &schema.Schema{
        Type: schema.TypeBool,
        Optional: true,
        Computed: true,
      },
      "public_exit" : &schema.Schema{
        Type: schema.TypeBool,
        Optional: true,
        Computed: true,
      },
      "allow_membership_requests" : &schema.Schema{
        Type: schema.TypeBool,
        Optional: true,
        Computed: true,
      },
      "full_name" : &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
        Computed: true,
      },
      "default_notification_level" : &schema.Schema{
        Type: schema.TypeInt,
        Optional: true,
        Computed: true,
      },
      "members_visibility_level" : &schema.Schema{
        Type: schema.TypeInt,
        Optional: true,
        Computed: true,
      },
    },
    SchemaVersion: 1,
  }
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics
  conf := m.(Config)

  // Building our request structs
  newGroup := &Group {
    Name: d.Get("name").(string),
  }

  if arg, argOk := d.GetOk("mentionable_level"); argOk {
    newGroup.MentionableLevel = arg.(int)
  }
  if arg, argOk := d.GetOk("messageable_level"); argOk {
    newGroup.MessageableLevel = arg.(int)
  }
  if arg, argOk := d.GetOk("visibility_level"); argOk {
    newGroup.VisibilityLevel = arg.(int)
  }
  if arg, argOk := d.GetOk("primary_group"); argOk {
    newGroup.PrimaryGroup = arg.(bool)
  }
  if arg, argOk := d.GetOk("title"); argOk {
    newGroup.Title = arg.(string)
  }
  if arg, argOk := d.GetOk("grant_trust_level"); argOk {
    newGroup.GrantTrustLevel = arg.(int)
  }
  if arg, argOk := d.GetOk("flair_url"); argOk {
    newGroup.FlairUrl = arg.(string)
  }
  if arg, argOk := d.GetOk("flair_bg_color"); argOk {
    newGroup.FlairBackgroundColor = arg.(string)
  }
  if arg, argOk := d.GetOk("flair_color"); argOk {
    newGroup.FlairColor = arg.(string)
  }
  if arg, argOk := d.GetOk("public_admission"); argOk {
    newGroup.PublicAdmission = arg.(bool)
  }
  if arg, argOk := d.GetOk("public_exit"); argOk {
    newGroup.PublicExit = arg.(bool)
  }
  if arg, argOk := d.GetOk("allow_membership_requests"); argOk {
    newGroup.AllowMembershipRequests = arg.(bool)
  }
  if arg, argOk := d.GetOk("full_name"); argOk {
    newGroup.FullName = arg.(string)
  }
  if arg, argOk := d.GetOk("default_notification_level"); argOk {
    newGroup.DefaultNotificationLevel = arg.(int)
  }
  if arg, argOk := d.GetOk("members_visibility_level"); argOk {
    newGroup.MembersVisibilityLevel = arg.(int)
  }

  groupMap := make(map[string]*Group)
  groupMap["group"] = newGroup

  // Turning resource struct into jsonBody
  groupJson, err := json.Marshal(groupMap)
  if err != nil {
    return diag.FromErr(err)
  }
  log(string(groupJson))

  request := &ApiRequest{
    Method: "POST",
    Endpoint: "/admin/groups.json",
    JsonBody: string(groupJson),
    Config: conf,
  }

  body, reqErr, ok := request.Call()
  if ! ok {
    diags = append(diags, reqErr)
    return diags
  }

  var group GroupResponse
	j := json.Unmarshal([]byte(body), &group)
	if j != nil {
		return diag.FromErr(j)
	}

  d.Set("name", group.Group.Name)
  d.Set("mentionable_level", group.Group.MentionableLevel)
  d.Set("messageable_level", group.Group.MessageableLevel)
  d.Set("visibility_level", group.Group.VisibilityLevel)
  d.Set("primary_group", group.Group.PrimaryGroup)
  d.Set("title", group.Group.Title)
  d.Set("grant_trust_level", group.Group.GrantTrustLevel)
  d.Set("flair_url", group.Group.FlairUrl)
  d.Set("flair_bg_color", group.Group.FlairBackgroundColor)
  d.Set("flair_color", group.Group.FlairColor)
  d.Set("public_admission", group.Group.PublicAdmission)
  d.Set("public_exit", group.Group.PublicExit)
  d.Set("allow_membership_requests", group.Group.AllowMembershipRequests)
  d.Set("full_name", group.Group.FullName)
  d.Set("default_notification_level", group.Group.DefaultNotificationLevel)
  d.Set("members_visibility_level", group.Group.MembersVisibilityLevel)
  // Terraform requires strings for Ids, so we have to convert the integer Discourse gives us
  log("ID is "+strconv.Itoa(group.Group.Id))
	d.SetId(strconv.Itoa(group.Group.Id))
  return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  var diags diag.Diagnostics
  conf := m.(Config)

  request := &ApiRequest{
    Method: "DELETE",
    Endpoint: "/admin/groups/"+d.Id()+".json",
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
