# Description:
#   Interfaces with github.com/fgeller/rad to search API documentation
#
# Configuration:
#   RAD_URL
#
# Commands:
#   hubot rad <pack> <entity> <member> - queries api documentation
#
# Author:
#   fgeller

module.exports = (robot) ->
  robot.respond /rad ([^ ]+) ([^ ]+) ([^ ]+)/i, (res) ->
    rad_url = process.env.RAD_URL
    pkg = res.match[1]
    ent = res.match[2]
    mem = res.match[3]
    robot.http("#{rad_url}/s?p=" + pkg + "&e=" + ent + "&m=" + mem)
      .get() (err, result, body) ->
        if err
          res.reply "Error while asking for pack [#{pkg}] entity [#{ent}] member [#{mem}]"
          return

        results = JSON.parse(body)
        res.reply "rad found #{results.length} entries:"
        for r in results
          robot.adapter.customMessage {
            channel: res.message.room
            attachments: [
              {
                title: "#{r.Entity} #{r.Member}",
                title_link: "#{rad_url}/ui/?doc=#{r.Target.substring(1)}",
                mrkdwn_in: ["text"]
              }
            ]
          }
