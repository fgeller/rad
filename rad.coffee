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
        for r in results
          res.reply "#{r.Entity} #{r.Member} #{rad_url}/ui/?doc=#{r.Target}"

