resource "dynatrace_document" "#name#" {
  name = "#name#"
  type = "launchpad"
  content = jsonencode({
    "background" : "default",
    "containerList" : {
      "containers" : [
        {
          "blocks" : [
            {
              "appearance" : "list",
              "content" : [
                {
                  "action" : {
                    "type" : "openExternalLink",
                    "url" : "https://www.google.at"
                  },
                  "description" : "",
                  "id" : "3b7b55b0-fe97-432a-b80a-7f39595903a0",
                  "title" : "#name#",
                  "type" : "link"
                }
              ],
              "contentType" : "static",
              "id" : "ac3d92c2-ed43-43d7-9b00-7b846038382e",
              "properties" : {
                "expanded" : true
              },
              "type" : "links"
            }
          ],
          "horizontalLayoutWeight" : 1
        }
      ]
    },
    "icon" : "default",
    "schemaVersion" : 2
  })
  private = false
}
