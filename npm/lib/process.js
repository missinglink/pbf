
var os = require('os'),
    util = require('util'),
    path = require('path'),
    child = require('child_process');

function spawn( flags ){

  // select correct executable to use for this system
  var exec = path.join(__dirname, '../../build', util.format( 'pbf.%s-%s', os.platform(), os.arch() ) );

  // spawn child process
  return child.spawn( exec, flags );
}

module.exports = spawn;
