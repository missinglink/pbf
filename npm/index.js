
var proc = require('./lib/process'),
    jsonStream = require('./lib/jsonStream');

function api( config ){

  // format flags
  var flags = [];
  // flags.push( util.format( '-tags=%s', config.tags ) );
  // if( config.hasOwnProperty( 'leveldb' ) ){
  //   flags.push( util.format( '-leveldb=%s', config.leveldb ) );
  // }
  // flags.push( config.file );

  // spawn child process
  var child = proc( flags );

  // create a json read stream
  return jsonStream( child );
}

module.exports = api;
