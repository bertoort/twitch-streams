$(document).ready(function () {
  var streams = $('.streamName')
  for (var i = 0; i < streams.length; i++) {
    var stream = streams[i].innerHTML
    $.get("https://api.twitch.tv/kraken/streams/" + stream)
        .done(function( data ) {
        if (data.stream) {
          $('.' + data.stream.channel.name + 'Status').css("color", "#43AC6A");
        }
      });
  }
})
