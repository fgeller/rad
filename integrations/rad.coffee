# Description:
#   Interfaces with github.com/fgeller/rad to search API documentation
#
# Dependencies:
#  "ws": "0.8.0"
#
# Configuration:
#   RAD_URL
#   RAD_WS_URL
#
# Commands:
#   hubot rad <pack> <entity> <member> - queries api documentation
#
# Author:
#   fgeller

module.exports = (robot) ->
  robot.respond /rad ([^ ]+) ([^ ]+) ([^ ]+)/i, (res) ->
    rad_url = process.env.RAD_URL
    rad_ws_url = process.env.RAD_WS_URL
    pkg = res.match[1]
    pth = res.match[2]
    mem = res.match[3]

    WebSocket = require('ws')
    conn = new WebSocket(rad_ws_url)
    req = {"Limit": 100, "Pack": pkg, "Path": pth, "Member": mem}
    results = 0

    conn.on('message', (data) ->
      results++
      if results <= 5
        entry = JSON.parse(data)
        robot.adapter.customMessage {
          channel: res.message.room
          attachments: [
            {
              title: "#{entry.Namespace} #{entry.Member}",
              title_link: "#{rad_url}/?doc=#{entry.Target}",
              mrkdwn_in: ["text"]
            }
          ]
        })

    conn.on('open', () -> conn.send(JSON.stringify(req)))

    conn.on('close', () ->
      if results == 0
        res.reply "rad found no entries"
      else
        res.reply "rad found at least #{results} entries (display limited to 5).")

    conn.on('error', (err) ->
      res.reply "Error while asking for pack [#{pkg}] path [#{pth}] member [#{mem}]: " + err)
